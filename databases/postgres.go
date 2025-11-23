package databases

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var (
	NAME     string
	PORT     string
	HOST     string
	USER     string
	PASSWORD string
)

var DatabaseQuery *sql.DB

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env")
	}

	NAME = os.Getenv("DB_NAME")
	PORT = os.Getenv("DB_PORT")
	HOST = os.Getenv("DB_HOST")
	USER = os.Getenv("DB_USER")
	PASSWORD = os.Getenv("DB_PASSWORD")
}


func ConnectToPostgres(app * fiber.App) error {

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, USER, PASSWORD, NAME,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	// cek koneksi
	if err := db.Ping(); err != nil {
		return err
	}

	DatabaseQuery = db
	fmt.Println("✅ Connected to Postgres!")
	return nil

}
