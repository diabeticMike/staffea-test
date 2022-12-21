package entity

import "github.com/google/uuid"

type Invite struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
}
