package models

import (
	"gogpt/config"
)

// Ayar yapısını tanımlama
type Setting struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Ayarları veritabanına ekleme veya güncelleme
func SetSetting(key, value string) error {
	query := `
        INSERT INTO settings (setting_key, setting_value)
        VALUES (?, ?)
        ON DUPLICATE KEY UPDATE setting_value = ?
    `
	_, err := config.DB.Exec(query, key, value, value)
	return err
}

// Ayarı veritabanından alma
func GetSetting(key string) (string, error) {
	var value string
	query := "SELECT setting_value FROM settings WHERE setting_key = ?"
	err := config.DB.QueryRow(query, key).Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil
}

// Tüm ayarları listeleme
func GetAllSettings() ([]Setting, error) {
	rows, err := config.DB.Query("SELECT setting_key, setting_value FROM settings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []Setting
	for rows.Next() {
		var setting Setting
		if err := rows.Scan(&setting.Key, &setting.Value); err != nil {
			return nil, err
		}
		settings = append(settings, setting)
	}

	return settings, nil
}
