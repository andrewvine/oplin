// Package wiring provides the wiring for the lineage service
package wiring

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"oplin/internal/lineage/api"
	"oplin/internal/lineage/htmx/datasets"
	"oplin/internal/lineage/htmx/jobs"
	"oplin/internal/lineage/htmx/requests"
	"oplin/internal/lineage/htmx/runs"
	"oplin/internal/lineage/ops"
	"oplin/resources"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var dbHost string
var dbName string
var dbUser string
var dbPassword string
var dbSslmode string
var dbPort int

// init parses the command line flags
func init() {
	flag.StringVar(&dbHost, "db_host", "", "the name of the host")
	flag.StringVar(&dbName, "db_name", "", "the name of the database")
	flag.StringVar(&dbUser, "db_user", "", "the name of the user")
	flag.StringVar(&dbPassword, "db_password", "", "the database users password")
	flag.StringVar(&dbSslmode, "db_sslmode", "", "the sslmode (disable)")
	flag.IntVar(&dbPort, "db_port", 0, "the database port")
}

// firstSet returns the first non-empty string in the slice of strings
func firstSet(xs ...string) string {
	for _, x := range xs {
		if x != "" {
			return x
		}
	}
	return xs[len(xs)-1]
}

// buildDSN builds the data source name for the database
func buildDSN() string {
	portString := ""
	if dbPort > 0 {
		portString = strconv.Itoa(dbPort)
	}

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		firstSet(dbHost, os.Getenv("OPLIN_DB_HOST"), "localhost"),
		firstSet(dbUser, os.Getenv("OPLIN_DB_USER"), "oplin"),
		firstSet(dbPassword, os.Getenv("OPLIN_DB_PASSWORD"), "topsecret"),
		firstSet(dbName, os.Getenv("OPLIN_DB_NAME"), "oplin"),
		firstSet(portString, os.Getenv("OPLIN_DB_PORT"), "5432"),
		firstSet(dbSslmode, os.Getenv("OPLIN_DB_SSLMODE"), "disable"))
}

// NewGinEngine creates a new gin.Engine
func NewGinEngine() *gin.Engine {
	r := gin.Default()
	r.Use(gin.CustomRecovery(ErrorHandler))
	return r
}

// ErrorHandler handles errors
func ErrorHandler(c *gin.Context, err any) {
	e, ok := err.(error)
	if !ok {
		log.Printf("In recovery could not convert any[%v] to error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("%v", err),
		})
	} else {
		c.Error(e)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": e.Error(),
		})
	}
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func bytesToString(b []byte) string {
	return fmt.Sprintf("%s", b)
}

// SetupLineage sets up the lineage service
func SetupLineage(r *gin.Engine) error {
	dsn := buildDSN()
	return setupLineage(r, dsn)
}

func setupLineage(r *gin.Engine, dsn string) error {
	DB, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	deps := &WiringDeps{DB: DB}

	// Initialize the schema if it does not exist
	ctx := context.Background()
	err = ops.InitializeDB(ctx, deps)
	if err != nil {
		return err
	}

	SetupRouter(r, deps)
	return nil
}

// SetupRouter sets up the router for the lineage service
func SetupRouter(r *gin.Engine,
	deps Deps,
) {
	// API
	r.POST("/api/v1/lineage", api.MakeCreateWithOpenLineageRunEvent(deps))

	// Static
	static, err := fs.Sub(resources.Static, "static")
	if err != nil {
		log.Fatalf("could not load static")
		return
	}
	r.StaticFS("/static", http.FS(static))

	// Templates
	r.SetFuncMap(template.FuncMap{
		"formatTime":    formatTime,
		"bytesToString": bytesToString,
	})
	templ := template.Must(template.New("").Funcs(r.FuncMap).ParseFS(resources.Templates, "templates/**/*.html"))
	r.SetHTMLTemplate(templ)

	listDatasets := datasets.MakeListDatasets(deps)

	// Datasets
	r.GET("/lineage/datasets/versions/fields", datasets.MakeGetDatasetVersionFields(deps))
	r.GET("/lineage/datasets/versions/lineage", datasets.MakeGetDatasetVersionLineage(deps))
	r.GET("/lineage/datasets/:id", datasets.MakeGetDataset(deps))
	r.GET("/lineage/datasets/:id/fields", datasets.MakeGetDatasetFields(deps))
	r.GET("/lineage/datasets/:id/lineage", datasets.MakeGetDatasetLineage(deps))
	r.GET("/lineage/datasets/:id/ownership", datasets.MakeGetDatasetOwnership(deps))
	r.GET("/lineage/datasets/:id/quality", datasets.MakeGetDatasetQuality(deps))
	r.GET("/lineage/datasets/:id/more", datasets.MakeGetDatasetMore(deps))
	r.GET("/lineage/datasets", listDatasets)

	// Jobs
	r.GET("/lineage/jobs/:id/runs", jobs.MakeGetJobRuns(deps))
	r.GET("/lineage/jobs/:id/ownership", jobs.MakeGetJobOwnership(deps))
	r.GET("/lineage/jobs/:id/sourcecode", jobs.MakeGetJobSourceCode(deps))
	r.GET("/lineage/jobs/:id", jobs.MakeGetJob(deps))
	r.GET("/lineage/jobs", jobs.MakeListJobs(deps))

	// Requests
	r.GET("/lineage/requests", requests.MakeGetRequests(deps))

	// Runs
	r.GET("/lineage/runs/:id", runs.MakeGetRun(deps))

	// Home
	r.GET("/index", listDatasets)
	r.GET("/", listDatasets)
}
