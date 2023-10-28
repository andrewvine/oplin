package openlineage_test

import (
	"context"
	"fmt"
	"log"
	"oplin/internal/lineage"
	ops "oplin/internal/lineage/ops"
	ol_ops "oplin/internal/lineage/ops/openlineage"
	"oplin/internal/openlineage"
	"oplin/internal/utils"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupSuite(tb testing.TB) (*ops.TestDeps, func(tb testing.TB)) {
	ctx := context.Background()
	db := utils.GetTestDB()
	deps := &ops.TestDeps{DB: db}

	err := ops.InitializeTestDB(ctx, deps)
	if err != nil {
		log.Fatalf("Cannot setup controller[%v]", err)
	}
	return deps, func(tb testing.TB) {
	}
}

func assertTimeEqual(t *testing.T, a, b time.Time) {
	f := "2006/01/02 10:08:01.123456"
	assert.Equal(t, a.Format(f), b.Format(f))
}

func getRunEvent(runUUID uuid.UUID, eventTime time.Time) openlineage.RunEvent {
	producer := "https://github.com/OpenLineage/OpenLineage/blob/v1-0-0/client"
	schemaURL := "https://openlineage.io/spec/0-0-1/OpenLineage.json"

	ev := openlineage.RunEvent{}
	ev.EventTime = eventTime
	ev.EventType = "start"
	ev.Job = openlineage.NewJob("airflow", "orders.monthly_summary", nil)
	ev.Run = openlineage.NewRun(runUUID)
	ev.Producer = producer
	ev.SchemaURL = schemaURL

	return ev
}

func TestNominalTimeAndParentFacets(t *testing.T) {
	deps, teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	ctx := context.Background()

	start := time.Now().UTC()
	runUUID := uuid.New()

	ev := getRunEvent(runUUID, start)
	ev.Run.Facets = []byte(
		`{"nominalTime": {
			"_producer": "https://some.producer.com/version/1.0",
			"_schemaURL": "https://github.com/OpenLineage/OpenLineage/blob/main/spec/facets/SQLJobFacet.json",
			"nominalStartTime": "2020-12-17T03:00:00.000Z",
			"nominalEndTime": "2020-12-17T03:05:00.000Z"
		}}`,
	)
	runEvent, err := ol_ops.CreateWithOpenLineageRunEvent(ctx, deps, &ev)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), runEvent.ID)
	assert.Equal(t, int64(1), runEvent.RunID)
	assert.Equal(t, lineage.RunEventTypeStart, runEvent.EventType)
	assertTimeEqual(t, start, runEvent.EventTime)

	run, err := ops.GetRunWithID(ctx, deps, runEvent.RunID)
	assert.Nil(t, err)

	// Check LastEventType
	assert.Equal(t, runEvent.EventType, run.LastEventType)

	// Check StartedAt
	assertTimeEqual(t, runEvent.EventTime, run.StartedAt)

	// Check nominalTime facet
	assert.Equal(t, "2020-12-17T03:00:00Z", run.NominalStartedAt.Format(time.RFC3339))
	assert.Equal(t, "2020-12-17T03:05:00Z", run.NominalEndedAt.Format(time.RFC3339))

	end := start.Add(time.Second * 5)
	ev = getRunEvent(runUUID, end)
	ev.EventType = "complete"

	runEvent, err = ol_ops.CreateWithOpenLineageRunEvent(ctx, deps, &ev)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), runEvent.ID)
	assert.Equal(t, int64(1), runEvent.RunID)
	assert.Equal(t, lineage.RunEventTypeComplete, runEvent.EventType)
	assertTimeEqual(t, end, runEvent.EventTime)

	run, err = ops.GetRunWithID(ctx, deps, runEvent.RunID)
	assert.Nil(t, err)
	assert.Equal(t, runEvent.EventType, run.LastEventType)
	assertTimeEqual(t, start, run.StartedAt)
	assertTimeEqual(t, end, run.EndedAt)

	start = time.Now().UTC()
	parentUUID := runUUID
	runUUID = uuid.New()

	// Check parent facet
	ev = getRunEvent(runUUID, start)
	ev.Run.Facets = []byte(
		fmt.Sprintf(`{"parent": {
			"job": {
				"name": "the-execution-parent-job", 
				"namespace": "the-namespace"
			  },
			  "run": {
				"runId": "%s"
			  }
		}}`, parentUUID),
	)

	runEvent, err = ol_ops.CreateWithOpenLineageRunEvent(ctx, deps, &ev)
	assert.Nil(t, err)
	assert.Equal(t, int64(3), runEvent.ID)
	assert.Equal(t, int64(2), runEvent.RunID)

	run, err = ops.GetRunWithID(ctx, deps, runEvent.RunID)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), run.ParentRunID)
}

