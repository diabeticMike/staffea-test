package entity

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	Login    string    `gorm:"type:string;column:login" json:"login"`
	Password string    `gorm:"type:string;column:password" json:"password"`
	Secret   [16]byte  `gorm:"type:bytea;column:secret" json:"-"`
}
