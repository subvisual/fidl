package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/subvisual/fidl/types"
)

type AuthorizationStatus int8

const (
	AuthorizationOpen AuthorizationStatus = iota + 1
	AuthorizationLocked
)

func (a AuthorizationStatus) String() string {
	switch a {
	case AuthorizationOpen:
		return "Open"
	case AuthorizationLocked:
		return "Locked"
	default:
		return "Unknown" // nolint:goconst
	}
}

type Authorization struct {
	ID        int64               `db:"id"`
	UUID      uuid.UUID           `db:"uuid"`
	Balance   types.FIL           `db:"balance"`
	Proxy     string              `db:"proxy"`
	Status    AuthorizationStatus `db:"status_id"`
	CreatedAt time.Time           `db:"created_at"`
	UpdatedAt time.Time           `db:"updated_at"`
}
