package openlineage

import "time"

type ErrorMessageFacet struct {
	Message             string `json:"message"`
	ProgrammingLanguage string `json:"programmingLanguage"`
	Stacktrace          string `json:"stackTrace"`
	Producer            string `json:"_producer"`
	SchemaURL           string `json:"_schemaURL"`
}

type ExternalQueryFacet struct {
	ExternalQueryID string `json:"externalQueryId"`
	Source          string `json:"source"`
	Producer        string `json:"_producer"`
	SchemaURL       string `json:"_schemaURL"`
}

type NominalTimeFacet struct {
	StartTime time.Time `json:"nominalStartTime"`
	EndTime   time.Time `json:"nominalEndTime"`
	Producer  string    `json:"_producer"`
	SchemaURL string    `json:"_schemaURL"`
}

type ParentJob struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type ParentRun struct {
	RunID string `json:"runId"`
}

type ParentRunFacet struct {
	ParentJob ParentJob `json:"job"`
	ParentRun ParentRun `json:"run"`
	Producer  string    `json:"_producer"`
	SchemaURL string    `json:"_schemaURL"`
}

type RunFacets struct {
	ErrorMessage  ErrorMessageFacet  `json:"errorMessage"`
	ExternalQuery ExternalQueryFacet `json:"externalQuery"`
	NominalTime   NominalTimeFacet   `json:"nominalTime"`
	Parent        ParentRunFacet     `json:"parent"`
	Producer      string             `json:"_producer"`
	SchemaURL     string             `json:"_schemaURL"`
}

type DocumentationFacet struct {
	Description string `json:"description"`
	Producer    string `json:"_producer"`
	SchemaURL   string `json:"_schemaURL"`
}

type Owner struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type OwnershipFacet struct {
	Owners    []Owner `json:"owners"`
	Producer  string  `json:"_producer"`
	SchemaURL string  `json:"_schemaURL"`
}

type SourceCodeFacet struct {
	Language   string `json:"language"`
	SourceCode string `json:"sourceCode"`
	Producer   string `json:"_producer"`
	SchemaURL  string `json:"_schemaURL"`
}

type SourceCodeLocationFacet struct {
	Type      string `json:"type"`
	URL       string `json:"url"`
	RepoURL   string `json:"repoUrl"`
	Path      string `json:"path"`
	Version   string `json:"version"`
	Tag       string `json:"tag"`
	Producer  string `json:"_producer"`
	SchemaURL string `json:"_schemaURL"`
}

type SQLJobFacet struct {
	Query     string `json:"query"`
	Producer  string `json:"_producer"`
	SchemaURL string `json:"_schemaURL"`
}

type JobFacets struct {
	Documentation      DocumentationFacet      `json:"documentation"`
	Ownership          OwnershipFacet          `json:"ownership"`
	SourceCode         SourceCodeFacet         `json:"sourceCode"`
	SourceCodeLocation SourceCodeLocationFacet `json:"sourceCodeLocation"`
	SQLJob             SQLJobFacet             `json:"sql"`
}

type InputField struct {
	Namespace                 string `json:"namespace"`
	Name                      string `json:"name"`
	Field                     string `json:"field"`
	TransformationType        string `json:"transformationType"`
	TransformationDescription string `json:"transformationDescription"`
}

type InputFields struct {
	InputFields []InputField `json:"inputFields"`
}

type Fields map[string]InputFields

type ColumnLineageFacet struct {
	Fields    Fields `json:"fields"`
	Producer  string `json:"_producer"`
	SchemaURL string `json:"_schemaURL"`
}

type SchemaField struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type SchemaFacet struct {
	Fields    []SchemaField `json:"fields"`
	Producer  string        `json:"_producer"`
	SchemaURL string        `json:"_schemaURL"`
}

type Assertion struct {
	Assertion string `json:"assertion"`
	Success   bool   `json:"success"`
	Column    string `json:"column"`
}

type DataQualityAssertionsFacet struct {
	Assertions []Assertion `json:"assertions"`
	Producer   string      `json:"_producer"`
	SchemaURL  string      `json:"_schemaURL"`
}

type DatasourceFacet struct {
	Name      string `json:"name"`
	URL       string `json:"uri"`
	Producer  string `json:"_producer"`
	SchemaURL string `json:"_schemaURL"`
}

type StorageFacet struct {
	StorageLayer string `json:"storageLayer"`
	FileFormat   string `json:"fileFormat"`
	Producer     string `json:"_producer"`
	SchemaURL    string `json:"_schemaURL"`
}

type Identifier struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Type      string `json:"type"`
}

type LifecycleStateChangeFacet struct {
	LifecycleStateChange string     `json:"lifecycleStateChange"`
	PreviousIdentifier   Identifier `json:"previousIdentifier"`
	Producer             string     `json:"_producer"`
	SchemaURL            string     `json:"_schemaURL"`
}

type SymlinksFacet struct {
	Identifiers []Identifier `json:"identifiers"`
	Producer    string       `json:"_producer"`
	SchemaURL   string       `json:"_schemaURL"`
}

type DatasetFacets struct {
	Datasource            DatasourceFacet            `json:"dataSource"`
	DataQualityAssertions DataQualityAssertionsFacet `json:"dataQualityAssertions"`
	LifecycleStateChange  LifecycleStateChangeFacet  `json:"lifecycleStateChange"`
	Schema                SchemaFacet                `json:"schema"`
	Storage               StorageFacet               `json:"storage"`
	Ownership             OwnershipFacet             `json:"ownership"`
	ColumnLineage         ColumnLineageFacet         `json:"columnLineage"`
	Symlinks              SymlinksFacet              `json:"symlinks"`
}

func NewDatasetFacets() *DatasetFacets {
	return &DatasetFacets{
		ColumnLineage: ColumnLineageFacet{Fields: make(Fields)},
	}
}

type QuantileMap map[string]int64

type ColumnMetric struct {
	NullCount     int64       `json:"nullCount"`
	DistinctCount int64       `json:"distinctCount"`
	Sum           int64       `json:"sum"`
	Count         int64       `json:"count"`
	Min           int64       `json:"min"`
	Max           int64       `json:"max"`
	Quantiles     QuantileMap `json:"quantiles"`
}

type ColumnMetricMap map[string]ColumnMetric

type DataQualityMetricsFacet struct {
	RowCount      int64           `json:"rowCount"`
	Size          int64           `json:"size"`
	ColumnMetrics ColumnMetricMap `json:"columnMetrics"`
	Producer      string          `json:"_producer"`
	SchemaURL     string          `json:"_schemaURL"`
}

type InputDatasetFacets struct {
	DataQualityMetrics DataQualityMetricsFacet `json:"dataQualityMetrics"`
}

func NewInputDatasetFacets() *InputDatasetFacets {
	return &InputDatasetFacets{
		DataQualityMetrics: DataQualityMetricsFacet{
			ColumnMetrics: make(ColumnMetricMap),
		},
	}
}

type OutputStatisticsFacet struct {
	RowCount  int64  `json:"rowCount"`
	Size      int64  `json:"size"`
	Producer  string `json:"_producer"`
	SchemaURL string `json:"_schemaURL"`
}

type OutputDatasetFacets struct {
	OutputStatistics OutputStatisticsFacet `json:"outputStatistics"`
}
