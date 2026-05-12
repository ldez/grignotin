package version

import (
	"context"
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
	return GetBuildWithContext(context.Background())
}

// GetBuildWithContext gets build information.
func GetBuildWithContext(ctx context.Context) (*Build, error) {
	endpoint, err := url.Parse(baseBuildURL)
	if err != nil {
		return nil, err
	}

	query := endpoint.Query()
	query.Set("mode", "json")
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

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
