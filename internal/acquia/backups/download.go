package backups

import (
	"compress/gzip"
	"fmt"
	"github.com/previousnext/acquia-cli/internal/acquia/api"
	"github.com/previousnext/acquia-cli/internal/acquia/auth"
	"io"
	"net/http"
)

// Download the latest database backup.
func Download(token *auth.Client, w io.Writer, backup Backup) error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/environments/%s/databases/%s/backups/%d/actions/download", api.BasePath, backup.Environment.ID, backup.Database.Name, backup.ID), nil)
	if err != nil {
		return err
	}

	token.WrapRequest(req)

	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return fmt.Errorf("failed to setup http client: %w", err)
	}
	defer resp.Body.Close()

	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to uncompress backup: %w", err)
	}
	defer reader.Close()

	_, err = io.Copy(w, reader)
	if err != nil {
		return fmt.Errorf("failed to save backup: %w", err)
	}

	return nil
}