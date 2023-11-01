package main

import (
	_ "embed"
	"log"
	"oplin/internal/utils"
	"strings"

	"github.com/rotisserie/eris"
)

//go:embed init.sql
var InitSQL string


func main() {
	db := utils.GetTestDB()
	stmts := strings.Split(InitSQL, ";")
	err := utils.RunStatements(db, stmts)
	if err != nil {
		log.Fatal(eris.Wrap(err, "Failed to initialize db"))
	}
}