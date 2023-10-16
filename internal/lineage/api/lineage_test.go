package api_test

import (
	"oplin/internal/env"
	"oplin/internal/lineage/wiring"
	"oplin/internal/utils"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupSuite(tb testing.TB) (*gin.Engine, func(tb testing.TB)) {
	db := utils.GetTestDB()
	filename := env.PrependProjectPath("internal/lineage/db/schema.sql")
	err := utils.RunSQLFile(db, filename)
	if err != nil {
		log.Fatalf("Cannot load schema[%v]", err)
	}

	r := wiring.NewGinEngine()
	err = wiring.SetupTestLineage(r)
	if err != nil {
		log.Fatalf("Cannot setup controller[%v]", err)
	}
	return r, func(tb testing.TB) {
	}
}

func TestCreateWithOpenLineageRunEvent(t *testing.T) {
	r, teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	w := httptest.NewRecorder()
	payload := `
	{"run": 
		{"runId": "6871ff97-c518-4081-97aa-8520f7e634b7", 
		"facets": {"a": {"b": "1"}}}, 
	"job": {"namespace": "abs", "name": "xyz"}, 
	"inputs": [{"namespace": "dns1", "name": "table1"}], 
	"outputs": [{"namespace": "dns2", "name": "table2"}], 
	"eventType": "start", 
	"eventTime": "2023-02-05T15:48:28.660754+02:00", 
	"Producer": "R-Kelly"}
	`

	req, _ := http.NewRequest(
		"POST", "/api/v1/lineage", strings.NewReader(payload),
	)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
