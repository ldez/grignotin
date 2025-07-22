// Package goproxy simple client for go modules proxy
// https://golang.org/cmd/go/#hdr-Module_proxy_protocol
// https://docs.gomods.io/intro/protocol/
// https://go.googlesource.com/proposal/+/master/design/25530-sumdb.md
package goproxy

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

const defaultProxyURL = "https://proxy.golang.org"

// VersionInfo is the representation of a version.
type VersionInfo struct {
	Name    string
	Short   string
	Version string
	Time    time.Time
}

// Client is the go modules proxy client.
type Client struct {
	proxyURL   *url.URL
	HTTPClient *http.Client
}

// NewClient creates a new Client.
func NewClient(proxyURL string) *Client {
	client := &Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}

	client.proxyURL, _ = url.Parse(defaultProxyURL)
	if proxyURL != "" {
		var err error

		client.proxyURL, err = url.Parse(proxyURL)
		if err != nil {
			// Use a panic to be non-breaking, but the [NewClient] signature must be changed.
			panic(err)
		}
	}

	return client
}

// GetSources gets the contents of the archive file.
func (c *Client) GetSources(moduleName, version string) ([]byte, error) {
	return c.GetSourcesWithContext(context.Background(), moduleName, version)
}

// GetSourcesWithContext gets the contents of the archive file.
func (c *Client) GetSourcesWithContext(ctx context.Context, moduleName, version string) ([]byte, error) {
	endpoint := c.proxyURL.JoinPath(mustEscapePath(moduleName), "@v", version+".zip")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode/100 != 2 {
		return nil, handleError(resp)
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return raw, nil
}

// DownloadSources returns an io.ReadCloser that reads the contents of the archive file.
// It is the caller's responsibility to close the ReadCloser.
func (c *Client) DownloadSources(moduleName, version string) (io.ReadCloser, error) {
	return c.DownloadSourcesWithContext(context.Background(), moduleName, version)
}

// DownloadSourcesWithContext returns an io.ReadCloser that reads the contents of the archive file.
// It is the caller's responsibility to close the ReadCloser.
func (c *Client) DownloadSourcesWithContext(ctx context.Context, moduleName, version string) (io.ReadCloser, error) {
	endpoint := c.proxyURL.JoinPath(mustEscapePath(moduleName), "@v", version+".zip")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		return nil, handleError(resp)
	}

	return resp.Body, nil
}

// GetModFile gets go.mod file.
func (c *Client) GetModFile(moduleName, version string) (*modfile.File, error) {
	return c.GetModFileWithContext(context.Background(), moduleName, version)
}

// GetModFileWithContext gets go.mod file.
func (c *Client) GetModFileWithContext(ctx context.Context, moduleName, version string) (*modfile.File, error) {
	endpoint := c.proxyURL.JoinPath(mustEscapePath(moduleName), "@v", version+".mod")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode/100 != 2 {
		return nil, handleError(resp)
	}

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return modfile.Parse("go.mod", all, nil)
}

// GetVersions gets all available module versions.
//
//	<proxy URL>/<module name>/@v/list
func (c *Client) GetVersions(moduleName string) ([]string, error) {
	return c.GetVersionsWithContext(context.Background(), moduleName)
}

// GetVersionsWithContext gets all available module versions.
//
//	<proxy URL>/<module name>/@v/list
func (c *Client) GetVersionsWithContext(ctx context.Context, moduleName string) ([]string, error) {
	endpoint := c.proxyURL.JoinPath(mustEscapePath(moduleName), "@v", "list")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode/100 != 2 {
		return nil, handleError(resp)
	}

	var versions []string

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		versions = append(versions, line)
	}

	return versions, nil
}

// GetInfo gets information about a module version.
//
//	<proxy URL>/<module name>/@v/<version>.info
func (c *Client) GetInfo(moduleName, version string) (*VersionInfo, error) {
	return c.GetInfoWithContext(context.Background(), moduleName, version)
}

// GetInfoWithContext gets information about a module version.
//
//	<proxy URL>/<module name>/@v/<version>.info
func (c *Client) GetInfoWithContext(ctx context.Context, moduleName, version string) (*VersionInfo, error) {
	return c.getInfo(ctx, c.proxyURL.JoinPath(mustEscapePath(moduleName), "@v", version+".info"))
}

// GetLatest gets information about the latest module version.
//
//	<proxy URL>/<module name>/@latest
func (c *Client) GetLatest(moduleName string) (*VersionInfo, error) {
	return c.GetLatestWithContext(context.Background(), moduleName)
}

// GetLatestWithContext gets information about the latest module version.
//
//	<proxy URL>/<module name>/@latest
func (c *Client) GetLatestWithContext(ctx context.Context, moduleName string) (*VersionInfo, error) {
	return c.getInfo(ctx, c.proxyURL.JoinPath(mustEscapePath(moduleName), "@latest"))
}

func (c *Client) getInfo(ctx context.Context, endpoint *url.URL) (*VersionInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode/100 != 2 {
		return nil, handleError(resp)
	}

	info := VersionInfo{}
	//nolint:musttag // data from Go proxy.
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func mustEscapePath(path string) string {
	escapePath, err := module.EscapePath(path)
	if err != nil {
		panic(err)
	}

	return escapePath
}

func handleError(resp *http.Response) error {
	all, _ := io.ReadAll(resp.Body)

	return &APIError{
		StatusCode: resp.StatusCode,
		Message:    fmt.Sprintf("%s: %s", resp.Status, string(all)),
	}
}
