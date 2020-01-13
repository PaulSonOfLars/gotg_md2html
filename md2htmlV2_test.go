package tg_md2html

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMD2HTMLV2Basic(t *testing.T) {
	for _, x := range append(basicMD) {
		assert.Equal(t, x.out, MD2HTMLV2(x.in))
	}
	// new mdv2 stuff
	for _, x := range []struct {
		in  string
		out string
	}{
		{
			in:  "~hello~",
			out: "<s>hello</s>",
		}, {
			in:  "__hello__",
			out: "<u>hello</u>",
		}, {
			in:  "```hello```",
			out: "<pre>hello</pre>",
			// NOTE: Decide on whether this is just a sad casualty of markdown parsing, or if:
			//  The closing tag should be the last viable part, if in a sequence. (eg 3x'_', last two are underline closes)
			//  This means that all other nested items of that tag should be escaped, to avoid:
			//  __underline  __double underline____, which is impossible. The HTML for this should be
			//  <u>underline __double underline(__)</u>(__) where () are up to opinion.
			//  Following my opinion, it should be the first.
		}, {
			in:  "___italic underline___",
			out: "<u><i>italic underline</i></u>",
		}, {
			in:  "__underline __double____",
			out: "<u>underline <u>double</u></u>",
		}, {
			in:  "__underline \\_\\_not double____",
			out: "<u>underline __not double__</u>",
		}, {
			in:  "____double underline____",
			out: "<u><u>double underline</u></u>",
		}, {
			in:  "````coded code block````",
			out: "<pre><code>coded code block</code></pre>",
		},
	} {
		assert.Equal(t, x.out, MD2HTMLV2(x.in))
	}
}

func TestMD2HTMLV2Advanced(t *testing.T) {
	for _, x := range advancedMD {
		assert.Equal(t, x.out, MD2HTMLV2(x.in))
	}
}

func TestNotMD2HTMLV2(t *testing.T) {
	for _, x := range []struct {
		in  string
		out string
	}{
		{
			in:  "hello",
			out: "hello",
		}, {
			in:  "_hello",
			out: "_hello",
		}, {
			in:  "hello_",
			out: "hello_",
		}, {
			in:  "_hello_there",
			out: "_hello_there",
		}, {
			in:  "_hello__",
			out: "<i>hello_</i>",
		}, {
			in:  "__hello__there",
			out: "__hello__there",
		}, {
			in:  "[hello](test.com)",
			out: `<a href="test.com">hello</a>`,
		},
	} {
		assert.Equal(t, x.out, MD2HTMLV2(x.in))
	}
}

func TestMD2HTMLV2Buttons(t *testing.T) {
	for _, x := range []struct {
		in   string
		out  string
		btns []ButtonV2
	}{
		{
			in:  "[hello](buttonurl:test.com)",
			out: "",
			btns: []ButtonV2{{
				Name:    "hello",
				Content: "test.com",
			}},
		}, {
			in:  "Some text, some *bold*, and a button [hello](buttonurl://test.com)",
			out: "Some text, some <b>bold</b>, and a button ",
			btns: []ButtonV2{{
				Name:    "hello",
				Content: "test.com",
			}},
		}, {
			in:   "[hello](buttonurl://test.com\\)",
			out:  "[hello](buttonurl://test.com)",
			btns: nil,
		}, {
			in:  "[hello](buttonurl://test.com\\)\n[hello2](buttonurl:test.com)",
			out: "[hello](buttonurl://test.com)\n",
			btns: []ButtonV2{{
				Name:    "hello2",
				Content: "test.com",
			}},
		},
	} {
		txt, b := MD2HTMLButtonsV2(x.in)
		assert.Equal(t, x.out, txt)
		assert.ElementsMatch(t, x.btns, b)
	}
}

func BenchmarkMD2HTMLV2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v, bs2 = MD2HTMLButtonsV2(message)
	}
}
