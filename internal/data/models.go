package data

import (
	"database/sql"
	"errors"
)

var (
	// These package-level errors let handlers distinguish missing rows from stale updates.
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// Models groups the database access layers behind one shared entry point.
type Models struct {
	Movies      MovieModel
	Permissions PermissionModel
	Token       TokenModel
	Users       UserModel
}

// NewModels wires each model to the same database handle.
func NewModels(db *sql.DB) Models {
	return Models{
		Movies:      MovieModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Token:       TokenModel{DB: db},
		Users:       UserModel{DB: db},
	}
}
