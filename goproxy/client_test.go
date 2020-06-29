package goproxy

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_Lookup(t *testing.T) {
	client := NewClient("", "")

	raw, err := client.Lookup("github.com/ldez/grignotin", "v0.1.0")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(raw)
}

func TestClient_GetSources(t *testing.T) {
	client := NewClient("", "")

	raw, err := client.GetSources("github.com/ldez/grignotin", "v0.1.0")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(raw))
}

func TestClient_DownloadSources(t *testing.T) {
	client := NewClient("", "")

	reader, err := client.DownloadSources("github.com/ldez/grignotin", "v0.1.0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = reader.Close() }()

	raw, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(raw))
}

func TestClient_GetVersions(t *testing.T) {
	client := NewClient("", "")

	versions, err := client.GetVersions("github.com/hashicorp/consul/api")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(versions)
}

func TestClient_GetInfo(t *testing.T) {
	client := NewClient("", "")

	info, err := client.GetInfo("github.com/ijc25/Gotty", "a8b993ba6abdb0e0c12b0125c603323a71c7790c")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(info)
}

func TestClient_GetLatest(t *testing.T) {
	client := NewClient("", "")

	info, err := client.GetLatest("golang.org/x/lint")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(info)
}

func Test_safeModuleName(t *testing.T) {
	testCases := []struct {
		name     string
		expected string
	}{
		{
			name:     "AkamaiOPEN-edgegrid-golang",
			expected: "!akamai!o!p!e!n-edgegrid-golang",
		},
		{
			name:     "golang.org/x/lint",
			expected: "golang.org/x/lint",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actual := safeModuleName(test.name)

			assert.Equal(t, test.expected, actual)
		})
	}
}
