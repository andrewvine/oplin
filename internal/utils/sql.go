package utils

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"
)

func RunSQLFile(db *sql.DB, filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	stmts := strings.Split(string(content), ";")
	return RunStatements(db, stmts)
}

func RunStatements(db *sql.DB, stmts []string) error {
	for _, stmt := range stmts {
		stmt = strings.TrimSpace(stmt)
		if len(stmt) == 0 {
			continue
		}
		_, err := db.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func DropSchema(db *sql.DB, schemaName string) error {
	return RunStatements(db, []string{
		fmt.Sprintf("drop schema if exists %s cascade;", schemaName),
	})
}
