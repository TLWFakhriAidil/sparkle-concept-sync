package handlers

import (
	"database/sql"
	"sparkle-concept-sync/internal/models"

	"github.com/gofiber/fiber/v2"
)

type ProfileHandler struct {
	db *sql.DB
}

func NewProfileHandler(db *sql.DB) *ProfileHandler {
	return &ProfileHandler{db: db}
}

// GetProfile returns the current user's profile
func (h *ProfileHandler) GetProfile(c *fiber.Ctx) error {
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

// UpdateProfile updates the current user's profile
func (h *ProfileHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update user profile
	query := `UPDATE users SET full_name = $2, email = $3, updated_at = NOW() WHERE id = $1`
	_, err := h.db.Exec(query, userID, req.FullName, req.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update profile",
		})
	}

	// Return updated profile
	return h.GetProfile(c)
}
