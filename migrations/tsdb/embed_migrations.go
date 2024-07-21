package tsdb

import (
	"embed"

	"blog/internal/infrastructure/repository/tsdb"
)

//go:embed *.sql
var EmbedMigrations embed.FS

var CurrentRepo *tsdb.Repository
