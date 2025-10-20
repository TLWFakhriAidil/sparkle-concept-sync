package handlers

import (
	"database/sql"
	"sparkle-concept-sync/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db        *sql.DB
	jwtSecret string
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"full_name" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func NewAuthHandler(db *sql.DB, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

// Login authenticates a user and returns JWT token
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Find user by email
	var user models.User
	query := `SELECT id, email, full_name, password_hash, is_active FROM users WHERE email = $1`
	err := h.db.QueryRow(query, req.Email).Scan(
		&user.ID, &user.Email, &user.FullName, &user.PasswordHash, &user.IsActive,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid credentials",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Check if user is active
	if !user.IsActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Account is disabled",
		})
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Update last login
	_, err = h.db.Exec("UPDATE users SET last_login = NOW() WHERE id = $1", user.ID)
	if err != nil {
		// Log error but don't fail login
		// logger.Error("Failed to update last login", err)
	}

	// Generate JWT token
	token, err := h.generateJWT(user.ID, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Store session in database
	sessionID := uuid.New().String()
	_, err = h.db.Exec(`
		INSERT INTO user_sessions (id, user_id, token, expires_at) 
		VALUES ($1, $2, $3, $4)`,
		sessionID, user.ID, token, time.Now().Add(24*time.Hour),
	)
	if err != nil {
		// Log error but don't fail login
		// logger.Error("Failed to store session", err)
	}

	// Remove password hash from response
	user.PasswordHash = ""

	return c.JSON(AuthResponse{
		Token: token,
		User:  user,
	})
}

// Register creates a new user account
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Check if user already exists
	var exists bool
	err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email).Scan(&exists)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Email already registered",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Create user
	userID := uuid.New().String()
	query := `
		INSERT INTO users (id, email, full_name, password_hash) 
		VALUES ($1, $2, $3, $4)
		RETURNING created_at, updated_at`

	var createdAt, updatedAt time.Time
	err = h.db.QueryRow(query, userID, req.Email, req.FullName, string(hashedPassword)).Scan(&createdAt, &updatedAt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	// Create user object for response
	user := models.User{
		ID:        userID,
		Email:     req.Email,
		FullName:  req.FullName,
		IsActive:  true,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	// Generate JWT token
	token, err := h.generateJWT(user.ID, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Store session in database
	sessionID := uuid.New().String()
	_, err = h.db.Exec(`
		INSERT INTO user_sessions (id, user_id, token, expires_at) 
		VALUES ($1, $2, $3, $4)`,
		sessionID, user.ID, token, time.Now().Add(24*time.Hour),
	)
	if err != nil {
		// Log error but don't fail registration
		// logger.Error("Failed to store session", err)
	}

	return c.Status(fiber.StatusCreated).JSON(AuthResponse{
		Token: token,
		User:  user,
	})
}

// Logout invalidates the user session
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Get token from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Authorization header required",
		})
	}

	// Extract token (remove "Bearer " prefix)
	token := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	}

	// Delete session from database
	_, err := h.db.Exec("DELETE FROM user_sessions WHERE token = $1", token)
	if err != nil {
		// Log error but don't fail logout
		// logger.Error("Failed to delete session", err)
	}

	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

// JWTMiddleware validates JWT tokens
func (h *AuthHandler) JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header required",
			})
		}

		// Extract token (remove "Bearer " prefix)
		tokenString := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(h.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Extract claims
		claims, ok := token.Claims.(*Claims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		// Check if session exists in database
		var sessionExists bool
		err = h.db.QueryRow(`
			SELECT EXISTS(
				SELECT 1 FROM user_sessions 
				WHERE token = $1 AND user_id = $2 AND expires_at > NOW()
			)`,
			tokenString, claims.UserID,
		).Scan(&sessionExists)

		if err != nil || !sessionExists {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Session expired or invalid",
			})
		}

		// Store user info in context
		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)

		return c.Next()
	}
}

// generateJWT creates a new JWT token
func (h *AuthHandler) generateJWT(userID, email string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "sparkle-concept-sync",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}

// GetCurrentUser returns current authenticated user info
func (h *AuthHandler) GetCurrentUser(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var user models.User
	query := `SELECT id, email, full_name, is_active, created_at, updated_at, last_login FROM users WHERE id = $1`
	err := h.db.QueryRow(query, userID).Scan(
		&user.ID, &user.Email, &user.FullName, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	return c.JSON(user)
}

// CleanupExpiredSessions removes expired sessions (should be called periodically)
func (h *AuthHandler) CleanupExpiredSessions() error {
	_, err := h.db.Exec("DELETE FROM user_sessions WHERE expires_at <= NOW()")
	return err
}