func TestErrorMessageFacets(t *testing.T) {
	deps, teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	ctx := context.Background()

	start := time.Now().UTC()
	runUUID := uuid.New()

	ev := getRunEvent(runUUID, start)
	ev.EventType = "fail"
	ev.Run.Facets = []byte(
		`{"errorMessage": {
			"message": "Bang!",
			"programmingLanguage": "Go",
			"stackTrace": "some stacktrace"
		}}`,
	)

	runEvent, err := ol_ops.CreateWithOpenLineageRunEvent(ctx, deps, &ev)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), runEvent.ID)
	assert.Equal(t, int64(1), runEvent.RunID)
	assert.Equal(t, lineage.RunEventTypeFail, runEvent.EventType)

	run, err := ops.GetRunWithID(ctx, deps, runEvent.RunID)
	assert.Nil(t, err)
	assert.Equal(t, runEvent.EventType, run.LastEventType)
	assert.Equal(t, "Bang!", run.ErrorMessage)
	assert.Equal(t, "Go", run.ProgrammingLanguage)
	assert.Equal(t, "some stacktrace", run.Stacktrace)
}

func TestColumnLineageFacets(t *testing.T) {
	deps, teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	ctx := context.Background()

	start := time.Now().UTC()
	runUUID := uuid.New()
	ev := getRunEvent(runUUID, start)
	facets := []byte(
		`{"columnLineage": {
		"_producer": "https://github.com/MarquezProject/marquez/blob/main/docker/metadata.json",
		"_schemaURL": "https://openlineage.io/spec/facets/1-0-1/ColumnLineageDatasetFacet.json",
		"fields": {
			"order_id": {
				"inputFields": [
					{
						"namespace": "food_delivery",
						"name": "public.delivery_7_days",
						"field": "order_id"
					}
				]
			},
			"order_placed_on": {
				"inputFields": [
					{
						"namespace": "food_delivery",
						"name": "public.delivery_7_days",
						"field": "order_placed_on"
					}
				]
			},
			"order_delivered_on": {
				"inputFields": [
					{
						"namespace": "food_delivery",
						"name": "public.delivery_7_days",
						"field": "order_delivered_on"
					}
				]
			},
			"order_delivery_time": {
				"inputFields": [
					{
						"namespace": "food_delivery",
						"name": "public.delivery_7_days",
						"field": "order_placed_on"
					},
					{
						"namespace": "food_delivery",
						"name": "public.delivery_7_days",
						"field": "order_delivered_on"
					}
				]
			}
		}
	}}`,
	)
	ev.Inputs = []openlineage.InputDataset{
		{
			Dataset: openlineage.Dataset{
				Namespace: "food_delivery",
				Name:      "public.delivery_7_days",
				Facets:    facets,
			},
		},
	}
	ev.Outputs = []openlineage.OutputDataset{
		{
			Dataset: openlineage.Dataset{
				Namespace: "food_delivery",
				Name:      "public.top_delivery_times",
				Facets:    facets,
			},
		},
	}
	ev.EventType = "complete"
	_, err := ol_ops.CreateWithOpenLineageRunEvent(ctx, deps, &ev)
	assert.Nil(t, err)
}
