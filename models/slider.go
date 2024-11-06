package models

import (
	"gogpt/config"
	"time"
)

// Slider yapısı
type Slider struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	ImageURL     string    `json:"image_url"`
	LanguageCode string    `json:"language_code"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Tüm slider öğelerini almak için fonksiyon
func GetAllSliders() ([]Slider, error) {
	var sliders []Slider
	rows, err := config.DB.Query("SELECT id, title, description, image_url, language_code, is_active, created_at, updated_at FROM slider")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var slider Slider
		err := rows.Scan(&slider.ID, &slider.Title, &slider.Description, &slider.ImageURL, &slider.LanguageCode, &slider.IsActive, &slider.CreatedAt, &slider.UpdatedAt)
		if err != nil {
			return nil, err
		}
		sliders = append(sliders, slider)
	}
	return sliders, nil
}

// Yeni bir slider öğesi eklemek için fonksiyon
func CreateSlider(slider Slider) error {
	query := "INSERT INTO slider (title, description, image_url, language_code, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, NOW(), NOW())"
	_, err := config.DB.Exec(query, slider.Title, slider.Description, slider.ImageURL, slider.LanguageCode, slider.IsActive)
	return err
}

// Slider öğesini güncellemek için fonksiyon
func UpdateSlider(slider Slider) error {
	query := "UPDATE slider SET title = ?, description = ?, image_url = ?, language_code = ?, is_active = ?, updated_at = NOW() WHERE id = ?"
	_, err := config.DB.Exec(query, slider.Title, slider.Description, slider.ImageURL, slider.LanguageCode, slider.IsActive, slider.ID)
	return err
}

// Slider öğesini silmek için fonksiyon
func DeleteSlider(id int) error {
	query := "DELETE FROM slider WHERE id = ?"
	_, err := config.DB.Exec(query, id)
	return err
}
