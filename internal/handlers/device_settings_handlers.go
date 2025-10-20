package handlers

import (
	"sparkle-concept-sync/internal/models"
	"sparkle-concept-sync/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type DeviceSettingsHandler struct {
	service *services.DeviceSettingsService
}

func NewDeviceSettingsHandler(service *services.DeviceSettingsService) *DeviceSettingsHandler {
	return &DeviceSettingsHandler{service: service}
}

// GetDevices returns all devices for the authenticated user
func (h *DeviceSettingsHandler) GetDevices(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	devices, err := h.service.GetDevicesByUser(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch devices",
		})
	}

	return c.JSON(devices)
}

// GetDevice returns a specific device by ID
func (h *DeviceSettingsHandler) GetDevice(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)

	device, err := h.service.GetDeviceByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Device not found",
		})
	}

	// Check ownership
	if device.UserID == nil || *device.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
		})
	}

	return c.JSON(device)
}

// CreateDevice creates a new device
func (h *DeviceSettingsHandler) CreateDevice(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req models.DeviceSetting
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Generate ID and set user
	req.ID = uuid.New().String()
	req.UserID = &userID

	if err := h.service.CreateDevice(&req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create device",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(req)
}

// UpdateDevice updates an existing device
func (h *DeviceSettingsHandler) UpdateDevice(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)

	// Check if device exists and user owns it
	existingDevice, err := h.service.GetDeviceByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Device not found",
		})
	}

	if existingDevice.UserID == nil || *existingDevice.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
		})
	}

	var req models.DeviceSetting
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Preserve ID and user
	req.ID = id
	req.UserID = &userID

	if err := h.service.UpdateDevice(&req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update device",
		})
	}

	return c.JSON(req)
}

// DeleteDevice deletes a device
func (h *DeviceSettingsHandler) DeleteDevice(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)

	// Check if device exists and user owns it
	existingDevice, err := h.service.GetDeviceByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Device not found",
		})
	}

	if existingDevice.UserID == nil || *existingDevice.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
		})
	}

	if err := h.service.DeleteDevice(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete device",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Device deleted successfully",
	})
}
