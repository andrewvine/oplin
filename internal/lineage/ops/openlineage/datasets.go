package openlineage

import (
	"bytes"
	"context"
	"encoding/json"
	"oplin/internal/lineage"
	"oplin/internal/lineage/db"
	"oplin/internal/openlineage"
	"oplin/internal/utils"
	"sort"

	"github.com/rotisserie/eris"
)

type IODataset struct {
	openlineage.Dataset
	IOFacets json.RawMessage
	Type     lineage.IOType
}

func createDatasetNamespaceIfNotExists(ctx context.Context, qtx *db.Queries, name string) (*db.LineageDatasetNamespace, error) {
	ns, err := qtx.GetDatasetNamespaceByName(ctx, name)
	if err != nil && !utils.IsNoRowsError(err) {
		return nil, eris.Wrapf(err, "getting dataset namespace[%s] failed", name)
	}
	if err == nil {
		return &ns, nil
	}

	ns, err = qtx.CreateDatasetNamespace(ctx, db.CreateDatasetNamespaceParams{
		Name:      name,
		CreatedAt: utils.NowUTC(),
	})
	if err != nil {
		return nil, eris.Wrapf(err, "creating dataset namespace failed[%s]", name)
	}
	return &ns, err
}

func createDatasetIfNotExists(
	ctx context.Context, qtx *db.Queries, nsID int64, datasetName string, msg json.RawMessage,
) (*db.LineageDataset, error) {
	getParams := db.GetDatasetByNamespaceIDAndNameParams{
		NamespaceID: nsID, Name: datasetName,
	}
	ds, err := qtx.GetDatasetByNamespaceIDAndName(ctx, getParams)
	if err != nil && !utils.IsNoRowsError(err) {
		return nil, eris.Wrapf(err, "get dataset by params[%v] failed", getParams)
	}
	if err == nil {
		return &ds, nil
	}

	params := db.CreateDatasetParams{
		NamespaceID: nsID,
		Name:        datasetName,
		Facets:      utils.ToPQRawMessageType(msg),
		CreatedAt:   utils.NowUTC(),
	}
	ds, err = qtx.CreateDataset(ctx, params)
	if err != nil {
		return nil, eris.Wrapf(err, "create dataset[%v] failed", params)
	}
	return &ds, nil
}

func updateDataset(
	ctx context.Context, qtx *db.Queries, dsID int64, msg json.RawMessage,
) (*db.LineageDataset, error) {
	params := db.UpdateDatasetParams{
		ID:        dsID,
		Facets:    utils.ToPQRawMessageType(msg),
		UpdatedAt: utils.NowUTCAsNullTime(),
	}
	ds, err := qtx.UpdateDataset(ctx, params)
	if err != nil {
		return nil, eris.Wrapf(err, "update dataset[%v] failed", params)
	}
	return &ds, nil
}

func createDatasetVersionIfNotExists(
	ctx context.Context, qtx *db.Queries, ds *db.LineageDataset,
) (*db.LineageDatasetVersion, error) {

	if ds.CurrentVersionID.Valid {
		dv, err := qtx.GetDatasetVersionByID(ctx, ds.CurrentVersionID.Int64)
		if err != nil {
			return nil, eris.Wrapf(err, "get dataset version[%v] failed", ds.CurrentVersionID.Int64)
		}
		return &dv, err
	}

	return createDatasetVersion(ctx, qtx, ds)
}

func createDatasetVersion(
	ctx context.Context, qtx *db.Queries, ds *db.LineageDataset,
) (*db.LineageDatasetVersion, error) {

	params := db.CreateDatasetVersionParams{
		NamespaceID: ds.NamespaceID,
		DatasetID:   ds.ID,
		Name:        ds.Name,
		CreatedAt:   utils.NowUTC(),
	}
	dv, err := qtx.CreateDatasetVersion(ctx, params)
	if err != nil {
		return nil, eris.Wrapf(err, "create dataset version[%v] failed", params)
	}
	return &dv, nil
}

