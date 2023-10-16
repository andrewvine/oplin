package openlineage

import (
	"oplin/internal/lineage"
	"oplin/internal/lineage/db"
	"oplin/internal/openlineage"
	"oplin/internal/utils"
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/rotisserie/eris"
)

func createRunIfNotExists(
	ctx context.Context, qtx *db.Queries, runUUID uuid.UUID, eventTime time.Time, jobVersionID int64, msg json.RawMessage,
) (*db.LineageRun, error) {
	run, err := qtx.GetRunByUUID(ctx, runUUID)
	if err != nil && !utils.IsNoRowsError(err) {
		return nil, eris.Wrapf(err, "get run[%s] failed", runUUID)
	}
	if err == nil {
		return &run, nil
	}

	fs := openlineage.RunFacets{}

	if len(msg) > 0 {
		err = json.Unmarshal(msg, &fs)
		if err != nil {
			return nil, eris.Wrapf(err, "unmarshalling facets[%v] failed", msg)
		}
	}

	var parentRunID *int64
	if len(fs.Parent.ParentRun.RunID) > 0 {
		u, err := uuid.Parse(fs.Parent.ParentRun.RunID)
		if err != nil {
			return nil, eris.Wrapf(err, "parent run id[%v] conversion to uuid failed", fs.Parent.ParentRun.RunID)
		}
		parent, err := qtx.GetRunByUUID(ctx, u)
		if err != nil {
			return nil, eris.Wrapf(err, "get parent by run uuid[%v] failed", u)
		}
		parentRunID = &parent.ID
	}

	params := db.CreateRunParams{
		RunUuid:          runUUID,
		JobVersionID:     jobVersionID,
		Facets:           utils.ToPQRawMessageType(msg),
		ParentRunID:      utils.NullInt64(parentRunID),
		NominalStartedAt: utils.NullTime(fs.NominalTime.StartTime),
		NominalEndedAt:   utils.NullTime(fs.NominalTime.EndTime),
		StartedAt:        utils.NullTime(eventTime),
		CreatedAt:        utils.NowUTC(),
	}
	run, err = qtx.CreateRun(ctx, params)
	if err != nil {
		return nil, eris.Wrapf(err, "create run[%v] failed", params)
	}
	return &run, nil
}

func updateRun(
	ctx context.Context, qtx *db.Queries, run *db.LineageRun, runEvent *db.LineageRunEvent,
) (*db.LineageRun, error) {
	endedAt := run.EndedAt.Time
	if runEvent.EventType == int32(lineage.RunEventTypeComplete) {
		endedAt = runEvent.EventTime
	}

	msg, err := utils.MergeFacets(run.Facets.RawMessage, runEvent.Facets.RawMessage)
	if err != nil {
		return nil, eris.Wrapf(err, "update run failed cannot merge facets[%v], [%v]", run.Facets.RawMessage, runEvent.Facets.RawMessage)
	}

	fs := openlineage.RunFacets{}
	if len(msg) > 0 {
		err = json.Unmarshal(msg, &fs)
		if err != nil {
			return nil, eris.Wrapf(err, "unmarshalling facets[%v] failed", msg)
		}
	}

	r, err := qtx.UpdateRun(ctx, db.UpdateRunParams{
		ID:                  run.ID,
		Facets:              utils.ToPQRawMessageType(msg),
		EndedAt:             utils.NullTime(endedAt),
		LastEventType:       runEvent.EventType,
		ErrorMessage:        utils.NullString(fs.ErrorMessage.Message),
		ProgrammingLanguage: utils.NullString(fs.ErrorMessage.ProgrammingLanguage),
		Stacktrace:          utils.NullString(fs.ErrorMessage.Stacktrace),
		UpdatedAt:           utils.NullTime(time.Now()),
	})
	if err != nil {
		return nil, eris.Wrapf(err, "could not update run[%d]", run.ID)
	}
	return &r, nil
}
