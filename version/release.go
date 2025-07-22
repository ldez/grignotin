package version

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const baseDLURL = "https://golang.org/dl/"

// Release represents a release on the golang.org downloads page.
type Release struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
	Files   []File `json:"files"`
}

// File represents a file on the golang.org downloads page.
type File struct {
	Filename string `json:"filename"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Version  string `json:"version"`
	SHA256   string `json:"sha256"`
	Size     int    `json:"size"`
	Kind     string `json:"kind"`
}

// GetReleases gets build information.
func GetReleases(all bool) ([]Release, error) {
	return GetReleasesWithContext(context.Background(), all)
}

// GetReleasesWithContext gets build information.
func GetReleasesWithContext(ctx context.Context, all bool) ([]Release, error) {
	dlURL, err := url.Parse(baseDLURL)
	if err != nil {
		return nil, err
	}

	query := dlURL.Query()
	query.Set("mode", "json")

	if all {
		query.Set("include", "all")
	}

	dlURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, dlURL.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("invalid response, status code: %d", resp.StatusCode)
	}

	var releases []Release

	err = json.NewDecoder(resp.Body).Decode(&releases)
	if err != nil {
		return nil, err
	}

	return releases, nil
}
