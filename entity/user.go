package entity

import (
	"errors"
	"github.com/google/uuid"
	"regexp"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;unique"`
	Login    string    `gorm:"type:string;column:login;unique" json:"login"`
	Password string    `gorm:"type:string;column:password" json:"password"`
}

var (
	Secret = uuid.New()
)

type UserAuthRequest struct {
	ID       uuid.UUID `json:"uuid"`
	Login    string    `json:"login"`
	Password string    `json:"password"`
}

type UserAuthResponse struct {
	Login        string `json:"login"`
	RefreshToken string `json:"refreshToken"`
	AccessToken  string `json:"accessToken"`
	Provider     int64  `json:"provider"`
	Locale       string `json:"locale"`
}

func (u UserAuthRequest) Validate() error {
	_, err := regexp.Match(
		`^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$`,
		[]byte(u.ID.String()))
	if err != nil {
		return err
	}

	_, err = regexp.Match(
		`^[A-Za-z0-9]+$`,
		[]byte(u.Login))
	if err != nil {
		return err
	}

	if len(u.Login) < 6 || len(u.Login) > 255 {
		return errors.New("login length is incorrect")
	}

	if len(u.Password) < 8 || len(u.Password) > 255 {
		return errors.New("password length is incorrect")
	}

	return nil
}
