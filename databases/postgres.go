package databases

import (
	"database/sql"
	"fmt"
	"os"

	
	"github.com/joho/godotenv"
	 _ "github.com/lib/pq"
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

func ConnectToPostgres() (*sql.DB, error) {

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, USER, PASSWORD, NAME,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	// cek koneksi
	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	fmt.Println("Postgres Connected ....")

	DatabaseQuery = db
	return db, err
}
