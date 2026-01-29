package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	if host == "" || user == "" || pass == "" || name == "" || port == "" {
		log.Fatal("❌ Database env vars missing. Check .env loading")
	}

	fmt.Println("DB_HOST:", host)
	fmt.Println("DB_PORT:", port)
	fmt.Println("DB_PASSWORD length:", len(pass))

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
		host, user, pass, name, port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect db:", err)
	}

	DB = db
	fmt.Println("✅ Database Connected Successfully")
}
