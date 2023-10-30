package utils

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"
)

// RunSQLFile executes the SQL statements in the file named filename
func RunSQLFile(db *sql.DB, filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	stmts := strings.Split(string(content), ";")
	return RunStatements(db, stmts)
}

// RunStatements executes the SQL statements in the slice of strings
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

// DropSchema drops the schema named schemaName
func DropSchema(db *sql.DB, schemaName string) error {
	return RunStatements(db, []string{
		fmt.Sprintf("drop schema if exists %s cascade;", schemaName),
	})
}
