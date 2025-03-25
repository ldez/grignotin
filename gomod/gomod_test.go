package gomod

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetModuleInfo(t *testing.T) {
	info, err := GetModuleInfo(context.Background())
	require.NoError(t, err)

	require.Len(t, info, 1)

	assert.Equal(t, "github.com/ldez/grignotin", info[0].Path)
	assert.Equal(t, "1.23.0", info[0].GoVersion)
	assert.True(t, info[0].Main)
}

func TestGetModulePath(t *testing.T) {
	p, err := GetModulePath(context.Background())
	require.NoError(t, err)

	assert.Equal(t, "github.com/ldez/grignotin", p)
}
