package tg_md2html_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	tg_md2html "github.com/PaulSonOfLars/gotg_md2html"
)

var basicMDv2 = []struct {
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
		in:  "||hello||",
		out: "<span class=\"tg-spoiler\">hello</span>",
	}, {
		in:  "||<hello>||",
		out: "<span class=\"tg-spoiler\">&lt;hello&gt;</span>",
	}, {
		in:  "```content```",
		out: "<pre>content</pre>",
	}, {
		in:  "```code\ncontent```",
		out: "<pre><code class=\"language-code\">content</code></pre>",
	}, {
		in:  "```spaced words\ncontent```",
		out: "<pre><code class=\"language-spaced words\">content</code></pre>",
	}, {
		in:  "```quoted\"words\ncontent```",
		out: "<pre><code class=\"language-quoted&#34;words\">content</code></pre>",
	}, {
		// NOTE: Decide on whether this is just a sad casualty of markdown parsing, or if:
		//  The closing tag should be the last viable part, if in a sequence. (eg 3x'_', last two are underline closes)
		//  This means that all other nested items of that tag should be escaped, to avoid:
		//  __underline  __double underline____, which is impossible. The HTML for this should be
		//  <u>underline __double underline(__)</u>(__) where () are up to opinion.
		//  Following my opinion, it should be the first.
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
		in:  "[hello](test.com)",
		out: `<a href="test.com">hello</a>`,
	}, {
		in:  "inline[url](test.com)test",
		out: `inline<a href="test.com">url</a>test`,
	}, {
		// pre and code dont support nested, so we dont parse the nested data.
		in:  "````coded code block````",
		out: "<pre>`coded code block`</pre>",
	}, { // ensure that premium stickers can get converted
		in:  `![👍](tg://emoji?id=5368324170671202286)`,
		out: `<tg-emoji emoji-id="5368324170671202286">👍</tg-emoji>`,
	}, {
		in:  "> quote",
		out: "<blockquote>quote</blockquote>",
	}, {
		in:  ">multi\n> line",
		out: "<blockquote>multi\nline</blockquote>",
	}, {
		in:  ">expandable multi\n>line\n>quote||",
		out: "<blockquote expandable>expandable multi\nline\nquote</blockquote>",
	}, {
		in:  ">expandable multi\n>line\n>quote||\nMore text on another line",
		out: "<blockquote expandable>expandable multi\nline\nquote</blockquote>\nMore text on another line",
	}, {
		in:  "**>expandable multi with star prefix\n>line\n>quote||",
		out: "<blockquote expandable>expandable multi with star prefix\nline\nquote</blockquote>",
	}, {
		in:  ">normal quote\n**>expandable multi\n>idk||",
		out: "<blockquote>normal quote</blockquote>\n<blockquote expandable>expandable multi\nidk</blockquote>",
	},
}

func TestMD2HTMLV2Basic(t *testing.T) {
	for _, x := range append(basicMD, basicMDv2...) {
		t.Run(x.in, func(t *testing.T) {
			assert.Equal(t, x.out, tg_md2html.MD2HTMLV2(x.in))
		})
	}
}

func TestMD2HTMLV2Advanced(t *testing.T) {
	for _, x := range advancedMD {
		t.Run(x.in, func(t *testing.T) {
			assert.Equal(t, x.out, tg_md2html.MD2HTMLV2(x.in))
		})
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
			in:  "||bad spoiler",
			out: "||bad spoiler",
		}, {
			in:  "|noop|",
			out: "|noop|",
		}, {
			in:  "end with >",
			out: "end with &gt;",
		}, {
			in:  "no premium ! in text", // confirm that a '!' doesnt break premiums
			out: "no premium ! in text",
		}, {
			in:  "no premium!", // confirm that ending with '!' doesn't break premiums
			out: "no premium!",
		}, {
			in:  `test !`, // Also check ' !'
			out: `test !`,
		}, {
			// Premium Emoji matching should be greedy
			in:  "Some text ![😎](tg://emoji?id=6026150900848923116)andstuckhere![😎](tg://emoji?id=6026150900848923116), more text",
			out: "Some text <tg-emoji emoji-id=\"6026150900848923116\">😎</tg-emoji>andstuckhere<tg-emoji emoji-id=\"6026150900848923116\">😎</tg-emoji>, more text",
		},
	} {
		t.Run(x.in, func(t *testing.T) {
			txt := tg_md2html.MD2HTMLV2(x.in)
			assert.Equal(t, x.out, txt)
		})
	}
}

