package dump

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"

	"github.com/previousnext/acquia-cli/internal/acquia/auth"
	"github.com/previousnext/acquia-cli/internal/acquia/backups"
)

type command struct {
	Key string
	Secret string
	Environment string
	Database string
}

// Helper function to run a database dump from an Acquia backup.
func (cmd *command) run(c *kingpin.ParseContext) error {
	fmt.Fprintln(os.Stderr, "Intializing Acquia client")

	client, err := auth.New(cmd.Key, cmd.Secret)
	if err != nil {
		return fmt.Errorf("failed creating authentication client: %w", err)
	}

	fmt.Fprintln(os.Stderr, "Listing backups")

	list, err := backups.List(client, cmd.Environment, cmd.Database)
	if err != nil {
		return fmt.Errorf("failed to download backup: %w", err)
	}

	backup := backups.GetLatest(list)

	fmt.Fprintf(os.Stderr, "Downloading and extracting the latest backup (%d)\n", backup.ID)

	err = backups.Download(client, os.Stdout, backup)
	if err != nil {
		return fmt.Errorf("failed to download backup: %w", err)
	}

	fmt.Fprintln(os.Stderr, "Complete!")

	return nil
}

// Command which dumps a database from backups.
func Command(app *kingpin.CmdClause) {
	cmd := new(command)

	command := app.Command("dump", "Dumps the latest database backup").Action(cmd.run)
	command.Flag("key", "Key used to authenticate to the Acquia API").Envar("ACQUIA_KEY").StringVar(&cmd.Key)
	command.Flag("secret", "Secret used to authenticate to the Acquia API").Envar("ACQUIA_SECRET").StringVar(&cmd.Secret)
	command.Flag("environment", "ID of the environment").Envar("ACQUIA_ENVIRONMENT").StringVar(&cmd.Environment)
	command.Flag("database", "Name of the database").Envar("ACQUIA_DATABASE").StringVar(&cmd.Database)
}

