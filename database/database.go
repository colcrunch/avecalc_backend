package database

import (
	"log"
	"os"
	"strings"

	"github.com/colcrunch/avecalc_backend/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

func buildDsn() string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	b := strings.Builder{}

	b.Write([]byte(os.Getenv("DB_USER")))
	b.Write([]byte(":"))
	b.Write([]byte(os.Getenv("DB_PASS")))
	b.Write([]byte("@tcp("))
	b.Write([]byte(os.Getenv("DB_HOST")))
	b.Write([]byte(":"))
	b.Write([]byte(os.Getenv("DB_PORT")))
	b.Write([]byte(")/"))
	b.Write([]byte(os.Getenv("DB_NAME")))
	b.Write([]byte("?charset=utf8mb4&parseTime=True&loc=Local"))

	return b.String()

}

func ConnectDB() {

	dsn := buildDsn()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Could not connect to the database! \n", err.Error())
		os.Exit(1)
	}

	db.Logger = logger.Default.LogMode(logger.Info)
	db.AutoMigrate(&models.User{}, &models.Contract{})

	Db = db
}
