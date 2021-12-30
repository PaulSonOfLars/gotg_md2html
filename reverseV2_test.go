package tg_md2html_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	tg_md2html "github.com/PaulSonOfLars/gotg_md2html"
)

func TestReverseV2(t *testing.T) {
	for _, test := range reverseTest {
		out, err := tg_md2html.ReverseV2(tg_md2html.MD2HTMLV2(test), nil)
		assert.NoError(t, err, "Error for:\n%s", test)
		assert.Equal(t, tg_md2html.MD2HTMLV2(test), tg_md2html.MD2HTMLV2(out))
	}

	for _, test := range []string{
		"___________test_______",
		"|||||spoiler|||",
	} {
		out, err := tg_md2html.ReverseV2(tg_md2html.MD2HTMLV2(test), nil)
		assert.NoError(t, err, "Error for:\n%s", test)
		assert.Equal(t, tg_md2html.MD2HTMLV2(test), tg_md2html.MD2HTMLV2(out))
	}
}

func TestReverseV2Buttons(t *testing.T) {
	for _, x := range []struct {
		in   string
		out  string
		btns []tg_md2html.ButtonV2
	}{
		{
			in:  "[hello](buttonurl://test.com)",
			out: "",
			btns: []tg_md2html.ButtonV2{{
				Name:    "hello",
				Content: "test.com",
			}},
		}, {
			in:  "Some text, some *bold*, and a button\n[hello](buttonurl://test.com)",
			out: "Some text, some <b>bold</b>, and a button",
			btns: []tg_md2html.ButtonV2{{
				Name:    "hello",
				Content: "test.com",
			}},
		},
	} {

		txt, b := tg_md2html.MD2HTMLButtonsV2(x.in)
		txt = strings.TrimSpace(txt)
		assert.Equal(t, x.out, txt)
		assert.ElementsMatch(t, x.btns, b)
		out, err := tg_md2html.ReverseV2(txt, x.btns)
		assert.NoError(t, err, "no error expected")
		assert.Equal(t, x.in, strings.TrimSpace(out))
	}
}
