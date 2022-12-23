package web

import (
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/staffea-test/entity"
	"github.com/staffea-test/repository"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type Controller struct {
	repo repository.Repo
}

func NewController(repo repository.Repo) *Controller {
	return &Controller{repo: repo}
}

func (ctl *Controller) HandleAuthentication(c *gin.Context) {
	var userInfo entity.UserAuthRequest
	err := c.BindJSON(&userInfo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	err = userInfo.Validate()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	_, err = ctl.repo.GetInviteByID(userInfo.ID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "you have no invite with provided uuid"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := entity.User{
		ID:       uuid.New(),
		Login:    userInfo.Login,
		Password: userInfo.Password,
	}

	tx := ctl.repo.GetDB().Begin()
	txDB := ctl.repo.WithTrx(tx)

	err = txDB.RemoveInviteByID(userInfo.ID.String())
	if err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = txDB.CreateUser(user)
	if err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = tx.Commit().Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tokenString, refreshString, err := GenerateNewTokens(user.ID.String())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("token", "Bearer "+tokenString)
	c.Header("refresh", "Bearer "+refreshString)

	fmt.Println("OK")
}

func GenerateNewTokens(userID string) (string, string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().UTC().Add(time.Minute * 1).Unix(),
	})

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().UTC().Add(time.Minute * 6).Unix(),
	})

	tokenString, err := token.SignedString([]byte(entity.Secret.String()))
	if err != nil {
		return "", "", err
	}
	refreshString, err := refresh.SignedString([]byte(entity.Secret.String()))
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshString, nil
}
