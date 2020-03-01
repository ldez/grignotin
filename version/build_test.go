package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBuild(t *testing.T) {
	build, err := GetBuild()
	require.NoError(t, err)

	require.NotNil(t, build)
	assert.NotEmpty(t, build.Builders)
	assert.NotEmpty(t, build.Revisions)
}
