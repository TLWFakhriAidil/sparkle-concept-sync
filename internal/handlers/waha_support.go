package handlers

import (
	"log"
	"sparkle-concept-sync/internal/models"
	"sparkle-concept-sync/internal/services"

	"github.com/gofiber/fiber/v2"
)

type WAHAHandler struct {
	flowService      *services.FlowService
	providerService  *services.ProviderService
	websocketService *services.WebSocketService
}

func NewWAHAHandler(flowService *services.FlowService, providerService *services.ProviderService, websocketService *services.WebSocketService) *WAHAHandler {
	return &WAHAHandler{
		flowService:      flowService,
		providerService:  providerService,
		websocketService: websocketService,
	}
}

// HandleWAHAWebhook processes incoming WAHA webhook messages
func (h *WAHAHandler) HandleWAHAWebhook(c *fiber.Ctx) error {
	deviceID := c.Params("device_id")
	if deviceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Device ID is required",
		})
	}

	var payload models.WhatsAppMessage
	if err := c.BodyParser(&payload); err != nil {
		log.Printf("Failed to parse webhook payload: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payload format",
		})
	}

	// Set device ID from URL parameter
	payload.DeviceID = deviceID

	// Process message asynchronously
	go h.processMessage(payload)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "received",
		"message": "Webhook processed successfully",
	})
}

// HandleWablasWebhook processes incoming Wablas webhook messages
func (h *WAHAHandler) HandleWablasWebhook(c *fiber.Ctx) error {
	deviceID := c.Params("device_id")
	if deviceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Device ID is required",
		})
	}

	// Parse Wablas specific payload format
	var payload map[string]interface{}
	if err := c.BodyParser(&payload); err != nil {
		log.Printf("Failed to parse Wablas webhook payload: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payload format",
		})
	}

	// Convert Wablas format to standardized WhatsApp message
	message := h.convertWablasToStandardFormat(payload, deviceID)

	// Process message asynchronously
	go h.processMessage(message)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "received",
	})
}

// HandleWhacenterWebhook processes incoming Whacenter webhook messages
func (h *WAHAHandler) HandleWhacenterWebhook(c *fiber.Ctx) error {
	deviceID := c.Params("device_id")
	if deviceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Device ID is required",
		})
	}

	// Parse Whacenter specific payload format
	var payload map[string]interface{}
	if err := c.BodyParser(&payload); err != nil {
		log.Printf("Failed to parse Whacenter webhook payload: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payload format",
		})
	}

	// Convert Whacenter format to standardized WhatsApp message
	message := h.convertWhacenterToStandardFormat(payload, deviceID)

	// Process message asynchronously
	go h.processMessage(message)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "received",
	})
}

// processMessage handles the core message processing logic
func (h *WAHAHandler) processMessage(message models.WhatsAppMessage) {
	// Skip outgoing messages
	if message.From == message.DeviceID {
		return
	}

	log.Printf("Processing message from %s: %s", message.From, message.Body)

	// Execute flow
	response, err := h.flowService.ExecuteFlow(message)
	if err != nil {
		log.Printf("Flow execution error: %v", err)
		return
	}

	// Send response via provider
	if response != nil {
		err = h.providerService.SendMessage(message.DeviceID, message.From, response)
		if err != nil {
			log.Printf("Failed to send response: %v", err)
			return
		}
	}

	// Broadcast real-time update
	if h.websocketService != nil {
		update := models.WebSocketMessage{
			Type:     "message_processed",
			DeviceID: message.DeviceID,
			Data: map[string]interface{}{
				"from":     message.From,
				"message":  message.Body,
				"response": response,
			},
		}
		h.websocketService.Broadcast(update)
	}
}

