// Package goenv A set of functions to get information from `go env`.
package goenv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// GetAll gets information from "go env".
func GetAll() (map[string]string, error) {
	cmd := exec.Command("go", "env", "-json")

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("command %q: %w: %s", strings.Join(cmd.Args, " "), err, string(out))
	}

	v := map[string]string{}
	err = json.NewDecoder(bytes.NewBuffer(out)).Decode(&v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// GetOne gets information from "go env" for one environment variable.
func GetOne(name string) (string, error) {
	cmd := exec.Command("go", "env", "-json", name)

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("command %q: %w: %s", strings.Join(cmd.Args, " "), err, string(out))
	}

	v := map[string]string{}
	err = json.NewDecoder(bytes.NewBuffer(out)).Decode(&v)
	if err != nil {
		return "", err
	}

	return v[name], nil
}

// Get gets information from "go env" for one or several environment variables.
func Get(name ...string) (map[string]string, error) {
	args := append([]string{"env", "-json"}, name...)
	cmd := exec.Command("go", args...) //nolint:gosec // The env var names must be checked by the user.

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("command %q: %w: %s", strings.Join(cmd.Args, " "), err, string(out))
	}

	v := map[string]string{}
	err = json.NewDecoder(bytes.NewBuffer(out)).Decode(&v)
	if err != nil {
		return nil, err
	}

	return v, nil
}