func equalStringSlices(a, b []string) bool {
	sort.Strings(a)
	sort.Strings(b)

	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func fieldNamesEqual(rows []db.LineageField, fields []openlineage.SchemaField) bool {
	var a []string
	var b []string

	for i := range rows {
		a = append(a, rows[i].Name)
	}
	for i := range fields {
		b = append(b, fields[i].Name)
	}
	return equalStringSlices(a, b)
}

func createFields(
	ctx context.Context, qtx *db.Queries, dsvID int64, fields []openlineage.SchemaField,
) ([]db.LineageField, error) {
	var rows []db.LineageField
	for _, f := range fields {
		params := db.CreateFieldParams{
			DatasetVersionID: dsvID,
			Name:             f.Name,
			DataType:         f.Type,
			Description:      utils.NullString(f.Description),
			CreatedAt:        utils.NowUTC(),
		}
		row, err := qtx.CreateField(ctx, params)
		if err != nil {
			return nil, eris.Wrapf(err, "could not create field[%v]", params)
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func createRunDatasetVersionIfNotExists(
	ctx context.Context, qtx *db.Queries, runID int64, versionID int64, ioMsg json.RawMessage, dsMsg json.RawMessage, io lineage.IOType,
) (*db.LineageRunDatasetVersion, error) {

	getParams := db.GetRunDatasetVersionByRunIDAndDatasetVersionIDParams{
		RunID:            runID,
		DatasetVersionID: versionID,
	}
	rdv, err := qtx.GetRunDatasetVersionByRunIDAndDatasetVersionID(ctx, getParams)
	if err != nil && !utils.IsNoRowsError(err) {
		return nil, eris.Wrapf(err, "get run dataset version[%v] failed", getParams)
	}
	if err == nil {
		return &rdv, nil
	}

	params := db.CreateRunDatasetVersionParams{
		RunID:            runID,
		DatasetVersionID: versionID,
		IoType:           int32(io),
		IoFacets:         utils.ToPQRawMessageType(ioMsg),
		DatasetFacets:    utils.ToPQRawMessageType(dsMsg),
		CreatedAt:        utils.NowUTC(),
	}
	rdv, err = qtx.CreateRunDatasetVersion(ctx, params)
	if err != nil {
		return nil, eris.Wrapf(err, "create run dataset version[%v] failed", params)
	}
	return &rdv, nil
}

func updateCurrentDatasetVersion(ctx context.Context, qtx *db.Queries, ds *db.LineageDataset, dsVersion *db.LineageDatasetVersion) (
	*db.LineageDataset, error) {
	row, err := qtx.UpdateCurrentDatasetVersion(ctx, db.UpdateCurrentDatasetVersionParams{
		CurrentVersionID: utils.NullInt64(&dsVersion.ID),
		UpdatedAt:        utils.NowUTCAsNullTime(),
		ID:               ds.ID,
	})
	if err != nil {
		return nil, eris.Wrapf(err, "update current dataset version[%v] failed", dsVersion.ID)
	}
	return &row, nil
}

func handleIO(
	ctx context.Context, qtx *db.Queries, dsIO IODataset, runEvent *db.LineageRunEvent,
) (*db.LineageRunDatasetVersion, error) {
	ns, err := createDatasetNamespaceIfNotExists(ctx, qtx, dsIO.Namespace)
	if err != nil {
		return nil, eris.Wrapf(err, "create dataset namespace[%v] failed", dsIO.Namespace)
	}

	ds, err := createDatasetIfNotExists(ctx, qtx, ns.ID, dsIO.Name, dsIO.Facets)
	if err != nil {
		return nil, err
	}

	// only update when facets change
	if !bytes.Equal(ds.Facets.RawMessage, dsIO.Facets) {
		ds, err = updateDataset(ctx, qtx, ds.ID, dsIO.Facets)
		if err != nil {
			return nil, err
		}
	}

	dsVersion, err := createDatasetVersionIfNotExists(ctx, qtx, ds)
	if err != nil {
		return nil, err
	}

	// set current version to new version
	if ds.CurrentVersionID.Int64 != dsVersion.ID {
		ds, err = updateCurrentDatasetVersion(ctx, qtx, ds, dsVersion)
		if err != nil {
			return nil, err
		}
	}

	fs := openlineage.NewDatasetFacets()
	err = json.Unmarshal(dsIO.Facets, fs)
	if err != nil {
		return nil, eris.Wrapf(err, "could not unmarshall[%s]", dsIO.Facets)
	}

	rows, err := qtx.ListFieldsByDatasetVersionID(ctx, dsVersion.ID)
	if err != nil {
		return nil, eris.Wrapf(err, "Fethching fields failed for dsvID[%d]", dsVersion.ID)
	}

	if len(rows) == 0 {
		rows, err = createFields(ctx, qtx, dsVersion.ID, fs.Schema.Fields)
		if err != nil {
			return nil, err
		}
	}

	// create a new version if the schema has changed
	if !fieldNamesEqual(rows, fs.Schema.Fields) {
		dsVersion, err = createDatasetVersion(ctx, qtx, ds)
		if err != nil {
			return nil, err
		}
		rows, err = createFields(ctx, qtx, dsVersion.ID, fs.Schema.Fields)
		if err != nil {
			return nil, err
		}
		ds, err = updateCurrentDatasetVersion(ctx, qtx, ds, dsVersion)
		if err != nil {
			return nil, err
		}
	}

	return createRunDatasetVersionIfNotExists(ctx, qtx, runEvent.RunID, dsVersion.ID, dsIO.IOFacets, dsIO.Dataset.Facets, dsIO.Type)
}
