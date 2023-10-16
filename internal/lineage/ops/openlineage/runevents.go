package openlineage

import (
	"oplin/internal/lineage"
	"oplin/internal/lineage/db"
	"oplin/internal/openlineage"
	"oplin/internal/utils"
	"context"
	"encoding/json"
	"time"

	"github.com/rotisserie/eris"
)

func createRunEvent(
	ctx context.Context, qtx *db.Queries, runID int64, eventTime time.Time, eventType string,
	msg json.RawMessage,
) (*db.LineageRunEvent, error) {
	t, err := lineage.RunEventTypeFromString(eventType)
	if err != nil {
		return nil, eris.Wrapf(err, "could not convert event type[%s]", eventType)
	}

	params := db.CreateRunEventParams{
		RunID:     runID,
		EventType: int32(t),
		EventTime: eventTime,
		Facets:    utils.ToPQRawMessageType(msg),
		CreatedAt: utils.NowUTC(),
	}
	runEvent, err := qtx.CreateRunEvent(ctx, params)
	if err != nil {
		return nil, eris.Wrapf(err, "create run event[%v] failed", params)
	}
	return &runEvent, nil
}

func saveRequest(ctx context.Context, qtx *db.Queries, ev *openlineage.RunEvent) error {
	msg, err := json.Marshal(ev)
	if err != nil {
		return eris.Wrap(err, "could not marshal request")
	}
	params := db.CreateRequestParams{
		Payload:   msg,
		CreatedAt: utils.NowUTC(),
	}
	_, err = qtx.CreateRequest(ctx, params)
	if err != nil {
		return eris.Wrapf(err, "create request[%v] failed", params)
	}

	return nil
}

func CreateWithOpenLineageRunEvent(ctx context.Context, deps Deps, ev *openlineage.RunEvent) (*lineage.RunEvent, error) {
	pg := deps.GetDB()
	tx, err := pg.Begin()
	if err != nil {
		return nil, eris.Wrap(err, "begin transaction failed")
	}
	defer tx.Rollback()
	qtx := db.New(tx).WithTx(tx)

	err = saveRequest(ctx, qtx, ev)
	if err != nil {
		return nil, eris.Wrap(err, "could not save request")
	}

	ns, err := createJobNamespaceIfNotExists(ctx, qtx, ev.Job.Namespace)
	if err != nil {
		return nil, err
	}

	job, err := createJobIfNotExists(ctx, qtx, ns.ID, ev.Job.Name, ev.Job.Facets)
	if err != nil {
		return nil, err
	}

	jobVersion, err := createJobVersionIfNotExists(ctx, qtx, job)
	if err != nil {
		return nil, err
	}

	if jobVersion.ID != job.CurrentVersionID.Int64 {
		_, err := qtx.UpdateCurrentJobVersion(ctx, db.UpdateCurrentJobVersionParams{
			CurrentVersionID: utils.NullInt64(&jobVersion.ID),
			UpdatedAt:        utils.NowUTCAsNullTime(),
			ID:               job.ID,
		})
		if err != nil {
			return nil, eris.Wrapf(err, "update current job version[%v] failed", jobVersion.ID)
		}
	}

	run, err := createRunIfNotExists(ctx, qtx, ev.Run.ID, ev.EventTime, jobVersion.ID, ev.Run.Facets)
	if err != nil {
		return nil, err
	}

	runEvent, err := createRunEvent(ctx, qtx, run.ID, ev.EventTime, ev.EventType, ev.Run.Facets)
	if err != nil {
		return nil, err
	}

	run, err = updateRun(ctx, qtx, run, runEvent)
	if err != nil {
		return nil, err
	}

	if runEvent.EventType == int32(lineage.RunEventTypeComplete) {
		for _, dsInput := range ev.Inputs {
			dsIO := IODataset{
				Dataset:  dsInput.Dataset,
				IOFacets: dsInput.InputFacets,
				Type:     lineage.IOTypeInput,
			}
			_, err := handleIO(ctx, qtx, dsIO, runEvent)
			if err != nil {
				return nil, err
			}
		}
		for _, dsOutput := range ev.Outputs {
			dsIO := IODataset{
				Dataset:  dsOutput.Dataset,
				IOFacets: dsOutput.OutputFacets,
				Type:     lineage.IOTypeOutput,
			}
			_, err := handleIO(ctx, qtx, dsIO, runEvent)
			if err != nil {
				return nil, err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &lineage.RunEvent{
		ID:        runEvent.ID,
		EventType: lineage.RunEventType(runEvent.EventType),
		EventTime: runEvent.EventTime,
		RunID:     run.ID,
		CreatedAt: runEvent.CreatedAt,
		UpdatedAt: runEvent.UpdatedAt.Time,
	}, nil
}
