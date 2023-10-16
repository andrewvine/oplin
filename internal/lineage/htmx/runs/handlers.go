package runs

import (
	"oplin/internal/lineage/htmx"
	"oplin/internal/lineage/ops"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Breadcrumb struct {
	Text  string
	Href  string
	Class string
}

func buildBreadcrumbs(jobID int64, jobText string, runID int64, runText string) []Breadcrumb {
	return []Breadcrumb{
		{Href: "/lineage/jobs", Text: "Jobs"},
		{Href: fmt.Sprintf("/lineage/jobs/%d", jobID), Text: jobText},
		{Href: fmt.Sprintf("/lineage/runs/%d", runID), Text: runText, Class: "is-active"},
	}
}

func MakeGetRun(deps htmx.Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		s := c.Param("id")
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}

		run, err := ops.GetRunWithID(ctx, deps, id)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}

		jv, err := ops.GetJobVersionByID(ctx, deps, run.JobVersionID)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}

		jns, err := ops.GetJobWithNamespace(ctx, deps, jv.JobID)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}

		events, err := ops.ListRunEventsByRunID(ctx, deps, id)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}

		ioDatasets, err := ops.ListRunDatasetVersionsWithRelationshipsByRunID(ctx, deps, id)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}

		jobTitle := fmt.Sprintf("%s %s", jns.JobNamespace.Name, jns.Job.Name)

		c.HTML(http.StatusOK, "lineage/jobs-runs-events.html", gin.H{
			"Breadcrumbs": buildBreadcrumbs(jns.Job.ID, jobTitle, run.ID, run.StartedAt.Format("2006-01-02 15:04:05")),
			"Run":         run,
			"Events":      events,
			"IODatasets":  ioDatasets,
			"MenuItems":   htmx.BuildMenuItems("jobs"),
		})
	}
}
