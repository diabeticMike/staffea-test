package repository

import (
	"crypto/md5"
	"github.com/google/uuid"
	"github.com/gorm"
	"github.com/staffea-test/entity"
)

type UserDB interface {
	Create(u entity.User) error
}

type userDB struct {
	db *gorm.DB
}

func (uDB *userDB) Create(u entity.User) error {
	random, err := uuid.New().MarshalBinary()
	if err != nil {
		return err
	}

	u.Secret = md5.Sum(random)
	u.ID = uuid.New()
	err = uDB.db.Create(&u).Error
	if err != nil {
		return err
	}

	return nil
}
