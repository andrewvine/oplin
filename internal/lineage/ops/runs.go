package ops

import (
	"oplin/internal/lineage"
	"oplin/internal/lineage/db"
	ol "oplin/internal/openlineage"
	"oplin/internal/utils"
	"context"
	"encoding/json"

	"github.com/rotisserie/eris"
)

func GetRunWithID(ctx context.Context, deps Deps, runID int64) (*lineage.Run, error) {
	pg := deps.GetDB()
	qtx := db.New(pg)
	row, err := qtx.GetRunByID(ctx, runID)
	if utils.IsNoRowsError(err) {
		return nil, eris.Wrapf(err, "could not find run with id[%d]", runID)
	}
	if err != nil {
		return nil, err
	}

	f := &ol.RunFacets{}
	err = json.Unmarshal(row.Facets.RawMessage, f)
	if err != nil {
		return nil, eris.Wrapf(err, "could not unmarshall[%s]", row.Facets.RawMessage)
	}

	return &lineage.Run{
		ID:                  row.ID,
		JobVersionID:        row.JobVersionID,
		Facets:              *f,
		ParentRunID:         row.ParentRunID.Int64,
		LastEventType:       lineage.RunEventType(row.LastEventType),
		NominalStartedAt:    row.NominalStartedAt.Time,
		NominalEndedAt:      row.NominalEndedAt.Time,
		StartedAt:           row.StartedAt.Time,
		EndedAt:             row.EndedAt.Time,
		ErrorMessage:        row.ErrorMessage.String,
		ProgrammingLanguage: row.ProgrammingLanguage.String,
		Stacktrace:          row.Stacktrace.String,
		CreatedAt:           row.CreatedAt,
		UpdatedAt:           row.UpdatedAt.Time,
	}, nil
}

func ListRunsByJobVersionID(ctx context.Context, deps Deps, jvID int64) ([]lineage.Run, error) {
	pg := deps.GetDB()
	qtx := db.New(pg)
	rows, err := qtx.ListRunsByJobVersionID(ctx, jvID)
	if err != nil {
		return nil, eris.Wrap(err, "Failed to list dataset versions")
	}
	var res []lineage.Run

	for _, row := range rows {
		res = append(res, lineage.Run{
			ID:           row.ID,
			JobVersionID: row.JobVersionID,
			//Facets:              row.Facets.RawMessage,
			ParentRunID:         row.ParentRunID.Int64,
			LastEventType:       lineage.RunEventType(row.LastEventType),
			NominalStartedAt:    row.NominalStartedAt.Time,
			NominalEndedAt:      row.NominalEndedAt.Time,
			StartedAt:           row.StartedAt.Time,
			EndedAt:             row.EndedAt.Time,
			ErrorMessage:        row.ErrorMessage.String,
			ProgrammingLanguage: row.ProgrammingLanguage.String,
			Stacktrace:          row.Stacktrace.String,
			CreatedAt:           row.CreatedAt,
			UpdatedAt:           row.UpdatedAt.Time,
		})
	}
	return res, nil
}

func ListRunEventsByRunID(ctx context.Context, deps Deps, runID int64) ([]lineage.RunEvent, error) {
	pg := deps.GetDB()
	qtx := db.New(pg)
	rows, err := qtx.ListRunEventsByRunID(ctx, runID)
	if err != nil {
		return nil, eris.Wrap(err, "Failed to list run events")
	}
	var res []lineage.RunEvent

	for _, row := range rows {
		res = append(res, lineage.RunEvent{
			ID:        row.ID,
			RunID:     runID,
			EventType: lineage.RunEventType(row.EventType),
			EventTime: row.EventTime,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt.Time,
		})
	}
	return res, nil
}

func ListRunDatasetVersionsWithRelationshipsByRunID(ctx context.Context, deps Deps, id int64) ([]lineage.RunIODatasetWithRelationships, error) {
	pg := deps.GetDB()
	qtx := db.New(pg)
	rows, err := qtx.ListRunDatasetVersionsWithRelationshipsByRunID(ctx, id)
	if err != nil {
		return nil, eris.Wrap(err, "Failed to get dataset")
	}

	var res []lineage.RunIODatasetWithRelationships

	for _, row := range rows {

		inFacets := ol.NewInputDatasetFacets()
		outFacets := &ol.OutputDatasetFacets{}
		if len(row.IoFacets.RawMessage) > 0 {
			if row.IoType == int32(lineage.IOTypeInput) {
				err = json.Unmarshal(row.IoFacets.RawMessage, inFacets)
				if err != nil {
					return nil, eris.Wrapf(err, "could not unmarshall[%s]", row.IoFacets.RawMessage)
				}
			} else {
				err = json.Unmarshal(row.IoFacets.RawMessage, outFacets)
				if err != nil {
					return nil, eris.Wrapf(err, "could not unmarshall[%s]", row.IoFacets.RawMessage)
				}
			}
		}

		res = append(res, lineage.RunIODatasetWithRelationships{
			RunIODataset: lineage.RunIODataset{
				DatasetVersionID: row.VersionID,
				RunID:            row.RunID,
				IOType:           lineage.IOType(row.IoType),
				InputFacets:      *inFacets,
				OutputFacets:     *outFacets,
				CreatedAt:        row.CreatedAt,
			},
			DatasetVersion: lineage.DatasetVersion{
				ID:                 row.VersionID,
				DatasetID:          row.VersionDatasetID,
				DatasetNamespaceID: row.VersionNamespaceID,
				Name:               row.VersionName,
				CreatedAt:          row.VersionCreatedAt,
				UpdatedAt:          row.VersionUpdatedAt.Time,
			},
			DatasetNamespace: lineage.DatasetNamespace{
				ID:        row.NamespaceID,
				Name:      row.NamespaceName,
				CreatedAt: row.NamespaceCreatedAt,
				UpdatedAt: row.NamespaceUpdatedAt.Time,
			},
		})
	}
	return res, nil
}
