package goenv

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetOne(t *testing.T) {
	p, err := GetOne(GOMOD)
	require.NoError(t, err)

	abs, err := filepath.Abs("..")
	require.NoError(t, err)

	assert.Equal(t, filepath.Join(abs, "go.mod"), p)
}

func TestGet(t *testing.T) {
	values, err := Get(GOMOD, GO111MODULE)
	require.NoError(t, err)

	assert.NotEmpty(t, values[GOMOD])
	assert.NotEmpty(t, values[GO111MODULE])
	assert.Empty(t, values[GOEXE])
}
