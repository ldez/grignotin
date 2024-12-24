package gomod

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetModuleInfo(t *testing.T) {
	info, err := GetModuleInfo()
	require.NoError(t, err)

	require.Len(t, info, 1)

	assert.Equal(t, "github.com/ldez/grignotin", info[0].Path)
	assert.Equal(t, "1.22.0", info[0].GoVersion)
	assert.True(t, info[0].Main)
}

func TestGetGoModPath(t *testing.T) {
	p, err := GetGoModPath()
	require.NoError(t, err)

	abs, err := filepath.Abs("..")
	require.NoError(t, err)

	assert.Equal(t, filepath.Join(abs, "go.mod"), p)
}

func TestGetModulePath(t *testing.T) {
	p, err := GetModulePath()
	require.NoError(t, err)

	assert.Equal(t, "github.com/ldez/grignotin", p)
}
