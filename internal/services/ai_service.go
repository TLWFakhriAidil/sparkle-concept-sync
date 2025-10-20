package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"sparkle-concept-sync/internal/models"
)

type AIService struct {
	openRouterAPIKey string
	redisService     *RedisService
	httpClient       *http.Client
}

type OpenRouterRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenRouterResponse struct {
	Choices []Choice  `json:"choices"`
	Error   *APIError `json:"error,omitempty"`
}

type Choice struct {
	Message Message `json:"message"`
}

type APIError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

const (
	openRouterBaseURL = "https://openrouter.ai/api/v1/chat/completions"
	cacheTimeout      = 5 * time.Minute
)

func NewAIService(openRouterAPIKey string, redisService *RedisService) *AIService {
	return &AIService{
		openRouterAPIKey: openRouterAPIKey,
		redisService:     redisService,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// GetAIResponse generates AI response with caching and rate limiting
func (s *AIService) GetAIResponse(ctx context.Context, prompt, model, userID string) (*models.AIResponse, error) {
	// Create cache key
	cacheKey := fmt.Sprintf("ai_response:%s:%s:%s", userID, model, prompt)

	// Try to get from cache first
	if cached, err := s.redisService.Get(ctx, cacheKey); err == nil && cached != "" {
		var response models.AIResponse
		if err := json.Unmarshal([]byte(cached), &response); err == nil {
			return &response, nil
		}
	}

	// Check rate limit
	rateLimitKey := fmt.Sprintf("rate_limit:ai:%s", userID)
	if !s.redisService.CheckRateLimit(ctx, rateLimitKey, 100, time.Minute) {
		return nil, fmt.Errorf("rate limit exceeded for user %s", userID)
	}

	// Make API request
	response, err := s.makeOpenRouterRequest(ctx, prompt, model)
	if err != nil {
		return nil, err
	}

	// Parse AI response
	aiResponse, err := s.parseAIResponse(response)
	if err != nil {
		return nil, err
	}

	// Cache the response
	if responseBytes, err := json.Marshal(aiResponse); err == nil {
		s.redisService.Set(ctx, cacheKey, string(responseBytes), cacheTimeout)
	}

	return aiResponse, nil
}

func (s *AIService) makeOpenRouterRequest(ctx context.Context, prompt, model string) (string, error) {
	// Default to GPT-4 if model not specified
	if model == "" {
		model = "openai/gpt-4"
	}

	// Construct system prompt for chatbot context
	systemPrompt := `You are an AI assistant for a WhatsApp chatbot. You must respond in this exact JSON format:
{
  "Stage": "Current conversation stage",
  "Response": [
    {"type": "text", "content": "Your response message here"}
  ]
}

Available response types: text, image, audio, video, delay, condition
Keep responses conversational and helpful. Always include a Stage and Response array.`

	request := OpenRouterRequest{
		Model: model,
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: prompt},
		},
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", openRouterBaseURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.openRouterAPIKey)
	req.Header.Set("HTTP-Referer", "https://sparkle-concept-sync.com")
	req.Header.Set("X-Title", "Sparkle Concept Sync")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var openRouterResponse OpenRouterResponse
	if err := json.Unmarshal(body, &openRouterResponse); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if openRouterResponse.Error != nil {
		return "", fmt.Errorf("API error: %s", openRouterResponse.Error.Message)
	}

	if len(openRouterResponse.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return openRouterResponse.Choices[0].Message.Content, nil
}

func (s *AIService) parseAIResponse(response string) (*models.AIResponse, error) {
	var aiResponse models.AIResponse

	// Try to parse as JSON first
	if err := json.Unmarshal([]byte(response), &aiResponse); err == nil {
		// Validate that we have required fields
		if aiResponse.Stage != "" && len(aiResponse.Response) > 0 {
			return &aiResponse, nil
		}
	}

	// If JSON parsing fails or invalid structure, create a fallback response
	aiResponse = models.AIResponse{
		Stage: "General Response",
		Response: []models.AIMessage{
			{
				Type:    "text",
				Content: response,
			},
		},
	}

	return &aiResponse, nil
}

// GetAvailableModels returns list of supported AI models
func (s *AIService) GetAvailableModels() []string {
	return []string{
		"openai/gpt-5-chat",
		"openai/gpt-5-mini",
		"openai/chatgpt-4o-latest",
		"openai/gpt-4.1",
		"google/gemini-2.5-pro",
		"google/gemini-pro-1.5",
	}
}

// ValidateModel checks if a model is supported
func (s *AIService) ValidateModel(model string) bool {
	supportedModels := s.GetAvailableModels()
	for _, supportedModel := range supportedModels {
		if model == supportedModel {
			return true
		}
	}
	return false
}

// ProcessFlowPrompt processes a chatbot flow prompt with context
func (s *AIService) ProcessFlowPrompt(ctx context.Context, prompt, model, userID string, flowContext map[string]interface{}) (*models.AIResponse, error) {
	// Enhance prompt with flow context
	enhancedPrompt := s.buildContextualPrompt(prompt, flowContext)

	return s.GetAIResponse(ctx, enhancedPrompt, model, userID)
}

func (s *AIService) buildContextualPrompt(prompt string, context map[string]interface{}) string {
	contextStr := ""

	if context != nil {
		if stage, ok := context["stage"].(string); ok && stage != "" {
			contextStr += fmt.Sprintf("Current Stage: %s\n", stage)
		}

		if userName, ok := context["user_name"].(string); ok && userName != "" {
			contextStr += fmt.Sprintf("User Name: %s\n", userName)
		}

		if previousMessages, ok := context["previous_messages"].(string); ok && previousMessages != "" {
			contextStr += fmt.Sprintf("Previous Context: %s\n", previousMessages)
		}

		if flowData, ok := context["flow_data"].(map[string]interface{}); ok {
			if niche, ok := flowData["niche"].(string); ok && niche != "" {
				contextStr += fmt.Sprintf("Business Niche: %s\n", niche)
			}
		}
	}

	if contextStr != "" {
		return fmt.Sprintf("%s\n\nUser Message: %s", contextStr, prompt)
	}

	return prompt
}

// Circuit breaker for AI service reliability
type CircuitBreaker struct {
	failureCount    int
	lastFailureTime time.Time
	timeout         time.Duration
	threshold       int
}

func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		threshold: threshold,
		timeout:   timeout,
	}
}

func (cb *CircuitBreaker) Call(operation func() error) error {
	if cb.failureCount >= cb.threshold {
		if time.Since(cb.lastFailureTime) < cb.timeout {
			return fmt.Errorf("circuit breaker is open")
		}
		cb.failureCount = 0
	}

	err := operation()
	if err != nil {
		cb.failureCount++
		cb.lastFailureTime = time.Now()
		return err
	}

	cb.failureCount = 0
	return nil
}
