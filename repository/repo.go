package repository

import (
	"github.com/staffea-test/entity"
	"gorm.io/gorm"
)

type Repo interface {
	CreateUser(u entity.User) error
	GetUserByID(id string) (entity.User, error)
	GetInviteByID(id string) (entity.Invite, error)
	RemoveInviteByID(id string) error
	WithTrx(trxHandle *gorm.DB) Repo
	GetDB() *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &repo{db: db}
}

type repo struct {
	db *gorm.DB
}

func (r *repo) CreateUser(u entity.User) error {
	err := r.db.Create(&u).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) GetUserByID(id string) (entity.User, error) {
	var user entity.User
	err := r.db.Table("users").First(&user, "id = ?", id).Error
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *repo) GetInviteByID(id string) (entity.Invite, error) {
	var i entity.Invite
	err := r.db.Table("invites").First(&i, "id = ?", id).Error
	if err != nil {
		return entity.Invite{}, err
	}

	return i, nil
}

func (r *repo) RemoveInviteByID(id string) error {
	err := r.db.Delete(entity.Invite{}, "id = ?", id).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) WithTrx(trxHandle *gorm.DB) Repo {
	if trxHandle == nil {
		return r
	}
	return &repo{db: trxHandle}
}

func (r *repo) GetDB() *gorm.DB {
	return r.db
}
