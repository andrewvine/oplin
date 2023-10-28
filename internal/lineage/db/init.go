package db

import (
	"context"
	"database/sql"
	"oplin/internal/utils"
	"strings"

	"github.com/rotisserie/eris"
)

func InitializeDB(ctx context.Context, db *sql.DB) error {
	sql := "SELECT exists(select schema_name FROM information_schema.schemata WHERE schema_name = 'lineage');"
	rows, err := db.Query(sql)
	if err != nil {
		return eris.Wrap(err, "Failed when querying if schema exists")
	}
	defer rows.Close()

	var exists bool
	for rows.Next() {
		err := rows.Scan(&exists)
		if err != nil {
			return eris.Wrap(err, "Failed when scanning if schema exists")
		}
	}
	if !exists {
		stmts := strings.Split(SchemaSQL, ";")
		err := utils.RunStatements(db, stmts)
		if err != nil {
			return eris.Wrap(err, "Failed to initialize db")
		}
	}
	return nil
}

func InitializeTestDB(ctx context.Context, db *sql.DB) error {
	stmts := strings.Split(DropSQL, ";")
	err := utils.RunStatements(db, stmts)
	if err != nil {
		return eris.Wrap(err, "Failed to initialize db")
	}
	stmts = strings.Split(SchemaSQL, ";")
	err = utils.RunStatements(db, stmts)
	if err != nil {
		return eris.Wrap(err, "Failed to initialize db")
	}
	return nil
}
