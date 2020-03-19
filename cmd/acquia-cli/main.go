package main

import (
	"os"

	"github.com/alecthomas/kingpin"

	"github.com/previousnext/acquia-cli/cmd/acquia-cli/deploy"
	"github.com/previousnext/acquia-cli/cmd/acquia-cli/backup"
)

func main() {
	app := kingpin.New("acquia-cli", "Utility for interfacing with MySQL on Acquia")

	backup.Commands(app)
	deploy.Command(app)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}