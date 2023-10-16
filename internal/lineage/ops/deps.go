package ops

import "database/sql"

type Deps interface {
	GetDB() *sql.DB
}

type TestDeps struct {
	DB *sql.DB
}

func (d *TestDeps) GetDB() *sql.DB {
	return d.DB
}
