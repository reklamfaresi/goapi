package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Router oluşturuluyor
	router := mux.NewRouter()

	// Basit bir test rotası ekleyelim
	router.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Admin panel API çalışıyor!"))
	})

	// Sunucu başlatılıyor
	log.Println("Sunucu 8000 portunda dinleniyor...")
	log.Fatal(http.ListenAndServe(":8000", router))
}
