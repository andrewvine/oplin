package ops

import (
	"context"
	"oplin/internal/lineage"
	"oplin/internal/lineage/db"

	"github.com/rotisserie/eris"
)

func ListRequests(ctx context.Context, deps Deps) ([]lineage.Request, error) {
	pg := deps.GetDB()
	qtx := db.New(pg)
	rows, err := qtx.ListRequests(ctx)
	if err != nil {
		return nil, eris.Wrap(err, "Failed to list requests")
	}
	var res []lineage.Request

	for _, row := range rows {
		res = append(res, lineage.Request{
			ID:        row.ID,
			Payload:   row.Payload,
			CreatedAt: row.CreatedAt,
		})
	}
	return res, nil
}
