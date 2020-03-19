package deploy

import (
	"fmt"
	"github.com/alecthomas/kingpin"
	"github.com/previousnext/acquia-cli/internal/acquia/environment"
	"os"

	"github.com/previousnext/acquia-cli/internal/acquia/auth"
)

type command struct {
	Key string
	Secret string
	Environment string
	Branch string
}

// Helper function to run a database dump from an Acquia backup.
func (cmd *command) run(c *kingpin.ParseContext) error {
	fmt.Fprintln(os.Stderr, "Intializing Acquia client")

	a, err := auth.New(cmd.Key, cmd.Secret)
	if err != nil {
		return fmt.Errorf("failed creating authentication client: %w", err)
	}

	fmt.Fprintln(os.Stderr, "Triggering a code switch")

	_, err = environment.CodeSwitch(a, cmd.Environment, cmd.Branch)
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "Code switch submitted")

	return nil
}

// Command which dumps a database from backups.
func Command(app *kingpin.Application) {
	cmd := new(command)

	command := app.Command("deploy", "Triggers a deployment").Action(cmd.run)
	command.Flag("key", "Key used to authenticate to the Acquia API").Envar("ACQUIA_KEY").StringVar(&cmd.Key)
	command.Flag("secret", "Secret used to authenticate to the Acquia API").Envar("ACQUIA_SECRET").StringVar(&cmd.Secret)
	command.Arg("environment", "UUID of the environment which we will be triggering a deployment").Required().StringVar(&cmd.Environment)
	command.Arg("branch", "Branch to deploy").Required().StringVar(&cmd.Branch)
}

