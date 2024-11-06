package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gogpt/config"
	"gogpt/models"
	"gogpt/utils"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"time"
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

		token, err := utils.GenerateJWT(user.Username, user.Role)
		if err != nil {
			http.Error(w, "Token oluşturulamadı", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}).Methods("POST")

	// Kullanıcı eklemek için rota
	router.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var newUser models.User
			err := json.NewDecoder(r.Body).Decode(&newUser)
			if err != nil {
				http.Error(w, "Geçersiz veri formatı", http.StatusBadRequest)
				return
			}

			// Kullanıcının zaten var olup olmadığını kontrol et
			exists, err := models.CheckUserExists(newUser.Username, newUser.Email)
			if err != nil {
				http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
				return
			}
			if exists {
				http.Error(w, "Bu kullanıcı adı veya e-posta zaten kayıtlı", http.StatusConflict)
				return
			}

			// Varsayılan rol atanması
			if newUser.Role == "" {
				newUser.Role = "user"
			}

			// Şifreyi hashle ve kullanıcıyı oluştur
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "Şifre hashlenemedi", http.StatusInternalServerError)
				return
			}
			newUser.Password = string(hashedPassword)

			// Oluşturma tarihi ekleyin
			newUser.CreatedAt = time.Now()
			newUser.UpdatedAt = time.Now()

			err = models.CreateUser(newUser)
			if err != nil {
				http.Error(w, "Kullanıcı oluşturulamadı", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{"message": "Kullanıcı başarıyla oluşturuldu"})
		} else {
			http.Error(w, "Yöntem desteklenmiyor", http.StatusMethodNotAllowed)
		}
	}).Methods("POST")

	// Kullanıcıları listelemek için rota (JWT ve dinamik rol kontrolü ile korumalı)
	router.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Yetkisiz erişim", http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[len("Bearer "):len(authHeader)]

		role, err := utils.GetUserRoleFromJWT(tokenString)
		if err != nil {
			http.Error(w, "Geçersiz token", http.StatusUnauthorized)
			return
		}

		// Kullanıcının izin kontrolü
		hasPermission, err := models.CheckPermission(role, "list_users")
		if err != nil || !hasPermission {
			http.Error(w, "Yetkisiz erişim", http.StatusForbidden)
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

	// Kullanıcı güncelleme ve silme işlemleri
	router.HandleFunc("/api/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Yetkisiz erişim", http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[len("Bearer "):len(authHeader)]
		role, err := utils.GetUserRoleFromJWT(tokenString)
		if err != nil {
			http.Error(w, "Geçersiz token", http.StatusUnauthorized)
			return
		}

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Geçersiz kullanıcı ID'si", http.StatusBadRequest)
			return
		}

		// Kullanıcıyı yalnızca adminler ya da kendi bilgileri üzerinde işlem yapabilir
		if role != "admin" {
			http.Error(w, "Yetkisiz erişim", http.StatusForbidden)
			return
		}

		if r.Method == http.MethodPut {
			var updatedUser models.User
			err = json.NewDecoder(r.Body).Decode(&updatedUser)
			if err != nil {
				http.Error(w, "Geçersiz veri formatı", http.StatusBadRequest)
				return
			}

			updatedUser.ID = id
			updatedUser.UpdatedAt = time.Now() // Güncellenme tarihi

			err = models.UpdateUser(updatedUser)
			if err != nil {
				http.Error(w, "Kullanıcı güncellenemedi", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Kullanıcı başarıyla güncellendi"})
		} else if r.Method == http.MethodDelete {
			err = models.DeleteUser(id)
			if err != nil {
				http.Error(w, "Kullanıcı silinemedi", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Kullanıcı başarıyla silindi"})
		}
	}).Methods("PUT", "DELETE")

	// Kullanıcı profili güncelleme (JWT ile korumalı)
	router.HandleFunc("/api/profile", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Yetkisiz erişim", http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[len("Bearer "):len(authHeader)]

		// Kullanıcı adı JWT'den elde ediliyor
		username, err := utils.GetUsernameFromJWT(tokenString)
		if err != nil {
			http.Error(w, "Geçersiz token", http.StatusUnauthorized)
			return
		}

		if r.Method == http.MethodPut {
			var updatedProfile struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}

			err = json.NewDecoder(r.Body).Decode(&updatedProfile)
			if err != nil {
				http.Error(w, "Geçersiz veri formatı", http.StatusBadRequest)
				return
			}

			// Şifreyi hashleyelim
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedProfile.Password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "Şifre hashlenemedi", http.StatusInternalServerError)
				return
			}

			// Kullanıcıyı güncelle
			err = models.UpdateUserProfile(username, updatedProfile.Email, string(hashedPassword))
			if err != nil {
				http.Error(w, "Profil güncellenemedi", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Profil başarıyla güncellendi"})
		} else {
			http.Error(w, "Yöntem desteklenmiyor", http.StatusMethodNotAllowed)
		}
	}).Methods("PUT")

	// Entegrasyon ayarlarını almak ve güncellemek için rota
	router.HandleFunc("/api/integrations", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			integration, err := models.GetIntegrations()
			if err != nil {
				http.Error(w, "Entegrasyon ayarları alınamadı", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(integration)
		} else if r.Method == http.MethodPut {
			var updatedIntegration models.Integration
			err := json.NewDecoder(r.Body).Decode(&updatedIntegration)
			if err != nil {
				http.Error(w, "Geçersiz veri formatı", http.StatusBadRequest)
				return
			}

			updatedIntegration.UpdatedAt = time.Now() // Güncellenme tarihi

			err = models.UpdateIntegrations(updatedIntegration)
			if err != nil {
				http.Error(w, "Entegrasyon ayarları güncellenemedi", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Entegrasyon ayarları başarıyla güncellendi"})
		} else {
			http.Error(w, "Yöntem desteklenmiyor", http.StatusMethodNotAllowed)
		}
	}).Methods("GET", "PUT")

	// Slider CRUD işlemleri
	router.HandleFunc("/api/sliders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			sliders, err := models.GetAllSliders()
			if err != nil {
				http.Error(w, "Sliderlar alınamadı", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(sliders)
		} else if r.Method == http.MethodPost {
			var newSlider models.Slider
			err := json.NewDecoder(r.Body).Decode(&newSlider)
			if err != nil {
				http.Error(w, "Geçersiz veri formatı", http.StatusBadRequest)
				return
			}

			newSlider.CreatedAt = time.Now()
			newSlider.UpdatedAt = time.Now()

			err = models.CreateSlider(newSlider)
			if err != nil {
				http.Error(w, "Slider oluşturulamadı", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{"message": "Slider başarıyla oluşturuldu"})
		} else {
			http.Error(w, "Yöntem desteklenmiyor", http.StatusMethodNotAllowed)
		}
	}).Methods("GET", "POST")

	router.HandleFunc("/api/sliders/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Geçersiz slider ID'si", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodPut {
			var updatedSlider models.Slider
			err := json.NewDecoder(r.Body).Decode(&updatedSlider)
			if err != nil {
				http.Error(w, "Geçersiz veri formatı", http.StatusBadRequest)
				return
			}

			updatedSlider.ID = id
			updatedSlider.UpdatedAt = time.Now() // Güncellenme tarihi

			err = models.UpdateSlider(updatedSlider)
			if err != nil {
				http.Error(w, "Slider güncellenemedi", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Slider başarıyla güncellendi"})
		} else if r.Method == http.MethodDelete {
			err = models.DeleteSlider(id)
			if err != nil {
				http.Error(w, "Slider silinemedi", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Slider başarıyla silindi"})
		} else {
			http.Error(w, "Yöntem desteklenmiyor", http.StatusMethodNotAllowed)
		}
	}).Methods("PUT", "DELETE")

	// About Us CRUD işlemleri
	router.HandleFunc("/api/about-us", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			languageCode := r.URL.Query().Get("lang")
			if languageCode == "" {
				languageCode = "en" // Varsayılan dil İngilizce
			}
			aboutUsSections, err := models.GetAllAboutUs(languageCode)
			if err != nil {
				http.Error(w, "About Us bölümleri alınamadı", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(aboutUsSections)
		} else if r.Method == http.MethodPost {
			var newAboutUs models.AboutUs
			err := json.NewDecoder(r.Body).Decode(&newAboutUs)
			if err != nil {
				http.Error(w, "Geçersiz veri formatı", http.StatusBadRequest)
				return
			}
			err = models.CreateAboutUs(newAboutUs)
			if err != nil {
				http.Error(w, "About Us bölümü oluşturulamadı", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{"message": "About Us bölümü başarıyla oluşturuldu"})
		} else {
			http.Error(w, "Yöntem desteklenmiyor", http.StatusMethodNotAllowed)
		}
	}).Methods("GET", "POST")

	router.HandleFunc("/api/about-us/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Geçersiz About Us ID'si", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodPut {
			var updatedAboutUs models.AboutUs
			err := json.NewDecoder(r.Body).Decode(&updatedAboutUs)
			if err != nil {
				http.Error(w, "Geçersiz veri formatı", http.StatusBadRequest)
				return
			}
			updatedAboutUs.ID = id
			err = models.UpdateAboutUs(updatedAboutUs)
			if err != nil {
				http.Error(w, "About Us bölümü güncellenemedi", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "About Us bölümü başarıyla güncellendi"})
		} else if r.Method == http.MethodDelete {
			err = models.DeleteAboutUs(id)
			if err != nil {
				http.Error(w, "About Us bölümü silinemedi", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "About Us bölümü başarıyla silindi"})
		} else {
			http.Error(w, "Yöntem desteklenmiyor", http.StatusMethodNotAllowed)
		}
	}).Methods("PUT", "DELETE")

	// Sunucuyu başlat
	log.Println("Sunucu başlatılıyor: http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
