// Package goproxy simple client for go modules proxy
// https://docs.gomods.io/intro/protocol/
// https://go.googlesource.com/proposal/+/master/design/25530-sumdb.md
package goproxy

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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
	proxyURL   string
	HTTPClient *http.Client
}

// NewClient creates a new Client.
func NewClient(proxyURL string) *Client {
	client := &Client{
		HTTPClient: &http.Client{},
	}

	client.proxyURL = defaultProxyURL
	if proxyURL != "" {
		client.proxyURL = proxyURL
	}

	return client
}

// GetSources gets the contents of the archive file.
func (c *Client) GetSources(moduleName string, version string) ([]byte, error) {
	uri := fmt.Sprintf("%s/%s/@v/%s.zip", c.proxyURL, mustEscapePath(moduleName), version)

	resp, err := c.HTTPClient.Get(uri)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode/100 != 2 {
		return nil, handleError(resp)
	}

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return raw, nil
}

// DownloadSources returns an io.ReadCloser that reads the contents of the archive file.
// It is the caller's responsibility to close the ReadCloser.
func (c *Client) DownloadSources(moduleName string, version string) (io.ReadCloser, error) {
	uri := fmt.Sprintf("%s/%s/@v/%s.zip", c.proxyURL, mustEscapePath(moduleName), version)

	resp, err := c.HTTPClient.Get(uri)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		return nil, handleError(resp)
	}

	return resp.Body, nil
}

// GetModFile gets go.mod file.
func (c *Client) GetModFile(moduleName string, version string) (*modfile.File, error) {
	uri := fmt.Sprintf("%s/%s/@v/%s.mod", c.proxyURL, mustEscapePath(moduleName), version)

	resp, err := c.HTTPClient.Get(uri)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode/100 != 2 {
		return nil, handleError(resp)
	}

	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return modfile.Parse("go.mod", all, nil)
}

// GetVersions gets all available module versions.
//	<proxy URL>/<module name>/@v/list
func (c *Client) GetVersions(moduleName string) ([]string, error) {
	uri := fmt.Sprintf("%s/%s/@v/list", c.proxyURL, mustEscapePath(moduleName))

	resp, err := c.HTTPClient.Get(uri)
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
//	<proxy URL>/<module name>/@v/<version>.info
func (c *Client) GetInfo(moduleName string, version string) (*VersionInfo, error) {
	return c.getInfo(fmt.Sprintf("%s/%s/@v/%s.info", c.proxyURL, mustEscapePath(moduleName), version))
}

// GetLatest gets information about the latest module version.
//	<proxy URL>/<module name>/@latest
func (c *Client) GetLatest(moduleName string) (*VersionInfo, error) {
	return c.getInfo(fmt.Sprintf("%s/%s/@latest", c.proxyURL, mustEscapePath(moduleName)))
}

func (c *Client) getInfo(uri string) (*VersionInfo, error) {
	resp, err := c.HTTPClient.Get(uri)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode/100 != 2 {
		return nil, handleError(resp)
	}

	info := VersionInfo{}
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
	all, _ := ioutil.ReadAll(resp.Body)

	return &APIError{
		StatusCode: resp.StatusCode,
		Message:    fmt.Sprintf("%s: %s", resp.Status, string(all)),
	}
}
