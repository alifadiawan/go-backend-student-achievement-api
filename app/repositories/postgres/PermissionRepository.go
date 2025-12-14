package postgres

import "backendUAS/databases"

func LoadPermissions(userID string) (map[string]bool, error) {
	rows, err := databases.DatabaseQuery.Query(`
		SELECT p.resource, p.action
		FROM users u
		JOIN role_permissions rp ON rp.role_id = u.role_id
		JOIN permissions p ON p.id = rp.permission_id
		WHERE u.id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	perms := make(map[string]bool)

	for rows.Next() {
		var resource, action string
		rows.Scan(&resource, &action)

		perms[resource+":"+action] = true
	}

	return perms, nil
}
