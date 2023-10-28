package api

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	ol_ops "oplin/internal/lineage/ops/openlineage"
	"oplin/internal/openlineage"

	"github.com/gin-gonic/gin"
	"github.com/rotisserie/eris"
)

func writeError(c *gin.Context, httpStatus int, err error) {
	c.Error(err)
	formattedStr := eris.ToString(err, true)
	fmt.Println(formattedStr)

	c.JSON(httpStatus, gin.H{
		"error": err.Error(),
	})
}

func writeData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func MakeCreateWithOpenLineageRunEvent(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		ev := openlineage.NewRunEvent()

		body, _ := ioutil.ReadAll(c.Request.Body)
		fmt.Print(string(body))
		c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))

		if err := c.BindJSON(&ev); err != nil {
			c.Error(err)
		} else {
			id, err := ol_ops.CreateWithOpenLineageRunEvent(ctx, deps, ev)
			if err != nil {
				writeError(c, http.StatusInternalServerError, err)
			} else {
				writeData(c, id)
			}
		}
	}
}
