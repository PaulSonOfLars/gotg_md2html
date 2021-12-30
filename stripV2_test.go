package tg_md2html_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	tg_md2html "github.com/PaulSonOfLars/gotg_md2html"
)

func TestStripMDV2(t *testing.T) {
	for _, x := range stripMD {
		assert.Equal(t, x.output, tg_md2html.StripMDV2(x.in), "failed to strip all markdown")
	}
}
