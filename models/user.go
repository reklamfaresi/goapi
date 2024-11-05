package models

import (
	"fmt"
	"gogpt/config"
	"golang.org/x/crypto/bcrypt"
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

	return users, nil
}

// Yeni bir kullanıcı eklemek için fonksiyon (şifre hashleme ile)
func CreateUser(user User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := "INSERT INTO users (username, password, email, role) VALUES (?, ?, ?, ?)"
	_, err = config.DB.Exec(query, user.Username, string(hashedPassword), user.Email, user.Role)
	if err != nil {
		return err
	}
	return nil
}

// Kullanıcı kimlik doğrulama fonksiyonu (bcrypt kullanarak)
func AuthenticateUser(username, password string) (*User, error) {
	var user User
	query := "SELECT id, username, password, email, role FROM users WHERE username = ?"
	err := config.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Role)
	if err != nil {
		return nil, err
	}

	// Kullanıcının verdiği şifre ile veritabanındaki hash'i karşılaştır
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("geçersiz şifre")
	}

	return &user, nil
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

// Kullanıcı profilini güncellemek için fonksiyon
func UpdateUserProfile(username, email, hashedPassword string) error {
	query := "UPDATE users SET email = ?, password = ? WHERE username = ?"
	_, err := config.DB.Exec(query, email, hashedPassword, username)
	if err != nil {
		return err
	}
	return nil
}

// Kullanıcının mevcut olup olmadığını kontrol etmek için fonksiyon
func CheckUserExists(username, email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username = ? OR email = ?)"
	err := config.DB.QueryRow(query, username, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
