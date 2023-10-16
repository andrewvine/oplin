package requests

import (
	"oplin/internal/lineage/htmx"
	"oplin/internal/lineage/ops"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MakeGetRequests(deps htmx.Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		reqs, err := ops.ListRequests(ctx, deps)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}

		c.HTML(http.StatusOK, "lineage/requests-list.html", gin.H{
			"Requests":  reqs,
			"MenuItems": htmx.BuildMenuItems("events"),
		})
	}
}
