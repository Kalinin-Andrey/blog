package migrations

import (
	"embed"
	"github.com/Kalinin-Andrey/blog/internal/infrastructure/repository/tsdb"
)

//go:embed *.sql
var EmbedMigrations embed.FS

var CurrentRepo *tsdb.Repository
