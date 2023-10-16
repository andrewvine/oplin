package api

import "database/sql"

type Deps interface {
	GetDB() *sql.DB
}
