package models

import (
	"gogpt/config"
	"log"
)

// Kullanıcı yapısı (model)
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// Kullanıcıları veritabanından çekmek için bir fonksiyon
func GetAllUsers() ([]User, error) {
	rows, err := config.DB.Query("SELECT id, username, email, role FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Terminale kullanıcıların sayısını yazdırmak için
	log.Printf("Toplam %d kullanıcı bulundu.\n", len(users))

	return users, nil
}

// Yeni bir kullanıcı eklemek için fonksiyon
func CreateUser(user User) error {
	query := "INSERT INTO users (username, password, email, role) VALUES (?, ?, ?, ?)"
	_, err := config.DB.Exec(query, user.Username, user.Password, user.Email, user.Role)
	if err != nil {
		return err
	}
	return nil
}

// Kullanıcıyı güncellemek için fonksiyon
func UpdateUser(user User) error {
	query := "UPDATE users SET username = ?, email = ?, role = ? WHERE id = ?"
	_, err := config.DB.Exec(query, user.Username, user.Email, user.Role, user.ID)
	if err != nil {
		return err
	}
	return nil
}

// Kullanıcıyı silmek için fonksiyon
func DeleteUser(id int) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := config.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

// Kullanıcı kimlik doğrulama fonksiyonu
func AuthenticateUser(username, password string) (*User, error) {
	var user User
	query := "SELECT id, username, email, role FROM users WHERE username = ? AND password = ?"
	err := config.DB.QueryRow(query, username, password).Scan(&user.ID, &user.Username, &user.Email, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
