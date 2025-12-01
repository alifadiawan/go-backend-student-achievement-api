package postgres

import (
	models "backendUAS/app/models/postgres"
	"backendUAS/databases"
)

func GetAllUserRepository() ([]models.User, error) {

	Query, err := databases.DatabaseQuery.Query(`
		SELECT 
			u.id, 
			u.username, 
			u.email, 
			u.full_name, 
			u.role_id, 
			r.name,
			u.is_active,
			u.created_at,
			u.updated_at 
		FROM users as u
		JOIN roles as r on r.id = u.role_id
	`)
	if err != nil {
		return nil, err
	}
	defer Query.Close()

	var Users []models.User

	for Query.Next(){
		var item models.User

		foreach := Query.Scan(
			&item.ID,
			&item.Username,
			&item.Email,
			&item.FullName,
			&item.RoleID,
			&item.RoleName,
			&item.IsActive,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if foreach != err {
			return nil, err
		}

		Users = append(Users, item)
	}

	return Users, nil

}
