package api_test

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"oplin/internal/lineage/api"
	"oplin/internal/lineage/ops"
	"oplin/internal/lineage/wiring"
	"oplin/internal/utils"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupSuite(tb testing.TB) (*gin.Engine, func(tb testing.TB)) {
	ctx := context.Background()
	db := utils.GetTestDB()

	deps := api.TestDeps{DB: db}
	r := wiring.NewGinEngine()
	err := ops.InitializeTestDB(ctx, &deps)
	wiring.SetupRouter(r, &deps)
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
