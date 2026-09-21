// Package migrations embeds the SQL migration files so goose can run them
// from the compiled binary without needing the source tree on disk.
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
