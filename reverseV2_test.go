package tg_md2html_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	tg_md2html "github.com/PaulSonOfLars/gotg_md2html"
)

func TestReverseV2(t *testing.T) {
	for _, test := range reverseTest {
		t.Run(test, func(t *testing.T) {
			out, err := tg_md2html.ReverseV2(tg_md2html.MD2HTMLV2(test), nil)
			assert.NoError(t, err, "Error for:\n%s", test)
			assert.Equal(t, tg_md2html.MD2HTMLV2(test), tg_md2html.MD2HTMLV2(out))
		})
	}

	for _, test := range append(append(basicMD, basicMDv2...), advancedMD...) {
		t.Run(test.in, func(t *testing.T) {
			out, err := tg_md2html.ReverseV2(tg_md2html.MD2HTMLV2(test.in), nil)
			assert.NoError(t, err, " Error for:\n%s", test)
			assert.Equal(t, tg_md2html.MD2HTMLV2(test.in), tg_md2html.MD2HTMLV2(out))
		})
	}

	for _, test := range []string{
		"___________test_______",                  // uneven underlines
		"|||||spoiler|||",                         // uneven spoilers
		"||<spoiler>||",                           // spoilers, but with HTML bits inside
		"![👍](tg://emoji?id=5368324170671202286)", // premium emoji
		"> ",             // empty quotes
		"test\n>\ntest",  // multiline quotes
		"||||||||| test", // nested spoilers
	} {
		t.Run(test, func(t *testing.T) {
			htmlv2 := tg_md2html.MD2HTMLV2(test)
			out, err := tg_md2html.ReverseV2(htmlv2, nil)
			assert.NoError(t, err, "Error for:\n%s", test)
			assert.Equal(t, htmlv2, tg_md2html.MD2HTMLV2(out))
		})
	}
}

func TestReverseV2Buttons(t *testing.T) {
	for _, x := range md2HTMLV2Buttons {
		t.Run(x.in, func(t *testing.T) {
			cv := testConverter()

			// Button parsing is lossful; we apply opinionated changes to ensure stability. (eg, strip markdown from button names)
			// So, we don't compare the reverse to the input.
			// Instead, we ensure that the parsed input is equal to the parsed reverse.
			txt, b := cv.MD2HTMLButtons(x.in)
			assert.Equal(t, x.out, txt)
			assert.ElementsMatch(t, x.btns, b)

			out, err := cv.Reverse(txt, x.btns)
			assert.NoError(t, err, "no error expected")

			txt2, b2 := cv.MD2HTMLButtons(out)
			assert.Equal(t, txt, txt2)
			assert.ElementsMatch(t, b, b2)
		})
	}
}

func TestReverseV2Buttons_errors(t *testing.T) {
	for _, x := range []struct {
		name    string
		in      string
		wantErr error
	}{
		{
			name:    "valid button",
			in:      "[test](buttonurl://example.com)",
			wantErr: nil,
		}, {
			name:    "empty button name",
			in:      "[](buttonurl://example.com)",
			wantErr: tg_md2html.ErrNoButtonContent,
		}, {
			name:    "empty button url",
			in:      "[test](buttonurl://)",
			wantErr: tg_md2html.ErrNoButtonContent,
		}, {
			name:    "valid button style",
			in:      "[test](buttonurl#primary://blah)",
			wantErr: nil,
		}, {
			name:    "invalid button style",
			in:      "[test](buttonurl#green://blah)",
			wantErr: tg_md2html.ErrInvalidButtonStyle,
		},
	} {
		t.Run(x.name, func(t *testing.T) {
			cv := testConverter()

			// Button parsing is lossful; we apply opinionated changes to ensure stability. (eg, strip markdown from button names)
			// So, we don't compare the reverse to the input.
			// Instead, we ensure that the parsed input is equal to the parsed reverse.
			txt, b := cv.MD2HTMLButtons(x.in)
			_, err := cv.Reverse(txt, b)
			assert.ErrorIs(t, err, x.wantErr)
		})
	}
}
