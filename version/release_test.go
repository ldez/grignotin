package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetReleases(t *testing.T) {
	releases, err := GetReleases(false)
	require.NoError(t, err)

	assert.Len(t, releases, 2)
}
