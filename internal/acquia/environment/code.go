package environment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/previousnext/acquia-cli/internal/acquia/api"
	"github.com/previousnext/acquia-cli/internal/acquia/auth"
)

// CodeSwitchRequest used when executing a request.
type CodeSwitchRequest struct {
	Branch string `json:"branch"`
}

// CodeSwitchResponse returned from a request.
type CodeSwitchResponse struct {
	Links CodeSwitchResponseLinks `json:"_links"`
}

// CodeSwitchResponseLinks which contains a link to the notification.
type CodeSwitchResponseLinks struct {
	Notification CodeResponseLinkNotification `json:"notification"`
}

// CodeResponseLinkNotification which contains a link.
type CodeResponseLinkNotification struct {
	URL string `json:"href"`
}

// CodeSwitch (AKA deploy) an environment.
func CodeSwitch(token *auth.Client, environment, branch string) (string, error) {
	request := &CodeSwitchRequest{
		Branch: branch,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/environments/%s/code/actions/switch", api.BasePath, environment), bytes.NewReader(requestBody))
	if err != nil {
		return "", err
	}

	token.WrapRequest(req)

	req.Header.Set("Content-Type", "application/json")

	cli := &http.Client{}

	resp, err := cli.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to setup http client: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	msg := string(body)

	var response CodeSwitchResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	// Ensure that we have a valid url which has been returned.
	// If not then show the body of the response, this will contain more information.
	_, err = url.ParseRequestURI(response.Links.Notification.URL)
	if err != nil {
		return "", fmt.Errorf("failed to parse notification url: %s", msg)
	}

	// Read the last element from the notification URL to determine the notification UUID.
	sl := strings.Split(response.Links.Notification.URL, "/")
	notification := sl[len(sl)-1]

	return notification, nil
}