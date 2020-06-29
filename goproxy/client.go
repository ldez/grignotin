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
	"unicode"
)

const (
	defaultProxyURL            = "https://proxy.golang.org"
	defaultChecksumDatabaseURL = "https://sum.golang.org"
)

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
	sumDBURL   string
	HTTPClient *http.Client
}

// NewClient creates a new Client.
func NewClient(proxyURL string, sumDBURL string) *Client {
	client := &Client{
		HTTPClient: &http.Client{},
	}

	client.proxyURL = defaultProxyURL
	if proxyURL != "" {
		client.proxyURL = proxyURL
	}

	client.sumDBURL = defaultChecksumDatabaseURL
	if sumDBURL != "" {
		client.sumDBURL = sumDBURL
	}

	return client
}

// Lookup gets checksum information.
//	<sumDB URL>/lookup/<module name>@<version>
func (c *Client) Lookup(moduleName, version string) (string, error) {
	uri := fmt.Sprintf("%s/lookup/%s@%s", c.sumDBURL, safeModuleName(moduleName), version)
	resp, err := c.HTTPClient.Get(uri)
	if err != nil {
		return "", err
	}

	defer func() { _ = resp.Body.Close() }()

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode/100 != 2 {
		return "", fmt.Errorf("invalid response: %s [%d]: %s", resp.Status, resp.StatusCode, string(raw))
	}

	return string(raw), nil
}

// GetSources gets the contents of the archive file.
func (c *Client) GetSources(moduleName string, version string) ([]byte, error) {
	uri := fmt.Sprintf("%s/%s/@v/%s.zip", c.proxyURL, safeModuleName(moduleName), version)

	resp, err := c.HTTPClient.Get(uri)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode/100 != 2 {
		raw, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("invalid response: %s [%d]: %s", resp.Status, resp.StatusCode, string(raw))
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
	uri := fmt.Sprintf("%s/%s/@v/%s.zip", c.proxyURL, safeModuleName(moduleName), version)

	resp, err := c.HTTPClient.Get(uri)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		raw, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("invalid response: %s [%d]: %s", resp.Status, resp.StatusCode, string(raw))
	}

	return resp.Body, nil
}

// GetVersions gets all available module versions.
//	<proxy URL>/<module name>/@v/list
func (c *Client) GetVersions(moduleName string) ([]string, error) {
	uri := fmt.Sprintf("%s/%s/@v/list", c.proxyURL, safeModuleName(moduleName))

	resp, err := c.HTTPClient.Get(uri)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("invalid response: %s [%d]", resp.Status, resp.StatusCode)
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
	return c.getInfo(fmt.Sprintf("%s/%s/@v/%s.info", c.proxyURL, safeModuleName(moduleName), version))
}

// GetLatest gets information about the latest module version.
//	<proxy URL>/<module name>/@latest
func (c *Client) GetLatest(moduleName string) (*VersionInfo, error) {
	return c.getInfo(fmt.Sprintf("%s/%s/@latest", c.proxyURL, safeModuleName(moduleName)))
}

func (c *Client) getInfo(uri string) (*VersionInfo, error) {
	resp, err := c.HTTPClient.Get(uri)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("invalid response: %s [%d]", resp.Status, resp.StatusCode)
	}

	info := VersionInfo{}
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func safeModuleName(name string) string {
	var to []byte
	for _, r := range name {
		if 'A' <= r && r <= 'Z' {
			to = append(to, '!', byte(unicode.ToLower(r)))
		} else {
			to = append(to, byte(r))
		}
	}

	return string(to)
}
