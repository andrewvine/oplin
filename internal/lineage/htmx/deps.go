package htmx

import "database/sql"

type Deps interface {
	GetDB() *sql.DB
}
