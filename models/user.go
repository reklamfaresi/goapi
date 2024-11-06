package models

import (
	"database/sql"
	"errors"
	"gogpt/config"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// User yapısı
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Yeni kullanıcı oluşturmak için fonksiyon
func CreateUser(user User) error {
	query := "INSERT INTO users (username, email, password, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := config.DB.Exec(query, user.Username, user.Email, user.Password, user.Role, user.CreatedAt, user.UpdatedAt)
	return err
}

// Kullanıcı adı veya e-posta adresine göre kullanıcının zaten var olup olmadığını kontrol etmek için fonksiyon
func CheckUserExists(username, email string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE username = ? OR email = ?"
	var count int
	err := config.DB.QueryRow(query, username, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Kullanıcı doğrulamak için fonksiyon
func AuthenticateUser(username, password string) (*User, error) {
	var user User
	query := "SELECT id, username, email, password, role, created_at, updated_at FROM users WHERE username = ?"
	err := config.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("kullanıcı bulunamadı")
		}
		return nil, err
	}

	// Şifre doğrulaması
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("şifre hatalı")
	}

	return &user, nil
}

// Tüm kullanıcıları almak için fonksiyon
func GetAllUsers() ([]User, error) {
	var users []User
	rows, err := config.DB.Query("SELECT id, username, email, role, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// Kullanıcı güncellemek için fonksiyon
func UpdateUser(user User) error {
	query := "UPDATE users SET username = ?, email = ?, role = ?, updated_at = ? WHERE id = ?"
	_, err := config.DB.Exec(query, user.Username, user.Email, user.Role, user.UpdatedAt, user.ID)
	return err
}

// Kullanıcı profili güncellemek için fonksiyon
func UpdateUserProfile(username, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query := "UPDATE users SET email = ?, password = ?, updated_at = NOW() WHERE username = ?"
	_, err = config.DB.Exec(query, email, hashedPassword, username)
	return err
}

// Kullanıcı silmek için fonksiyon
func DeleteUser(id int) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := config.DB.Exec(query, id)
	return err
}
