package goproxy

import (
	"fmt"
	"io"
	"testing"
)

func TestClient_GetSources(t *testing.T) {
	client := NewClient("")

	raw, err := client.GetSources("github.com/ldez/grignotin", "v0.1.0")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(raw))
}

func TestClient_DownloadSources(t *testing.T) {
	client := NewClient("")

	reader, err := client.DownloadSources("github.com/ldez/grignotin", "v0.1.0")
	if err != nil {
		t.Fatal(err)
	}

	defer func() { _ = reader.Close() }()

	raw, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(raw))
}

func TestClient_GetVersions(t *testing.T) {
	client := NewClient("")

	versions, err := client.GetVersions("github.com/hashicorp/consul/api")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(versions)
}

func TestClient_GetInfo(t *testing.T) {
	client := NewClient("")

	info, err := client.GetInfo("github.com/ijc25/Gotty", "a8b993ba6abdb0e0c12b0125c603323a71c7790c")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(info)
}

func TestClient_GetModFile(t *testing.T) {
	client := NewClient("")

	file, err := client.GetModFile("github.com/ldez/grignotin", "v0.1.0")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(file)
}

func TestClient_GetLatest(t *testing.T) {
	client := NewClient("")

	info, err := client.GetLatest("golang.org/x/lint")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(info)
}
