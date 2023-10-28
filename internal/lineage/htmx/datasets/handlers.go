package datasets

import (
	"context"
	"fmt"
	"net/http"
	"oplin/internal/lineage"
	"oplin/internal/lineage/htmx"
	"oplin/internal/lineage/ops"
	"oplin/internal/openlineage"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FieldLineage struct {
	OutputNamespace           string
	OutputName                string
	OutputField               string
	InputNamespace            string
	InputName                 string
	InputField                string
	TransformationType        string
	TransformationDescription string
}

type Field struct {
	Name string
	ID   string
}

type Table struct {
	Name   string
	Fields []Field
}

type Line struct {
	Input  string
	Output string
}

type TabItem struct {
	Key  string
	Text string
	Href string
	Role string
	Icon string
}

type Breadcrumb struct {
	Text  string
	Href  string
	Class string
}

var TabItems = []TabItem{
	{Key: "fields", Text: "Fields", Href: "/lineage/datasets/%d/fields"},
	{Key: "lineage", Text: "Lineage", Href: "/lineage/datasets/%d/lineage"},
	{Key: "ownership", Text: "Ownership", Href: "/lineage/datasets/%d/ownership"},
	{Key: "quality", Text: "Quality", Href: "/lineage/datasets/%d/quality"},
	{Key: "more", Text: "More...", Href: "/lineage/datasets/%d/more"},
}

func buildDatasetBreadcrumbs(dsID int64, text string) []Breadcrumb {
	return []Breadcrumb{
		{Href: "/lineage/datasets", Text: "Datasets"},
		{Href: fmt.Sprintf("/lineage/datasets/%d", dsID), Text: text, Class: "is-active"},
	}
}

func buildTabItems(chosenKey string, dsID int64) []TabItem {
	var res []TabItem
	for _, it := range TabItems {
		it.Href = fmt.Sprintf(it.Href, dsID)
		if it.Key == chosenKey {
			it.Role = "button"
		}
		res = append(res, it)
	}
	return res
}

func buildTablesAndLines(fl []FieldLineage) ([]Table, []Table, []Line) {
	var inputMap = make(map[string]Table)
	var outputMap = make(map[string]Table)
	var lines []Line

	for i, l := range fl {
		inName := fmt.Sprintf("%s.%s", l.InputNamespace, l.InputName)
		t, ok := inputMap[inName]
		if !ok {
			t = Table{Name: inName}
		}
		inID := fmt.Sprintf("in.%d", i)
		t.Fields = append(t.Fields, Field{Name: l.InputField, ID: inID})
		inputMap[inName] = t

		outName := fmt.Sprintf("%s.%s", l.OutputNamespace, l.OutputName)
		t, ok = outputMap[outName]
		if !ok {
			t = Table{Name: outName}
		}
		outID := fmt.Sprintf("out.%d", i)
		t.Fields = append(t.Fields, Field{Name: l.OutputField, ID: outID})
		outputMap[outName] = t

		lines = append(lines, Line{Input: inID, Output: outID})
	}

	var inputTables []Table
	for _, t := range inputMap {
		inputTables = append(inputTables, t)
	}
	var outputTables []Table
	for _, t := range outputMap {
		outputTables = append(outputTables, t)
	}

	return inputTables, outputTables, lines
}

func MakeListDatasets(deps htmx.Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		dss, err := ops.ListDatasetsWithNamespaces(ctx, deps)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}
		c.HTML(http.StatusOK, "lineage/datasets-list.html", gin.H{
			"Title":     "Datasets",
			"Datasets":  dss,
			"MenuItems": htmx.BuildMenuItems("datasets"),
		})
	}
}

func MakeGetDataset(deps htmx.Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		s := c.Param("id")
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}

		ds, err := ops.GetDatasetWithNamespace(ctx, deps, id)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}
		fields, err := ops.ListFieldsForDatasetVersion(ctx, deps, ds.Dataset.CurrentVersionID)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}
		vs, err := ops.ListDatasetVersions(ctx, deps, ds.Dataset.ID)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}

		title := fmt.Sprintf("%s %s", ds.DatasetNamespace.Name, ds.Dataset.Name)

		c.HTML(http.StatusOK, "lineage/datasets-detail.html", gin.H{
			"Title":                title,
			"DatasetWithNamespace": ds,
			"Fields":               fields,
			"VersionID":            ds.Dataset.CurrentVersionID,
			"Versions":             vs,
			"TabItems":             buildTabItems("fields", ds.Dataset.ID),
			"MenuItems":            htmx.BuildMenuItems("datasets"),
			"Breadcrumbs":          buildDatasetBreadcrumbs(ds.Dataset.ID, title),
		})
	}
}

