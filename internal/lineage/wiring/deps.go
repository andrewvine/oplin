package wiring

import "database/sql"

type Deps struct {
	DB *sql.DB
}

func (d *Deps) GetDB() *sql.DB {
	return d.DB
}