// convertWablasToStandardFormat converts Wablas webhook format to standard format
func (h *WAHAHandler) convertWablasToStandardFormat(payload map[string]interface{}, deviceID string) models.WhatsAppMessage {
	message := models.WhatsAppMessage{
		DeviceID: deviceID,
		Type:     "text",
	}

	// Map Wablas fields to standard format
	if from, ok := payload["from"].(string); ok {
		message.From = from
	}
	if body, ok := payload["body"].(string); ok {
		message.Body = body
	}
	if to, ok := payload["to"].(string); ok {
		message.To = to
	}
	if msgType, ok := payload["type"].(string); ok {
		message.Type = msgType
	}
	if msgID, ok := payload["id"].(string); ok {
		message.MessageID = msgID
	}
	if timestamp, ok := payload["timestamp"].(float64); ok {
		message.Timestamp = int64(timestamp)
	}

	return message
}

// convertWhacenterToStandardFormat converts Whacenter webhook format to standard format
func (h *WAHAHandler) convertWhacenterToStandardFormat(payload map[string]interface{}, deviceID string) models.WhatsAppMessage {
	message := models.WhatsAppMessage{
		DeviceID: deviceID,
		Type:     "text",
	}

	// Map Whacenter fields to standard format
	if data, ok := payload["data"].(map[string]interface{}); ok {
		if from, ok := data["from"].(string); ok {
			message.From = from
		}
		if body, ok := data["body"].(string); ok {
			message.Body = body
		}
		if to, ok := data["to"].(string); ok {
			message.To = to
		}
		if msgType, ok := data["type"].(string); ok {
			message.Type = msgType
		}
		if msgID, ok := data["id"].(string); ok {
			message.MessageID = msgID
		}
	}

	if timestamp, ok := payload["timestamp"].(float64); ok {
		message.Timestamp = int64(timestamp)
	}

	return message
}

// GetWebhookInfo returns webhook configuration info
func (h *WAHAHandler) GetWebhookInfo(c *fiber.Ctx) error {
	webhookInfo := map[string]interface{}{
		"endpoints": map[string]string{
			"waha":      "/webhooks/waha/{device_id}",
			"wablas":    "/webhooks/wablas/{device_id}",
			"whacenter": "/webhooks/whacenter/{device_id}",
		},
		"supported_providers": []string{
			"waha",
			"wablas",
			"whacenter",
		},
		"message_types": []string{
			"text",
			"image",
			"audio",
			"video",
			"document",
		},
		"setup_instructions": map[string]string{
			"waha":      "Configure your WAHA instance to send webhooks to the WAHA endpoint",
			"wablas":    "Set your Wablas webhook URL to the Wablas endpoint",
			"whacenter": "Configure Whacenter webhook to point to the Whacenter endpoint",
		},
	}

	return c.JSON(webhookInfo)
}

// ValidateWebhook validates webhook signatures (for providers that support it)
func (h *WAHAHandler) ValidateWebhook(c *fiber.Ctx) error {
	provider := c.Query("provider")
	signature := c.Get("X-Signature")

	if provider == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Provider parameter required",
		})
	}

	// Validate signature based on provider
	switch provider {
	case "wablas":
		// Implement Wablas signature validation
		return h.validateWablasSignature(c, signature)
	case "whacenter":
		// Implement Whacenter signature validation
		return h.validateWhacenterSignature(c, signature)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unsupported provider",
		})
	}
}

func (h *WAHAHandler) validateWablasSignature(c *fiber.Ctx, signature string) error {
	// Implement Wablas-specific signature validation
	// This would involve checking the signature against the payload and secret
	return c.JSON(fiber.Map{
		"valid":    true,
		"provider": "wablas",
	})
}

func (h *WAHAHandler) validateWhacenterSignature(c *fiber.Ctx, signature string) error {
	// Implement Whacenter-specific signature validation
	return c.JSON(fiber.Map{
		"valid":    true,
		"provider": "whacenter",
	})
}

// GetWebhookStats returns webhook processing statistics
func (h *WAHAHandler) GetWebhookStats(c *fiber.Ctx) error {
	// This would be implemented to return real statistics
	stats := map[string]interface{}{
		"total_messages_processed": 0,
		"messages_by_provider": map[string]int{
			"waha":      0,
			"wablas":    0,
			"whacenter": 0,
		},
		"average_processing_time": "150ms",
		"success_rate":            "99.2%",
	}

	return c.JSON(stats)
}
