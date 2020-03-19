package backups

import (
	"encoding/json"
	"fmt"
	"github.com/previousnext/acquia-cli/internal/acquia/api"
	"github.com/previousnext/acquia-cli/internal/acquia/auth"
	"io/ioutil"
	"net/http"
)

// ListResponse is returned from the Acquia API.
type ListResponse struct {
	Embedded ListResponseEmbedded `json:"_embedded"`
}

// ListResponseEmbedded is returned from the Acquia API.
type ListResponseEmbedded struct {
	Items *BackupList `json:"items"`
}

// List all backups.
func List(token *auth.Client, environment, database string) (*BackupList, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/environments/%s/databases/%s/backups?sort=created", api.BasePath, environment, database), nil)
	if err != nil {
		return nil, err
	}

	token.WrapRequest(req)

	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var list ListResponse

	err = json.Unmarshal(body, &list)
	if err != nil {
		return nil, err
	}

	return list.Embedded.Items, nil
}