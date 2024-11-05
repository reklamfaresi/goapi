package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gogpt/config"
	"gogpt/models" // Aynı şekilde modül adını kullanarak models paketine erişiyoruz
	"log"
	"net/http"
)

func main() {
	// Veritabanına bağlan
	config.Connect()

	// Router oluştur
	router := mux.NewRouter()

	// Kullanıcıları listelemek için rota
	router.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
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
