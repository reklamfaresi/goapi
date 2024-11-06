package models

import (
	"gogpt/config"
	"time"
)

// AboutUs yapısı
type AboutUs struct {
	ID           int       `json:"id"`
	Section      string    `json:"section"`
	Content      string    `json:"content"`
	LanguageCode string    `json:"language_code"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Tüm About Us bölümlerini almak için fonksiyon
func GetAllAboutUs(languageCode string) ([]AboutUs, error) {
	var aboutUsSections []AboutUs
	query := "SELECT id, section, content, language_code, created_at, updated_at FROM about_us WHERE language_code = ?"
	rows, err := config.DB.Query(query, languageCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var aboutUs AboutUs
		err := rows.Scan(&aboutUs.ID, &aboutUs.Section, &aboutUs.Content, &aboutUs.LanguageCode, &aboutUs.CreatedAt, &aboutUs.UpdatedAt)
		if err != nil {
			return nil, err
		}
		aboutUsSections = append(aboutUsSections, aboutUs)
	}
	return aboutUsSections, nil
}

// Yeni About Us bölümü eklemek için fonksiyon
func CreateAboutUs(aboutUs AboutUs) error {
	query := "INSERT INTO about_us (section, content, language_code, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW())"
	_, err := config.DB.Exec(query, aboutUs.Section, aboutUs.Content, aboutUs.LanguageCode)
	return err
}

// About Us bölümünü güncellemek için fonksiyon
func UpdateAboutUs(aboutUs AboutUs) error {
	query := "UPDATE about_us SET section = ?, content = ?, language_code = ?, updated_at = NOW() WHERE id = ?"
	_, err := config.DB.Exec(query, aboutUs.Section, aboutUs.Content, aboutUs.LanguageCode, aboutUs.ID)
	return err
}

// About Us bölümünü silmek için fonksiyon
func DeleteAboutUs(id int) error {
	query := "DELETE FROM about_us WHERE id = ?"
	_, err := config.DB.Exec(query, id)
	return err
}
