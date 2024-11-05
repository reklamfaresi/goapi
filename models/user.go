package models

import (
	"gogpt/config"
)

// Kullanıcı yapısı (model)
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
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

	return users, nil
}
