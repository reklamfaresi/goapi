package models

import (
	"gogpt/config"
	"time"
)

// Integration yapısı
type Integration struct {
	ID                int       `json:"id"`
	GoogleAnalyticsID string    `json:"google_analytics_id"`
	FacebookPixelID   string    `json:"facebook_pixel_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// Entegrasyon ayarlarını almak için fonksiyon
func GetIntegrations() (*Integration, error) {
	var integration Integration
	query := "SELECT id, google_analytics_id, facebook_pixel_id, created_at, updated_at FROM integrations WHERE id = 1"
	err := config.DB.QueryRow(query).Scan(&integration.ID, &integration.GoogleAnalyticsID, &integration.FacebookPixelID, &integration.CreatedAt, &integration.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &integration, nil
}

// Entegrasyon ayarlarını güncellemek için fonksiyon
func UpdateIntegrations(integration Integration) error {
	query := "UPDATE integrations SET google_analytics_id = ?, facebook_pixel_id = ?, updated_at = NOW() WHERE id = ?"
	_, err := config.DB.Exec(query, integration.GoogleAnalyticsID, integration.FacebookPixelID, integration.ID)
	return err
}
