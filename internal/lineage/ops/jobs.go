package ops

import (
	"oplin/internal/lineage"
	"oplin/internal/lineage/db"
	ol "oplin/internal/openlineage"
	"context"
	"encoding/json"

	"github.com/rotisserie/eris"
)

func ListJobsWithNamespaces(ctx context.Context, deps Deps) ([]lineage.JobWithNamespace, error) {
	pg := deps.GetDB()
	qtx := db.New(pg)
	rows, err := qtx.ListJobsWithNamespaces(ctx)
	if err != nil {
		return nil, eris.Wrap(err, "Failed to list jobs")
	}
	var res []lineage.JobWithNamespace

	for _, row := range rows {
		res = append(res, lineage.JobWithNamespace{
			Job: lineage.Job{
				ID:             row.ID,
				JobNamespaceID: row.NamespaceID,
				Name:           row.Name,
				CreatedAt:      row.CreatedAt,
				UpdatedAt:      row.UpdatedAt.Time},
			JobNamespace: lineage.JobNamespace{
				ID:        row.NamespaceID,
				Name:      row.NamespaceName,
				CreatedAt: row.NamespaceCreatedAt,
				UpdatedAt: row.NamespaceUpdatedAt.Time},
		})
	}
	return res, nil
}

func ListJobVersions(ctx context.Context, deps Deps, jobID int64) ([]lineage.JobVersion, error) {
	pg := deps.GetDB()
	qtx := db.New(pg)
	rows, err := qtx.ListJobVersionsByJobID(ctx, jobID)
	if err != nil {
		return nil, eris.Wrap(err, "Failed to list dataset versions")
	}
	var res []lineage.JobVersion

	for _, row := range rows {
		res = append(res, lineage.JobVersion{
			ID:             row.ID,
			JobID:          row.JobID,
			JobNamespaceID: row.NamespaceID,
			Name:           row.Name,
			CreatedAt:      row.CreatedAt,
			UpdatedAt:      row.UpdatedAt.Time,
		})
	}
	return res, nil
}

func GetJobWithNamespace(ctx context.Context, deps Deps, id int64) (*lineage.JobWithNamespace, error) {
	pg := deps.GetDB()
	qtx := db.New(pg)
	row, err := qtx.GetJobWithNamespace(ctx, id)
	if err != nil {
		return nil, eris.Wrap(err, "Failed to get job")
	}

	f := &ol.JobFacets{}
	err = json.Unmarshal(row.Facets.RawMessage, f)
	if err != nil {
		return nil, eris.Wrapf(err, "could not unmarshall[%s]", row.Facets.RawMessage)
	}

	res := lineage.JobWithNamespace{
		Job: lineage.Job{
			ID:               row.ID,
			Facets:           *f,
			CurrentVersionID: row.CurrentVersionID.Int64,
			JobNamespaceID:   row.NamespaceID,
			Name:             row.Name,
			CreatedAt:        row.CreatedAt,
			UpdatedAt:        row.UpdatedAt.Time},
		JobNamespace: lineage.JobNamespace{
			ID:        row.NamespaceID,
			Name:      row.NamespaceName,
			CreatedAt: row.NamespaceCreatedAt,
			UpdatedAt: row.NamespaceUpdatedAt.Time},
	}
	return &res, nil
}

func GetJobVersionByID(ctx context.Context, deps Deps, id int64) (*lineage.JobVersion, error) {
	pg := deps.GetDB()
	qtx := db.New(pg)
	row, err := qtx.GetJobVersionByID(ctx, id)
	if err != nil {
		return nil, eris.Wrap(err, "Failed to get job")
	}
	res := lineage.JobVersion{
		ID:             row.ID,
		JobID:          row.JobID,
		JobNamespaceID: row.NamespaceID,
		Name:           row.Name,
		CreatedAt:      row.CreatedAt,
		UpdatedAt:      row.UpdatedAt.Time,
	}
	return &res, nil
}
