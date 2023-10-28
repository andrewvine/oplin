package db

import (
	_ "embed"
)

//go:embed schema.sql
var SchemaSQL string

//go:embed drop.sql
var DropSQL string
