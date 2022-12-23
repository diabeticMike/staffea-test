package web

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func ListenTo(ctl *Controller, mw *Middleware) {
	r := gin.Default()

	r.POST("/ping", ctl.HandleAuthentication)
	r.Use(mw.Authorization())
	r.GET("/test", func(c *gin.Context) {
		_, err := c.Writer.Write([]byte(`OK`))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	})
	err := r.Run(":8000")
	if err != nil {
		log.Println(err.Error())
	}
}
