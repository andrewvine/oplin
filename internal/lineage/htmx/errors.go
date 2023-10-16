package htmx

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rotisserie/eris"
)

type Err struct {
	E error
}

func (e *Err) Error() string {
	return eris.ToString(e.E, true)
}

func InternalServerError(c *gin.Context, err error) {
	c.Error(&Err{E: err})
	c.HTML(http.StatusInternalServerError, "lineage/error.html", gin.H{})
}