var md2HTMLV2Buttons = []struct {
	in   string
	out  string
	btns []tg_md2html.ButtonV2
}{
	{
		in:  "[hello](buttonurl:test.com)",
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
	}, {
		in:   "[hello](buttonurl://test.com\\)",
		out:  "[hello](buttonurl://test.com)",
		btns: nil,
	}, {
		in:  "[hello](buttonurl://test.com\\)\n[hello2](buttonurl:test.com)",
		out: "",
		btns: []tg_md2html.ButtonV2{{
			Name:    "hello](buttonurl://test.com)\n[hello2",
			Type:    "url",
			Content: "test.com",
		}},
	}, {
		in:  "[text](buttontext:This is some basic text)\n[hello2](buttonurl:test.com)",
		out: "",
		btns: []tg_md2html.ButtonV2{{
			Name:    "text",
			Type:    "text",
			Content: "This is some basic text",
		}, {
			Name:    "hello2",
			Type:    "url",
			Content: "test.com",
		}},
	}, {
		// This is not a valid URL
		in:   "[text](tg://emoji?id=6026150900848923116)",
		out:  "[text](tg://emoji?id=6026150900848923116)",
		btns: nil,
	}, {
		in:  "Some text [![😎](tg://emoji?id=6026150900848923116) text](buttonurl://example.com)",
		out: "Some text",
		btns: []tg_md2html.ButtonV2{
			{
				Name:     "😎 text",
				Type:     "url",
				Content:  "example.com",
				SameLine: false,
			},
		},
	}, {
		in:  "[![🌐](tg://emoji?id=5343789187172670307)Website](buttonurl://example.com)",
		out: "",
		btns: []tg_md2html.ButtonV2{
			{
				Name:     "🌐Website",
				Type:     "url",
				Content:  "example.com",
				SameLine: false,
			},
		},
	}, {
		// This one has a space between emoji and website! which causes different output...!
		in:  "[![🌐](tg://emoji?id=5343789187172670307) Website](buttonurl://example.com)",
		out: "",
		btns: []tg_md2html.ButtonV2{
			{
				Name:     "🌐 Website",
				Type:     "url",
				Content:  "example.com",
				SameLine: false,
			},
		},
	}, {
		// Handle the edge case where people are purposefully trying to break buttons
		in:  "[![🌐](buttonurl://example.com)",
		out: "",
		btns: []tg_md2html.ButtonV2{
			{
				Name:     "![🌐",
				Type:     "url",
				Content:  "example.com",
				SameLine: false,
			},
		},
	}, {
		in:  "text\n> quote\ntext",
		out: "text\n<blockquote>quote</blockquote>\ntext",
	}, {
		in:  "> `code quote`",
		out: "<blockquote><code>code quote</code></blockquote>",
	}, {
		in:  "```go\ntext\n> not quote\nmore text```",
		out: "<pre><code class=\"language-go\">text\n&gt; not quote\nmore text</code></pre>",
	}, {
		in: `[text](buttonurl#primary://example.com)
[text](buttonurl#success://example.com)
[text](buttonurl#danger://example.com)`,
		out: "",
		btns: []tg_md2html.ButtonV2{{
			Name:     "text",
			Type:     "url",
			Content:  "example.com",
			SameLine: false,
			Style:    "primary",
		}, {
			Name:     "text",
			Type:     "url",
			Content:  "example.com",
			SameLine: false,
			Style:    "success",
		}, {
			Name:     "text",
			Type:     "url",
			Content:  "example.com",
			SameLine: false,
			Style:    "danger",
		}},
	},
}

func TestMD2HTMLV2Buttons(t *testing.T) {
	for _, x := range md2HTMLV2Buttons {
		cv := testConverter()
		t.Run(x.in, func(t *testing.T) {
			txt, b := cv.MD2HTMLButtons(x.in)
			assert.Equal(t, x.out, txt)
			assert.ElementsMatch(t, x.btns, b)
		})
	}
}

func BenchmarkMD2HTMLV2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v, bs2 = tg_md2html.MD2HTMLButtonsV2(message)
	}
}
