package tg_md2html

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverseV2(t *testing.T) {
	for _, test := range reverseTest {
		out, err := ReverseV2(MD2HTMLV2(test), nil)
		assert.NoError(t, err, "Error for:\n%s", test)
		assert.Equal(t, MD2HTMLV2(test), MD2HTMLV2(out))
	}

	for _, test := range []string{
		"___________test_______",
	} {
		out, err := ReverseV2(MD2HTMLV2(test), nil)
		assert.NoError(t, err, "Error for:\n%s", test)
		assert.Equal(t, MD2HTMLV2(test), MD2HTMLV2(out))
	}
}
