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
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(entity.Secret.String()), nil
		})
		if err != nil {
			if err.Error() != "Token is expired" {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			refreshString := strings.TrimPrefix(c.GetHeader("refresh"), "Bearer ")
			refresh, err := jwt.Parse(refreshString, func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
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
			c.Header("token", "Bearer "+newToken)
			c.Header("refresh", "Bearer "+newRefreshToken)
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
