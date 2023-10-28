package wiring

import "database/sql"

type Deps interface {
	GetDB() *sql.DB
}

type WiringDeps struct {
	DB *sql.DB
}

func (d *WiringDeps) GetDB() *sql.DB {
	return d.DB
}
