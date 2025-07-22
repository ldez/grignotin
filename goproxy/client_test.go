package goproxy

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetSources(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("GET /github.com/ldez/grignotin/@v/v0.1.0.zip",
		func(rw http.ResponseWriter, _ *http.Request) {
			_, _ = fmt.Fprint(rw, "hello")
		})

	client := NewClient(server.URL)

	raw, err := client.GetSources("github.com/ldez/grignotin", "v0.1.0")
	require.NoError(t, err)

	assert.Equal(t, "hello", string(raw))
}

func TestClient_GetSources_integration(t *testing.T) {
	client := NewClient("")

	raw, err := client.GetSources("github.com/ldez/grignotin", "v0.1.0")
	require.NoError(t, err)

	fmt.Println(string(raw))
}

func TestClient_DownloadSources(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("GET /github.com/ldez/grignotin/@v/v0.1.0.zip",
		func(rw http.ResponseWriter, _ *http.Request) {
			_, _ = fmt.Fprint(rw, "hello")
		})

	client := NewClient(server.URL)

	reader, err := client.DownloadSources("github.com/ldez/grignotin", "v0.1.0")
	require.NoError(t, err)

	defer func() { _ = reader.Close() }()

	raw, err := io.ReadAll(reader)
	require.NoError(t, err)

	assert.Equal(t, "hello", string(raw))
}

func TestClient_DownloadSources_integration(t *testing.T) {
	client := NewClient("")

	reader, err := client.DownloadSources("github.com/ldez/grignotin", "v0.1.0")
	require.NoError(t, err)

	defer func() { _ = reader.Close() }()

	raw, err := io.ReadAll(reader)
	require.NoError(t, err)

	fmt.Println(string(raw))
}

func TestClient_GetVersions(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("GET /github.com/hashicorp/consul/api/@v/list",
		func(rw http.ResponseWriter, _ *http.Request) {
			_, _ = fmt.Fprint(rw, "v1.2.3\nv1.2.4\nv1.2.5\n")
		})

	client := NewClient(server.URL)

	versions, err := client.GetVersions("github.com/hashicorp/consul/api")
	require.NoError(t, err)

	expected := []string{"v1.2.3", "v1.2.4", "v1.2.5"}
	assert.Equal(t, expected, versions)
}

func TestClient_GetVersions_integration(t *testing.T) {
	client := NewClient("")

	versions, err := client.GetVersions("github.com/hashicorp/consul/api")
	require.NoError(t, err)

	fmt.Println(versions)
}

func TestClient_GetInfo(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("GET /github.com/ijc25/!gotty/@v/a8b993ba6abdb0e0c12b0125c603323a71c7790c.info",
		func(rw http.ResponseWriter, _ *http.Request) {
			_, _ = fmt.Fprint(rw, `{"Version":"v0.0.0-20170406111628-a8b993ba6abd","Time":"2017-04-06T11:16:28Z"}`)
		})

	client := NewClient(server.URL)

	info, err := client.GetInfo("github.com/ijc25/Gotty", "a8b993ba6abdb0e0c12b0125c603323a71c7790c")
	require.NoError(t, err)

	expected := &VersionInfo{
		Version: "v0.0.0-20170406111628-a8b993ba6abd",
		Time:    time.Date(2017, time.April, 6, 11, 16, 28, 0, time.UTC),
	}

	assert.Equal(t, expected, info)
}

func TestClient_GetInfo_integration(t *testing.T) {
	client := NewClient("")

	info, err := client.GetInfo("github.com/ijc25/Gotty", "a8b993ba6abdb0e0c12b0125c603323a71c7790c")
	require.NoError(t, err)

	fmt.Println(info)
}

func TestClient_GetModFile(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("GET /github.com/ldez/grignotin/@v/v0.1.0.mod",
		func(rw http.ResponseWriter, _ *http.Request) {
			_, _ = fmt.Fprint(rw, `module github.com/ldez/grignotin

go 1.13

require github.com/stretchr/testify v1.5.1
`)
		})

	client := NewClient(server.URL)

	file, err := client.GetModFile("github.com/ldez/grignotin", "v0.1.0")
	require.NoError(t, err)

	assert.NotNil(t, file)
}

func TestClient_GetModFile_integration(t *testing.T) {
	client := NewClient("")

	file, err := client.GetModFile("github.com/ldez/grignotin", "v0.1.0")
	require.NoError(t, err)

	fmt.Println(file)
}

func TestClient_GetLatest(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("GET /golang.org/x/lint/@latest",
		func(rw http.ResponseWriter, _ *http.Request) {
			_, _ = fmt.Fprint(rw, `{"Version":"v0.0.0-20241112194109-818c5a804067","Time":"2024-11-12T19:41:09Z","Origin":{"VCS":"git","URL":"https://go.googlesource.com/lint","Hash":"818c5a80406779e3ce2860365fc289de6d133b00"}}`)
		})

	client := NewClient(server.URL)

	info, err := client.GetLatest("golang.org/x/lint")
	require.NoError(t, err)

	expected := &VersionInfo{
		Version: "v0.0.0-20241112194109-818c5a804067",
		Time:    time.Date(2024, time.November, 12, 19, 41, 9, 0, time.UTC),
	}

	assert.Equal(t, expected, info)
}

func TestClient_GetLatest_integration(t *testing.T) {
	client := NewClient("")

	info, err := client.GetLatest("golang.org/x/lint")
	require.NoError(t, err)

	fmt.Println(info)
}
