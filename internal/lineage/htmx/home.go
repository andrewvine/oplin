package htmx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MakeGetIndex() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "main/index.html", gin.H{
			"Title": "Beluga",
		})
	}
}
