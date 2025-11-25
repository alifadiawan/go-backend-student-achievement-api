package postgres

import (
	"backendUAS/app/models/postgres"
	"backendUAS/databases"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func CheckPassword(raw string, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw))
	return err == nil
}

func Authenticate(email string, password string) (*postgres.User, error) {

	var User postgres.User

	err := databases.DatabaseQuery.QueryRow(`
			SELECT id, username, email, password_hash
			FROM users
			WHERE email = $1
		`, email).Scan(
		&User.ID, &User.Username, &User.Email, &User.PasswordHash,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("email tidak ditemukan")
	}

	if !CheckPassword(password, User.PasswordHash) {
		return nil, errors.New("password salah")
	}

	if err != nil {
		return nil, err
	}

	return &User, err

}
