package tg_md2html

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStripMDV2(t *testing.T) {
	for _, x := range stripMD {
		assert.Equal(t, x.output, StripMDV2(x.in), "failed to strip all markdown")
	}
}
