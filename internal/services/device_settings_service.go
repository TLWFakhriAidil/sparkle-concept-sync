package services

import (
	"database/sql"
	"sparkle-concept-sync/internal/models"
)

type DeviceSettingsService struct {
	db *sql.DB
}

func NewDeviceSettingsService(db *sql.DB) *DeviceSettingsService {
	return &DeviceSettingsService{db: db}
}

// GetDevicesByUser returns all devices for a user
func (s *DeviceSettingsService) GetDevicesByUser(userID string) ([]models.DeviceSetting, error) {
	query := `SELECT id, device_id, api_key_option, webhook_id, provider, phone_number, api_key, id_device, user_id, instance, created_at, updated_at FROM device_setting WHERE user_id = $1`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []models.DeviceSetting
	for rows.Next() {
		var device models.DeviceSetting
		err := rows.Scan(
			&device.ID, &device.DeviceID, &device.APIKeyOption, &device.WebhookID,
			&device.Provider, &device.PhoneNumber, &device.APIKey, &device.IDDevice,
			&device.UserID, &device.Instance, &device.CreatedAt, &device.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}

	return devices, nil
}

// GetDeviceByID returns a device by ID
func (s *DeviceSettingsService) GetDeviceByID(id string) (*models.DeviceSetting, error) {
	query := `SELECT id, device_id, api_key_option, webhook_id, provider, phone_number, api_key, id_device, user_id, instance, created_at, updated_at FROM device_setting WHERE id = $1`

	var device models.DeviceSetting
	err := s.db.QueryRow(query, id).Scan(
		&device.ID, &device.DeviceID, &device.APIKeyOption, &device.WebhookID,
		&device.Provider, &device.PhoneNumber, &device.APIKey, &device.IDDevice,
		&device.UserID, &device.Instance, &device.CreatedAt, &device.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &device, nil
}

// CreateDevice creates a new device
func (s *DeviceSettingsService) CreateDevice(device *models.DeviceSetting) error {
	query := `INSERT INTO device_setting (id, device_id, api_key_option, webhook_id, provider, phone_number, api_key, id_device, user_id, instance) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := s.db.Exec(query, device.ID, device.DeviceID, device.APIKeyOption, device.WebhookID, device.Provider, device.PhoneNumber, device.APIKey, device.IDDevice, device.UserID, device.Instance)
	return err
}

// UpdateDevice updates an existing device
func (s *DeviceSettingsService) UpdateDevice(device *models.DeviceSetting) error {
	query := `UPDATE device_setting SET device_id = $2, api_key_option = $3, webhook_id = $4, provider = $5, phone_number = $6, api_key = $7, id_device = $8, instance = $9, updated_at = NOW() WHERE id = $1`

	_, err := s.db.Exec(query, device.ID, device.DeviceID, device.APIKeyOption, device.WebhookID, device.Provider, device.PhoneNumber, device.APIKey, device.IDDevice, device.Instance)
	return err
}

// DeleteDevice deletes a device
func (s *DeviceSettingsService) DeleteDevice(id string) error {
	query := `DELETE FROM device_setting WHERE id = $1`
	_, err := s.db.Exec(query, id)
	return err
}
