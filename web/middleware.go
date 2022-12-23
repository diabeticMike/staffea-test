package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/staffea-test/entity"
	"net/http"
	"strings"
)

type Middleware struct {
}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (mw *Middleware) Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := strings.TrimPrefix(c.GetHeader("token"), "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(entity.Secret.String()), nil
		})
		if err != nil {
			if err.Error() != "Token is expired" {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			refreshString := strings.TrimPrefix(c.GetHeader("refresh"), "Bearer ")
			refresh, err := jwt.Parse(refreshString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				return []byte(entity.Secret.String()), nil
			})
			if err != nil {
				if err.Error() == "Token is expired" {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "resign in again please"})
					return
				}
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			claims, ok := refresh.Claims.(jwt.MapClaims)
			if !ok || !refresh.Valid {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid claim map"})
				return
			}
			userID, ok := claims["id"].(string)
			if !ok {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "no user id provided"})
				return
			}
			newToken, newRefreshToken, err := GenerateNewTokens(userID)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, &entity.UserAuthResponse{
				Login:        "",
				RefreshToken: newToken,
				AccessToken:  newRefreshToken,
				Provider:     1,
				Locale:       "fr_FR",
			})
			return
		}

		_, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid claim map"})
			return
		}

		c.Next()
	}
}
