package tg_md2html

import (
	"strings"
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

func TestReverseV2Buttons(t *testing.T) {
	for _, x := range []struct {
		in   string
		out  string
		btns []ButtonV2
	}{
		{
			in:  "[hello](buttonurl://test.com)",
			out: "",
			btns: []ButtonV2{{
				Name:    "hello",
				Content: "test.com",
			}},
		}, {
			in:  "Some text, some *bold*, and a button\n[hello](buttonurl://test.com)",
			out: "Some text, some <b>bold</b>, and a button",
			btns: []ButtonV2{{
				Name:    "hello",
				Content: "test.com",
			}},
		},
	} {

		txt, b := MD2HTMLButtonsV2(x.in)
		txt = strings.TrimSpace(txt)
		assert.Equal(t, x.out, txt)
		assert.ElementsMatch(t, x.btns, b)
		out, err := ReverseV2(txt, x.btns)
		assert.NoError(t, err, "no error expected")
		assert.Equal(t, x.in, strings.TrimSpace(out))
	}
}