package entity

import (
	"github.com/google/uuid"
)

type Label struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	Login    string    `gorm:"type:string;column:login"`
	Password string    `gorm:"type:string;column:password"`
}
