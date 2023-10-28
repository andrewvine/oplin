package ops

import (
	"context"
	"oplin/internal/lineage/db"

	"github.com/rotisserie/eris"
)

func InitializeDB(ctx context.Context, deps Deps) error {
	err := db.InitializeDB(ctx, deps.GetDB())
	if err != nil {
		return eris.Wrap(err, "Failed to initialize db")
	}
	return nil
}

func InitializeTestDB(ctx context.Context, deps Deps) error {
	err := db.InitializeTestDB(ctx, deps.GetDB())
	if err != nil {
		return eris.Wrap(err, "Failed to initialize test db")
	}
	return nil
}
