package encoder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceFileExtension(t *testing.T) {

	// out extension without a dot prefix
	assert.Equal(t, "/var/foo.mp4", ReplaceFileExtension("/var/foo.flv", "mp4"))

	// our extension with a dot prefix
	assert.Equal(t, "/var/foo.mp4", ReplaceFileExtension("/var/foo.flv", ".mp4"))
}

func TestDurationToInterval(t *testing.T) {
	assert.Equal(t, uint(127), DurationToInterval(float64(7993.92)))
	assert.Equal(t, uint(34), DurationToInterval(float64(500)))
}
