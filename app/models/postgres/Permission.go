package postgres

type Permission struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Description string `json:"description"`
}

type RolePermission struct {
	RoleId string `json:"role_id"`
	PermissionId string `json:"permission_id"`
}
