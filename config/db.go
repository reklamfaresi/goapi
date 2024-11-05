package config

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() {
	var err error
	// Veritabanı bağlantı bilgileri: kullanıcı adı, veritabanı ismi, şifre yok
	DB, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/goweb")
	if err != nil {
		fmt.Println("Veritabanına bağlanırken hata oluştu:", err)
		return
	}
	// Bağlantının başarılı olup olmadığını kontrol ediyoruz
	if err = DB.Ping(); err != nil {
		fmt.Println("Veritabanı bağlantı testi başarısız:", err)
		return
	}
	fmt.Println("Veritabanına başarıyla bağlanıldı!")
}
