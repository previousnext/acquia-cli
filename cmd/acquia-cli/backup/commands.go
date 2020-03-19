package backup

import (
	"github.com/alecthomas/kingpin"

	"github.com/previousnext/acquia-cli/cmd/acquia-cli/backup/dump"
)

// Commands used as part of a backups workflow.
func Commands(app *kingpin.Application) {
	cmd := app.Command("backup", "Backup commands")
	dump.Command(cmd)
}