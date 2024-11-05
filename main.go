package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gogpt/config"
	"gogpt/models"
	"gogpt/utils"
	"log"
	"net/http"
)

func main() {
	// Veritabanına bağlan
	config.Connect()

	// Router oluştur
	router := mux.NewRouter()

	// Kullanıcı giriş işlemi
	router.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		var loginDetails struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&loginDetails)
		if err != nil {
			http.Error(w, "Geçersiz veri formatı", http.StatusBadRequest)
			return
		}

		user, err := models.AuthenticateUser(loginDetails.Username, loginDetails.Password)
		if err != nil {
			http.Error(w, "Kullanıcı doğrulanamadı", http.StatusUnauthorized)
			return
		}

		token, err := utils.GenerateJWT(user.Username)
		if err != nil {
			http.Error(w, "Token oluşturulamadı", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}).Methods("POST")

	// Diğer endpoint'ler (GET, POST, PUT, DELETE kullanıcılar)
	// ...
	// Kullanıcıları listelemek için rota (JWT korumalı)
	router.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Yetkisiz erişim", http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[len("Bearer "):]

		_, err := utils.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Geçersiz token", http.StatusUnauthorized)
			return
		}

		users, err := models.GetAllUsers()
		if err != nil {
			http.Error(w, "Kullanıcılar alınamadı", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}).Methods("GET")

	// Sunucuyu başlat
	log.Println("Sunucu 8000 portunda dinleniyor...")
	log.Fatal(http.ListenAndServe(":8000", router))
}
