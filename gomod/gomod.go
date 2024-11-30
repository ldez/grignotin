// Package gomod A function to get information about module (go list).
package gomod

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
)

// ModInfo Module information.
//
//nolint:tagliatelle // temporary: the next version of golangci-lint will allow configuration by package.
type ModInfo struct {
	Path      string `json:"Path"`
	Dir       string `json:"Dir"`
	GoMod     string `json:"GoMod"`
	GoVersion string `json:"GoVersion"`
	Main      bool   `json:"Main"`
}

// GetModuleInfo gets modules information from `go list`.
func GetModuleInfo() ([]ModInfo, error) {
	// https://github.com/golang/go/issues/44753#issuecomment-790089020
	cmd := exec.Command("go", "list", "-m", "-json")

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("command go list: %w: %s", err, string(out))
	}

	var infos []ModInfo

	for dec := json.NewDecoder(bytes.NewBuffer(out)); dec.More(); {
		var v ModInfo
		if err := dec.Decode(&v); err != nil {
			return nil, fmt.Errorf("unmarshaling error: %w: %s", err, string(out))
		}

		if v.GoMod == "" {
			return nil, errors.New("working directory is not part of a module")
		}

		if !v.Main || v.Dir == "" {
			continue
		}

		infos = append(infos, v)
	}

	if len(infos) == 0 {
		return nil, errors.New("go.mod file not found")
	}

	return infos, nil
}
