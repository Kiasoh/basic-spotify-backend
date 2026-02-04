package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// FloatVector is a custom type for handling []float64 as JSONB in PostgreSQL
type FloatVector []float64

// Scan implements the sql.Scanner interface.
func (fv *FloatVector) Scan(src any) error {
	if src == nil {
		*fv = nil
		return nil
	}
	switch s := src.(type) {
	case []byte:
		return json.Unmarshal(s, fv)
	case string:
		return json.Unmarshal([]byte(s), fv)
	default:
		return errors.New("unsupported type for FloatVector scanning")
	}
}

// Value implements the driver.Valuer interface.
func (fv FloatVector) Value() (driver.Value, error) {
	if fv == nil {
		return nil, nil
	}
	return json.Marshal(fv)
}

type User struct {
	ID               int         `json:"id"`
	Username         string      `json:"username"`
	Password         Password    `json:"-"`
	AvgInterest      FloatVector `json:"avg_interest"` // Now []float64
	RecommPlaylistID int         `json:"recomm_playlist_id"`
	CreatedAt        time.Time   `json:"created_at"`
}
