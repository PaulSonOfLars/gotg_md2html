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

	for _, test := range append(append(basicMD, basicMDv2...), advancedMD...) {
		out, err := tg_md2html.ReverseV2(tg_md2html.MD2HTMLV2(test.in), nil)
		assert.NoError(t, err, "Error for:\n%s", test)
		assert.Equal(t, tg_md2html.MD2HTMLV2(test.in), tg_md2html.MD2HTMLV2(out))
	}

	for _, test := range []string{
		"___________test_______",
		"|||||spoiler|||",
		`![👍](tg://emoji?id=5368324170671202286)`,
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
				Type:    "url",
				Content: "test.com",
			}},
		}, {
			in:  "Some text, some *bold*, and a button\n[hello](buttonurl://test.com)",
			out: "Some text, some <b>bold</b>, and a button",
			btns: []tg_md2html.ButtonV2{{
				Name:    "hello",
				Type:    "url",
				Content: "test.com",
			}},
		}, {
			in:  "Some text, some *bold*, and a button\n[hello](buttontext://some text)",
			out: "Some text, some <b>bold</b>, and a button",
			btns: []tg_md2html.ButtonV2{{
				Name:    "hello",
				Type:    "text",
				Content: "some text",
			}},
		}, {
			in:  "Some text, some *bold*, and a button\n[hello](buttontext://some text:same)",
			out: "Some text, some <b>bold</b>, and a button",
			btns: []tg_md2html.ButtonV2{{
				Name:     "hello",
				Type:     "text",
				Content:  "some text",
				SameLine: true,
			}},
		},
	} {

		cv := tg_md2html.NewV2(map[string]string{
			"url":  "buttonurl:",
			"text": "buttontext:",
		})

		txt, b := cv.MD2HTMLButtons(x.in)
		txt = strings.TrimSpace(txt)
		assert.Equal(t, x.out, txt)
		assert.ElementsMatch(t, x.btns, b)
		out, err := cv.Reverse(txt, x.btns)
		assert.NoError(t, err, "no error expected")
		assert.Equal(t, x.in, strings.TrimSpace(out))
	}
}
