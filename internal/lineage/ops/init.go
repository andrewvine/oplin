package ops

import (
	"oplin/internal/lineage/db"
	"context"

	"github.com/rotisserie/eris"
)

func InitializeDB(ctx context.Context, deps Deps) error {
	err := db.InitializeDB(ctx, deps.GetDB())
	if err != nil {
		return eris.Wrap(err, "Failed to initialize db")
	}
	return nil
}
