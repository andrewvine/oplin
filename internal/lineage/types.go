package lineage

import (
	"errors"
	"fmt"
	"oplin/internal/openlineage"
	"strings"
	"time"
)

type RunEventType int

const (
	RunEventTypeUnknown  RunEventType = 0
	RunEventTypeStart    RunEventType = 1
	RunEventTypeRunning  RunEventType = 2
	RunEventTypeComplete RunEventType = 3
	RunEventTypeAbort    RunEventType = 4
	RunEventTypeFail     RunEventType = 5
	RunEventTypeOther    RunEventType = 6
	runEventTypeSentinal RunEventType = 7
)

type IOType int

const (
	IOTypeUnknown  IOType = 0
	IOTypeInput    IOType = 1
	IOTypeOutput   IOType = 2
	ioTypeSentinal IOType = 3
)

var runEventTypeMap = map[string]RunEventType{
	"start":    RunEventTypeStart,
	"running":  RunEventTypeRunning,
	"complete": RunEventTypeComplete,
	"abort":    RunEventTypeAbort,
	"fail":     RunEventTypeFail,
	"other":    RunEventTypeOther,
}

var runEventTypeToStringMap = map[RunEventType]string{
	RunEventTypeStart:    "start",
	RunEventTypeRunning:  "running",
	RunEventTypeComplete: "complete",
	RunEventTypeAbort:    "abort",
	RunEventTypeFail:     "fail",
	RunEventTypeOther:    "other",
}

func (r RunEventType) String() string {
	return strings.ToUpper(runEventTypeToStringMap[r])
}

func RunEventTypeFromString(str string) (RunEventType, error) {
	val, ok := runEventTypeMap[strings.ToLower(str)]
	if !ok {
		return RunEventTypeUnknown, errors.New(fmt.Sprintf("No type matching [%s]", str))
	}
	return val, nil
}

type RunEventDatasetType int

const (
	RunEventDatasetTypeUnkonwn  RunEventDatasetType = 0
	RunEventDatasetTypeInput    RunEventDatasetType = 1
	RunEventDatasetTypeOutput   RunEventDatasetType = 2
	runEventDatasetTypeSentinal RunEventDatasetType = 3
)

type Dataset struct {
	ID                 int64
	CurrentVersionID   int64
	DatasetNamespaceID int64
	Name               string
	Facets             openlineage.DatasetFacets
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type DatasetVersion struct {
	ID                 int64
	DatasetID          int64
	DatasetNamespaceID int64
	Name               string
	Facets             openlineage.DatasetFacets
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type RunEventDataset struct {
	ID                  int64
	DatasetID           int64
	RunEventID          int64
	RunEventDatasetType RunEventDatasetType
	Facets              openlineage.DatasetFacets
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type Job struct {
	ID               int64
	CurrentVersionID int64
	JobNamespaceID   int64
	Name             string
	Facets           openlineage.JobFacets
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type JobVersion struct {
	ID             int64
	JobID          int64
	JobNamespaceID int64
	Name           string
	Facets         openlineage.JobFacets
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewJob() *Job {
	return &Job{}
}

type Run struct {
	ID                  int64
	JobVersionID        int64
	Facets              openlineage.RunFacets
	ParentRunID         int64
	LastEventType       RunEventType
	NominalStartedAt    time.Time
	NominalEndedAt      time.Time
	StartedAt           time.Time
	EndedAt             time.Time
	ErrorMessage        string
	ProgrammingLanguage string
	Stacktrace          string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type JobNamespace struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DatasetNamespace struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RunEvent struct {
	ID        int64
	RunID     int64
	EventType RunEventType
	EventTime time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DatasetWithNamespace struct {
	Dataset          Dataset
	DatasetNamespace DatasetNamespace
}

type JobWithNamespace struct {
	Job          Job
	JobNamespace JobNamespace
}

type RunIODataset struct {
	DatasetVersionID int64
	RunID            int64
	IOType           IOType
	InputFacets      openlineage.InputDatasetFacets
	OutputFacets     openlineage.OutputDatasetFacets
	CreatedAt        time.Time
}

type RunIODatasetWithRelationships struct {
	RunIODataset     RunIODataset
	DatasetVersion   DatasetVersion
	DatasetNamespace DatasetNamespace
}

type Field struct {
	ID          int64
	Name        string
	DataType    string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Request struct {
	ID        int64
	Payload   []byte
	CreatedAt time.Time
}
