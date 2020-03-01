package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const baseBuildURL = "https://build.golang.org/"

// Build information.
type Build struct {
	Builders  []string   `json:"builders"`
	Revisions []Revision `json:"revisions"`
}

// Revision information.
type Revision struct {
	Repo       string    `json:"repo"`
	Revision   string    `json:"revision"`
	Date       time.Time `json:"date"`
	Branch     string    `json:"branch"`
	Author     string    `json:"author"`
	Desc       string    `json:"desc"`
	Results    []string  `json:"results"`
	GoRevision string    `json:"goRevision,omitempty"`
	GoBranch   string    `json:"goBranch,omitempty"`
}

// GetBuild gets build information.
func GetBuild() (*Build, error) {
	dlURL, err := url.Parse(baseBuildURL)
	if err != nil {
		return nil, err
	}

	query := dlURL.Query()
	query.Set("mode", "json")
	dlURL.RawQuery = query.Encode()

	resp, err := http.Get(dlURL.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("invalid response, status code: %d", resp.StatusCode)
	}

	var build Build
	err = json.NewDecoder(resp.Body).Decode(&build)
	if err != nil {
		return nil, err
	}

	return &build, nil
}
