package jobs

import (
	"oplin/internal/lineage/htmx"
	"oplin/internal/lineage/ops"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TabItem struct {
	Key   string
	Text  string
	Href  string
	Role  string
	Icon  string
}

type Breadcrumb struct {
	Text  string
	Href  string
}

var TabItems = []TabItem{
	{Key: "runs", Text: "Runs", Href: "/lineage/jobs/%d/runs"},
	{Key: "ownership", Text: "Ownership", Href: "/lineage/jobs/%d/ownership"},
	{Key: "sourcecode", Text: "Source Code", Href: "/lineage/jobs/%d/sourcecode"},
}

func buildBreadcrumbs(jobID int64, text string) []Breadcrumb {
	return []Breadcrumb{
		{Href: "/lineage/jobs", Text: "Jobs"},
		{Href: fmt.Sprintf("/lineage/jobs/%d", jobID), Text: text},
	}
}

func buildTabItems(chosenKey string, jobID int64) []TabItem {
	var res []TabItem
	for _, it := range TabItems {
		it.Href = fmt.Sprintf(it.Href, jobID)
		if it.Key == chosenKey {
			it.Role = "button"
		}
		res = append(res, it)
	}
	return res
}

func MakeListJobs(deps htmx.Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		jss, err := ops.ListJobsWithNamespaces(ctx, deps)
		if err != nil {
			c.HTML(http.StatusOK, "lineage/error.html", gin.H{})
			return
		}
		c.HTML(http.StatusOK, "lineage/jobs-list.html", gin.H{
			"Title":     "Jobs",
			"Jobs":      jss,
			"MenuItems": htmx.BuildMenuItems("jobs"),
		})
	}
}

func MakeGetJob(deps htmx.Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		s := c.Param("id")
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			c.HTML(http.StatusOK, "lineage/error.html", gin.H{})
			return
		}

		jns, err := ops.GetJobWithNamespace(ctx, deps, id)
		if err != nil {
			c.HTML(http.StatusOK, "lineage/error.html", gin.H{})
			return
		}
		runs, err := ops.ListRunsByJobVersionID(ctx, deps, jns.Job.CurrentVersionID)
		if err != nil {
			c.HTML(http.StatusOK, "lineage/error.html", gin.H{})
			return
		}

		title := fmt.Sprintf("%s %s", jns.JobNamespace.Name, jns.Job.Name)

		c.HTML(http.StatusOK, "lineage/jobs-detail.html", gin.H{
			"Breadcrumbs":      buildBreadcrumbs(jns.Job.ID, title),
			"Title":            title,
			"JobWithNamespace": jns,
			"Runs":             runs,
			"TabItems":         buildTabItems("runs", jns.Job.ID),
			"MenuItems":        htmx.BuildMenuItems("jobs"),
		})
	}
}

func MakeGetJobRuns(deps htmx.Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		s := c.Param("id")
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			c.HTML(http.StatusOK, "lineage/error.html", gin.H{})
			return
		}

		jns, err := ops.GetJobWithNamespace(ctx, deps, id)
		if err != nil {
			c.HTML(http.StatusOK, "lineage/error.html", gin.H{})
			return
		}
		runs, err := ops.ListRunsByJobVersionID(ctx, deps, jns.Job.CurrentVersionID)
		if err != nil {
			c.HTML(http.StatusOK, "lineage/error.html", gin.H{})
			return
		}

		title := fmt.Sprintf("%s %s", jns.JobNamespace.Name, jns.Job.Name)

		c.HTML(http.StatusOK, "lineage/jobs-runs.html", gin.H{
			"Breadcrumbs":      buildBreadcrumbs(jns.Job.ID, title),
			"Title":            title,
			"JobWithNamespace": jns,
			"Runs":             runs,
			"TabItems":         buildTabItems("runs", jns.Job.ID),
		})
	}
}

func MakeGetJobOwnership(deps htmx.Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		s := c.Param("id")
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			c.HTML(http.StatusOK, "lineage/error.html", gin.H{})
			return
		}

		jns, err := ops.GetJobWithNamespace(ctx, deps, id)
		if err != nil {
			c.HTML(http.StatusOK, "lineage/error.html", gin.H{})
			return
		}

		title := fmt.Sprintf("%s %s", jns.JobNamespace.Name, jns.Job.Name)

		c.HTML(http.StatusOK, "lineage/jobs-ownership.html", gin.H{
			"Breadcrumbs":      buildBreadcrumbs(jns.Job.ID, title),
			"Title":            title,
			"JobWithNamespace": jns,
			"TabItems":         buildTabItems("ownership", jns.Job.ID),
		})
	}
}

func MakeGetJobSourceCode(deps htmx.Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		s := c.Param("id")
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			c.HTML(http.StatusOK, "lineage/error.html", gin.H{})
			return
		}

		jns, err := ops.GetJobWithNamespace(ctx, deps, id)
		if err != nil {
			c.HTML(http.StatusOK, "lineage/error.html", gin.H{})
			return
		}

		title := fmt.Sprintf("%s %s", jns.JobNamespace.Name, jns.Job.Name)

		c.HTML(http.StatusOK, "lineage/jobs-sourcecode.html", gin.H{
			"Breadcrumbs":      buildBreadcrumbs(jns.Job.ID, title),
			"Title":            title,
			"JobWithNamespace": jns,
			"TabItems":         buildTabItems("sourcecode", jns.Job.ID),
		})
	}
}
