package postgres

import (
	models "backendUAS/app/models/postgres"
	"backendUAS/databases"
	"database/sql"
	"errors"
)

func getAllRole() (*models.Role, error) {

	var roles models.Role

	err := databases.DatabaseQuery.QueryRow(`
		SELECT id, name, description
		FROM roles
	`).Scan(
		&roles.ID, &roles.Name, &roles, roles.Description,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("email tidak ditemukan")
	}

	return &roles, err

}
