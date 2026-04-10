package data

import (
	"context"
	"database/sql"
	"slices"
	"time"

	"github.com/lib/pq"
)

type Permissions []string

func (p Permissions) Include(code string) bool {
	// Permissions are stored as plain codes, so membership checks stay simple.
	return slices.Contains(p, code)
}

type PermissionModel struct {
	DB *sql.DB
}

func (m PermissionModel) GetAllForUSer(userID int64) (Permissions, error) {
	// Join through the mapping table so we return the permission codes assigned to one user.
	query := `
		SELECT permissions.code
		FROM permissions
		INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
		INNER JOIN users ON users_permissions.user_id = users.id
		WHERE users.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var permissions Permissions

	for rows.Next() {
		var permission string

		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Return the collected codes in scan order; callers only care about the set of grants.
	return permissions, nil
}

func (m PermissionModel) AddForUSer(userID int64, codes ...string) error {
	// Insert only permissions that already exist; the database enforces the valid code list.
	query := `
		INSERT INTO users_permissions
		SELECT $1, permissions.id FROM permissions WHERE permissions.code = ANY($2)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.QueryContext(ctx, query, userID, pq.Array(codes))
	return err
}
