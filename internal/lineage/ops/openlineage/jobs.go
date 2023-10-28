package openlineage

import (
	"context"
	"encoding/json"
	"oplin/internal/lineage/db"
	"oplin/internal/utils"

	"github.com/rotisserie/eris"
)

func createJobNamespaceIfNotExists(ctx context.Context, qtx *db.Queries, name string) (*db.LineageJobNamespace, error) {
	ns, err := qtx.GetJobNamespaceByName(ctx, name)
	if err != nil && !utils.IsNoRowsError(err) {
		return nil, eris.Wrapf(err, "getting job namespace[%s] failed", name)
	}
	if err == nil {
		return &ns, nil
	}

	ns, err = qtx.CreateJobNamespace(ctx, db.CreateJobNamespaceParams{
		Name:      name,
		CreatedAt: utils.NowUTC(),
	})
	if err != nil {
		return nil, eris.Wrapf(err, "creating job namespace failed[%s]", name)
	}
	return &ns, err
}

func createJobIfNotExists(
	ctx context.Context, qtx *db.Queries, nsID int64, jobName string, msg json.RawMessage,
) (*db.LineageJob, error) {
	getParams := db.GetJobByNamespaceIDAndNameParams{
		NamespaceID: nsID, Name: jobName,
	}
	job, err := qtx.GetJobByNamespaceIDAndName(ctx, getParams)
	if err != nil && !utils.IsNoRowsError(err) {
		return nil, eris.Wrapf(err, "get job by params[%v] failed", getParams)
	}
	if err == nil {
		return &job, nil
	}

	params := db.CreateJobParams{
		NamespaceID: nsID,
		Name:        jobName,
		Facets:      utils.ToPQRawMessageType(msg),
		CreatedAt:   utils.NowUTC(),
	}
	job, err = qtx.CreateJob(ctx, params)
	if err != nil {
		return nil, eris.Wrapf(err, "create job[%v] failed", params)
	}
	return &job, nil
}

func createJobVersionIfNotExists(
	ctx context.Context, qtx *db.Queries, job *db.LineageJob,
) (*db.LineageJobVersion, error) {

	if job.CurrentVersionID.Valid {
		jv, err := qtx.GetJobVersionByID(ctx, job.CurrentVersionID.Int64)
		if err != nil {
			return nil, eris.Wrapf(err, "get job version[%v] failed", job.CurrentVersionID.Int64)
		}
		return &jv, err
	}

	params := db.CreateJobVersionParams{
		NamespaceID: job.NamespaceID,
		JobID:       job.ID,
		Name:        job.Name,
		Facets:      job.Facets,
		CreatedAt:   utils.NowUTC(),
	}
	jv, err := qtx.CreateJobVersion(ctx, params)
	if err != nil {
		return nil, eris.Wrapf(err, "create job version[%v] failed", params)
	}
	return &jv, nil
}
