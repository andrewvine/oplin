package ops

import (
	"oplin/internal/lineage"
	"oplin/internal/lineage/db"
	ol "oplin/internal/openlineage"
	"context"
	"encoding/json"

	"github.com/rotisserie/eris"
)

func ListDatasetsWithNamespaces(ctx context.Context, deps Deps) ([]lineage.DatasetWithNamespace, error) {
	pg := deps.GetDB()
	qtx := db.New(pg)
	rows, err := qtx.ListDatasetsWithNamespaces(ctx)
	if err != nil {
		return nil, eris.Wrap(err, "Failed to list datasets")
	}
	var res []lineage.DatasetWithNamespace

	for _, row := range rows {
		res = append(res, lineage.DatasetWithNamespace{
			Dataset: lineage.Dataset{
				ID:                 row.ID,
				CurrentVersionID:   row.CurrentVersionID.Int64,
				DatasetNamespaceID: row.NamespaceID,
				Name:               row.Name,
				CreatedAt:          row.CreatedAt,
				UpdatedAt:          row.UpdatedAt.Time},
			DatasetNamespace: lineage.DatasetNamespace{
				ID:        row.NamespaceID,
				Name:      row.NamespaceName,
				CreatedAt: row.NamespaceCreatedAt,
				UpdatedAt: row.NamespaceUpdatedAt.Time},
		})
	}
	return res, nil
}

func ListDatasetVersions(ctx context.Context, deps Deps, dsID int64) ([]lineage.DatasetVersion, error) {
	pg := deps.GetDB()
	qtx := db.New(pg)
	rows, err := qtx.ListDatasetVersionsByDatasetID(ctx, dsID)
	if err != nil {
		return nil, eris.Wrap(err, "Failed to list dataset versions")
	}
	var res []lineage.DatasetVersion

	for _, row := range rows {
		res = append(res, lineage.DatasetVersion{
			ID:                 row.ID,
			DatasetID:          row.DatasetID,
			DatasetNamespaceID: row.NamespaceID,
			Name:               row.Name,
			CreatedAt:          row.CreatedAt,
			UpdatedAt:          row.UpdatedAt.Time,
		})
	}
	return res, nil
}

func ListFieldsForDatasetVersion(ctx context.Context, deps Deps, dsvID int64) ([]lineage.Field, error) {
	pg := deps.GetDB()
	qtx := db.New(pg)
	rows, err := qtx.ListFieldsByDatasetVersionID(ctx, dsvID)
	if err != nil {
		return nil, eris.Wrapf(err, "Failed to list fileds for dataset version[%d]", dsvID)
	}
	var res []lineage.Field

	for _, row := range rows {
		res = append(res, lineage.Field{
			ID:          row.ID,
			Name:        row.Name,
			DataType:    row.DataType,
			Description: row.Description.String,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt.Time,
		})
	}
	return res, nil
}

func GetDatasetWithNamespace(ctx context.Context, deps Deps, id int64) (*lineage.DatasetWithNamespace, error) {
	pg := deps.GetDB()
	qtx := db.New(pg)
	row, err := qtx.GetDatasetWithNamespace(ctx, id)
	if err != nil {
		return nil, eris.Wrap(err, "Failed to get dataset")
	}

	f := ol.NewDatasetFacets()
	err = json.Unmarshal(row.Facets.RawMessage, f)
	if err != nil {
		return nil, eris.Wrapf(err, "could not unmarshall[%s]", row.Facets.RawMessage)
	}

	res := lineage.DatasetWithNamespace{
		Dataset: lineage.Dataset{
			ID:                 row.ID,
			Facets:             *f,
			CurrentVersionID:   row.CurrentVersionID.Int64,
			DatasetNamespaceID: row.NamespaceID,
			Name:               row.Name,
			CreatedAt:          row.CreatedAt,
			UpdatedAt:          row.UpdatedAt.Time},
		DatasetNamespace: lineage.DatasetNamespace{
			ID:        row.NamespaceID,
			Name:      row.NamespaceName,
			CreatedAt: row.NamespaceCreatedAt,
			UpdatedAt: row.NamespaceUpdatedAt.Time},
	}
	return &res, nil
}

func GetLatestFacetsByDatasetVersionID(ctx context.Context, deps Deps, id int64) (*ol.DatasetFacets, error) {
	f := ol.NewDatasetFacets()
	pg := deps.GetDB()
	qtx := db.New(pg)
	row, err := qtx.GetLatestRunDatasetVersionByDatasetVersionID(ctx, id)
	if err != nil {
		return f, eris.Wrap(err, "Failed to get latest dataset version facets")
	}

	err = json.Unmarshal(row.DatasetFacets.RawMessage, f)
	if err != nil {
		return nil, eris.Wrapf(err, "could not unmarshall[%s]", row.DatasetFacets.RawMessage)
	}
	return f, nil
}

func GetDatasetVersionByID(ctx context.Context, deps Deps, id int64) (*lineage.DatasetVersion, error) {
	pg := deps.GetDB()
	qtx := db.New(pg)
	row, err := qtx.GetDatasetVersionByID(ctx, id)
	if err != nil {
		return nil, eris.Wrap(err, "Failed to get latest dataset version facets")
	}
	dsv := lineage.DatasetVersion{
		ID:                 row.ID,
		DatasetID:          row.DatasetID,
		DatasetNamespaceID: row.NamespaceID,
		Name:               row.Name,
		CreatedAt:          row.CreatedAt,
		UpdatedAt:          row.UpdatedAt.Time,
	}

	return &dsv, nil
}
