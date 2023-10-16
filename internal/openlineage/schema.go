package openlineage

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Dataset struct {
	Namespace string          `json:"namespace"`
	Name      string          `json:"name"`
	Facets    json.RawMessage `json:"facets"`
}

func NewDataset(ns string, name string, msg json.RawMessage) *Dataset {
	return &Dataset{Namespace: ns, Name: name, Facets: msg}
}

type InputDataset struct {
	Dataset
	InputFacets json.RawMessage `json:"inputFacets"`
}

type OutputDataset struct {
	Dataset
	OutputFacets json.RawMessage `json:"outputFacets"`
}

type Job struct {
	Namespace string          `json:"namespace"`
	Name      string          `json:"name"`
	Facets    json.RawMessage `json:"facets"`
}

func NewJob(ns string, name string, msg json.RawMessage) *Job {
	return &Job{Namespace: ns, Name: name, Facets: msg}
}

type Run struct {
	ID     uuid.UUID       `json:"runId"`
	Facets json.RawMessage `json:"facets"`
}

func NewRun(id uuid.UUID) *Run {
	return &Run{ID: id}
}

func NewRunEvent() *RunEvent {
	ev := RunEvent{}
	return &ev
}

type RunEvent struct {
	EventType string          `json:"eventType"`
	EventTime time.Time       `json:"eventTime"`
	Run       *Run            `json:"run"`
	Job       *Job            `json:"job"`
	Inputs    []InputDataset  `json:"inputs"`
	Outputs   []OutputDataset `json:"outputs"`
	Producer  string          `json:"producer"`
	SchemaURL string          `json:"schemaURL"`
}