func MakeGetDatasetFields(deps htmx.Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		s := c.Param("id")
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}

		ds, err := ops.GetDatasetWithNamespace(ctx, deps, id)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}
		fields, err := ops.ListFieldsForDatasetVersion(ctx, deps, ds.Dataset.CurrentVersionID)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}
		vs, err := ops.ListDatasetVersions(ctx, deps, ds.Dataset.ID)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}
		c.HTML(http.StatusOK, "lineage/datasets-fields.html", gin.H{
			"Fields":    fields,
			"VersionID": ds.Dataset.CurrentVersionID,
			"Versions":  vs,
			"TabItems":  buildTabItems("fields", ds.Dataset.ID),
		})
	}
}

func MakeGetDatasetVersionFields(deps htmx.Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		s := c.Query("version")
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}

		dsv, err := ops.GetDatasetVersionByID(ctx, deps, id)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}
		fields, err := ops.ListFieldsForDatasetVersion(ctx, deps, dsv.ID)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}
		vs, err := ops.ListDatasetVersions(ctx, deps, dsv.DatasetID)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}
		c.HTML(http.StatusOK, "lineage/datasets-fields.html", gin.H{
			"Fields":    fields,
			"VersionID": dsv.ID,
			"Versions":  vs,
			"TabItems":  buildTabItems("fields", dsv.DatasetID),
		})
	}
}

func MakeGetDatasetLineage(deps htmx.Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		s := c.Param("id")
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}

		ds, err := ops.GetDatasetWithNamespace(ctx, deps, id)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}
		vs, err := ops.ListDatasetVersions(ctx, deps, id)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}
		f, err := ops.GetLatestFacetsByDatasetVersionID(ctx, deps, ds.Dataset.CurrentVersionID)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}
		fl := buildFieldLineages(f, ds)
		inputTables, outputTables, lines := buildTablesAndLines(fl)

		c.HTML(http.StatusOK, "lineage/datasets-lineage.html", gin.H{
			"InputTables":  inputTables,
			"OutputTables": outputTables,
			"Lines":        lines,
			"VersionID":    vs[0].ID,
			"Versions":     vs,
			"TabItems":     buildTabItems("lineage", ds.Dataset.ID),
		})
	}
}

func MakeGetDatasetVersionLineage(deps htmx.Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		s := c.Query("version")
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}

		dsv, err := ops.GetDatasetVersionByID(ctx, deps, id)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}
		ds, err := ops.GetDatasetWithNamespace(ctx, deps, dsv.DatasetID)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}
		vs, err := ops.ListDatasetVersions(ctx, deps, dsv.DatasetID)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}
		f, err := ops.GetLatestFacetsByDatasetVersionID(ctx, deps, dsv.ID)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}

		fl := buildFieldLineages(f, ds)
		inputTables, outputTables, lines := buildTablesAndLines(fl)

		c.HTML(http.StatusOK, "lineage/datasets-lineage.html", gin.H{
			"InputTables":  inputTables,
			"OutputTables": outputTables,
			"Lines":        lines,
			"VersionID":    dsv.ID,
			"Versions":     vs,
			"TabItems":     buildTabItems("lineage", ds.Dataset.ID),
		})
	}
}

func buildFieldLineages(f *openlineage.DatasetFacets, ds *lineage.DatasetWithNamespace) []FieldLineage {
	var fl []FieldLineage
	for name, fields := range f.ColumnLineage.Fields {
		for _, inputField := range fields.InputFields {
			fl = append(fl, FieldLineage{
				OutputNamespace:           ds.DatasetNamespace.Name,
				OutputName:                ds.Dataset.Name,
				OutputField:               name,
				InputNamespace:            inputField.Namespace,
				InputName:                 inputField.Name,
				InputField:                inputField.Field,
				TransformationType:        inputField.TransformationType,
				TransformationDescription: inputField.TransformationDescription,
			})
		}
	}
	return fl
}

func MakeGetDatasetOwnership(deps htmx.Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		s := c.Param("id")
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}

		ds, err := ops.GetDatasetWithNamespace(ctx, deps, id)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}
		c.HTML(http.StatusOK, "lineage/datasets-ownership.html", gin.H{
			"DatasetWithNamespace": ds,
			"TabItems":             buildTabItems("ownership", ds.Dataset.ID),
		})
	}
}

func MakeGetDatasetQuality(deps htmx.Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		s := c.Param("id")
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}

		ds, err := ops.GetDatasetWithNamespace(ctx, deps, id)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}
		c.HTML(http.StatusOK, "lineage/datasets-quality.html", gin.H{
			"DatasetWithNamespace": ds,
			"TabItems":             buildTabItems("quality", ds.Dataset.ID),
		})
	}
}

func MakeGetDatasetMore(deps htmx.Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		s := c.Param("id")
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}

		ds, err := ops.GetDatasetWithNamespace(ctx, deps, id)
		if err != nil {
			htmx.InternalServerError(c, err)
			return
		}
		c.HTML(http.StatusOK, "lineage/datasets-more.html", gin.H{
			"DatasetWithNamespace": ds,
			"TabItems":             buildTabItems("more", ds.Dataset.ID),
		})
	}
}
