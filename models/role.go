package models

import (
	"gogpt/config"
)

type Role struct {
	ID       int    `json:"id"`
	RoleName string `json:"role_name"`
}

type Permission struct {
	ID             int    `json:"id"`
	RoleID         int    `json:"role_id"`
	PermissionName string `json:"permission_name"`
}

// Kullanıcının rolüne göre izin kontrolü
func CheckPermission(roleName, permissionName string) (bool, error) {
	query := `
        SELECT p.id
        FROM permissions p
        JOIN roles r ON p.role_id = r.id
        WHERE r.role_name = ? AND p.permission_name = ?
    `
	var permissionID int
	err := config.DB.QueryRow(query, roleName, permissionName).Scan(&permissionID)
	if err != nil {
		return false, err
	}
	return true, nil
}
