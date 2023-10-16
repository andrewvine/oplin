package resources

import (
	"embed"
)

//go:embed  all:static
var Static embed.FS

//go:embed  all:templates
var Templates embed.FS
