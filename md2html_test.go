package tg_md2html_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	tg_md2html "github.com/PaulSonOfLars/gotg_md2html"
)

var basicMD = []struct {
	in  string
	out string
}{
	{
		in:  "hello",
		out: "hello",
	}, {
		in:  "_hello_",
		out: "<i>hello</i>",
	}, {
		in:  "*hello*",
		out: "<b>hello</b>",
	}, {
		in:  "`hello`",
		out: "<code>hello</code>",
	}, {
		in:  "[hello](test.com)",
		out: `<a href="test.com">hello</a>`,
	}, {
		in:  `_hello\_ there_`,
		out: "<i>hello_ there</i>",
	}, {
		in:  "`code and <brackets>`",
		out: "<code>code and &lt;brackets&gt;</code>",
	},
}

var advancedMD = []struct {
	in  string
	out string
}{
	{
		in:  "hello there",
		out: "hello there",
	}, {
		in:  "_hello_ there",
		out: "<i>hello</i> there",
	}, {
		in:  "hello _there_",
		out: "hello <i>there</i>",
	}, {
		in:  "_hello there_",
		out: "<i>hello there</i>",
	}, {
		in:  "_hello_ there_",
		out: "<i>hello</i> there_",
	}, {
		in:  "_hello _there_",
		out: "<i>hello _there</i>",
	}, {
		in:  "_hello _ there_",
		out: "<i>hello _ there</i>",
	}, {
		in:  "so_hello _there_",
		out: "so_hello <i>there</i>",
	}, {
		in:  "_hello you_there_",
		out: "<i>hello you_there</i>",
	}, {
		in:  "`hello` there",
		out: "<code>hello</code> there",
	}, {
		in:  "*hello* there",
		out: "<b>hello</b> there",
	}, {
		in:  "hello [there](link.com)",
		out: `hello <a href="link.com">there</a>`,
	}, {
		in:  "hello [there](buttonurl://link.com)",
		out: `hello <a href="buttonurl://link.com">there</a>`,
	}, {
		in:  "hello [there[]](link.com)",
		out: `hello <a href="link.com">there[]</a>`,
	}, {
		in:  "[hello] soo] () [there](link.com)",
		out: `<a href="link.com">hello] soo] () [there</a>`,
	}, {
		in:  "[hello] soo] () [there](link.com)",
		out: `<a href="link.com">hello] soo] () [there</a>`,
	}, {
		in:  "]]]]]]] )))))))  ((((([link](example.com) [link2](example2.com) [link3](example3.com) ]]]]](((())))",
		out: `]]]]]]] )))))))  (((((<a href="example.com">link</a> <a href="example2.com">link2</a> <a href="example3.com">link3</a> ]]]]](((())))`,
	}, {
		in:  "[reallybiglink\\](example3.com) [insidelink](exampleLink.com)",
		out: `<a href="exampleLink.com">reallybiglink](example3.com) [insidelink</a>`,
	}, {
		in:  "[link](example.com) [escapedlink2]\\(example2.com) \\[escapedlink3](example3.com) [notalink] (no.com) [reallybiglink\\](example3.com) [insidelink](example3.com)",
		out: `<a href="example.com">link</a> <a href="example3.com">escapedlink2](example2.com) [escapedlink3</a> <a href="example3.com">notalink] (no.com) [reallybiglink](example3.com) [insidelink</a>`,
	}, {
		in:  "hello there _friend_ how * are _ you? [link[with a sub box!]](example.com) emoji [emoji link ](example.com)",
		out: `hello there <i>friend</i> how * are _ you? <a href="example.com">link[with a sub box!]</a> emoji <a href="example.com">emoji link </a>`,
	}, {
		in:  "_hello_1",
		out: "_hello_1",
	}, {
		in:  `*\**`,
		out: "<b>*</b>",
	}, {
		in:  "hell_o [there[]](link.com/this_isfine)",
		out: `hell_o <a href="link.com/this_isfine">there[]</a>`,
	}, {
		in:  "\\",
		out: "\\",
	}, {
		in:  "_| _ '_ ` _ _`",
		out: "<i>| _ &#39;</i> ` _ _`",
	}, {
		in:  "_| _ '_ ` _ _` _| _ '_ ` _ _` _| _ '_ ` _ _` _| _ '_ ` _ _`",
		out: "<i>| _ &#39;</i> ` _ <i>` _| _ &#39;</i> ` _ <i>` _| _ &#39;</i> ` _ <i>` _| _ &#39;</i> ` _ _`",
	}, {
		in:  "_hey\\_ there_",
		out: "<i>hey_ there</i>",
	}, {
		in:  "\\_input_",
		out: "_input_",
	},
}

var reverseTest = []string{
	"hello there",
	"_hello_ there",
	"hello _there_",
	"_hello there_",
	"_hello_ there_",
	"_hello _there_",
	"_hello _ there_",
	"so_hello _there_",
	"_hello you_there_",
	"`hello` there",
	"*hello* there",
	"hello [there](link.com)",
	"hello [there](buttonurl://link.com)",
	"hello [there[]](link.com)",
	"[hello] soo] () [there](link.com)",
	"_hello_ `there` *bold* [url](link.com) _`notcode`_ *_notitalic_* [weird not italic _](morelink.co.uk)_",
	"[hello] soo] () [there](link.com)",
	"]]]]]]] )))))))  ((((([link](example.com) [link2](example2.com) [link3](example3.com) ]]]]](((())))",
	"[reallybiglink\\](example3.com) [insidelink](exampleLink.com)",
	"[link](example.com) [escapedlink2]\\(example2.com) \\[escapedlink3](example3.com) [notalink] (no.com) [reallybiglink\\](example3.com) [insidelink](example3.com)",
	"hello there _friend_ how * are _ you? [link[with a sub box!]](example.com) emoji [emoji link ](example.com)",
	"_hello_1",
	`*\**`,
	"hell_o [there[]](link.com/this_isfine)",
	"hello \\[there](link.com)",
	"hello \\[there\\]\\(link.com\\)",
	"hello \\[there]\\(link.com)",
	"hello [there more \\[text](link.com)",
	"hello \\(link.com)",
	"hello \\\\(who knowsss\\)",
	"hello \\[ text",
	"don't know",
	"don't \\[ text",
	"don't hey_o_there whastup \"",
	"don't hey_o_there [link](idek.com)",
	"don't hey_o_there [link](buttonurl://idek.com)",
	"don't hey_o_there \\[link](buttonurl://idek.com)",
	"don't hey_o_there \\[link](buttonurl://idek.com) [nolink]",
	"don't hey_o_there \\[link](buttonurl://idek.com) [stillink](buttonurl://test.com)",
	"don't _hey'quotes_",
	"`test \\`escaped` backticks",
	"`test\\` escaped` backticks",
	"\\_test_",
	"\\\\_test_",
	"_italics\\_ *maybebold* still italics_",
	"`code and <brackets>`",
	"_italics and <brackets>_",
	"__strikethrough__", // mdv1 wont parse this, but thats ok
	"||spoiler||",       // mdv1 wont parse this, but thats ok
}

func TestMD2HTMLBasic(t *testing.T) {
	for _, x := range basicMD {
		assert.Equal(t, x.out, tg_md2html.MD2HTML(x.in))
	}
}

func TestMD2HTMLAdvanced(t *testing.T) {
	for _, test := range advancedMD {
		assert.Equal(t, test.out, tg_md2html.MD2HTML(test.in))
	}

	assert.Equal(t,
		"<i>hello</i> <code>there</code> <b>bold</b> <a href=\"link.com\">url</a> <i>`notcode`</i> <b>_notitalic_</b> <a href=\"morelink.co.uk\">weird not italic _</a>_",
		tg_md2html.MD2HTML("_hello_ `there` *bold* [url](link.com) _`notcode`_ *_notitalic_* [weird not italic _](morelink.co.uk)_"),
	)
}

func TestMD2HTMLButtons(t *testing.T) {
	type mdTestStruct struct {
		input  string
		output string
		btns   []tg_md2html.Button
	}

	for _, test := range []mdTestStruct{
		{
			input:  "hello [there](buttonurl://link.com)",
			output: `hello`,
			btns: []tg_md2html.Button{{
				Name:     "there",
				Content:  "link.com",
				SameLine: false,
			}},
		}, {
			input:  "hey there! My name is @MissRose_bot. go to [Rules](buttonurl://t.me/MissRose_bot?start=idek_12345)",
			output: "hey there! My name is @MissRose_bot. go to",
			btns: []tg_md2html.Button{{
				Name:     "Rules",
				Content:  "t.me/MissRose_bot?start=idek_12345",
				SameLine: false,
			}},
		}, {
			input:  "no [1](buttonurl://link.com)[2](buttonurl://link.com)[3](buttonurl://link.com)",
			output: `no`,
			btns: []tg_md2html.Button{{
				Name:     "1",
				Content:  "link.com",
				SameLine: false,
			}, {
				Name:     "2",
				Content:  "link.com",
				SameLine: false,
			}, {
				Name:     "3",
				Content:  "link.com",
				SameLine: false,
			}},
		}, {
			input:  "*bold [box]* [1](buttonurl://link.com)[2](buttonurl://link.com)[3](buttonurl://link.com)",
			output: `<b>bold [box]</b>`,
			btns: []tg_md2html.Button{{
				Name:     "1",
				Content:  "link.com",
				SameLine: false,
			}, {
				Name:     "2",
				Content:  "link.com",
				SameLine: false,
			}, {
				Name:     "3",
				Content:  "link.com",
				SameLine: false,
			}},
		}, {
			input:  "*a [box]*[link](buttonurl://link.com)",
			output: "<b>a [box]</b>",
			btns: []tg_md2html.Button{{
				Name:     "link",
				Content:  "link.com",
				SameLine: false,
			}},
		},
	} {
		out, btns := tg_md2html.MD2HTMLButtons(test.input)
		assert.Equal(t, test.output, out)
		assert.ElementsMatch(t, test.btns, btns)
	}
}

func TestReverse(t *testing.T) {
	for _, test := range reverseTest {
		assert.Equal(t, tg_md2html.MD2HTML(test), tg_md2html.MD2HTML(tg_md2html.Reverse(tg_md2html.MD2HTML(test), nil)))
	}
}

func TestReverseBtns(t *testing.T) {
	type TestRevBtn struct {
		text    string
		buttons []tg_md2html.Button
		out     string
	}

	for _, test := range []TestRevBtn{
		{
			text:    "Hello there",
			buttons: []tg_md2html.Button{},
			out:     "Hello there",
		}, {
			text:    "Hello there <i>italic</i>",
			buttons: []tg_md2html.Button{},
			out:     "Hello there _italic_",
		}, {
			text: "Hello there",
			buttons: []tg_md2html.Button{
				{
					Name:     "Test",
					Content:  "link.com",
					SameLine: false,
				},
			},
			out: "Hello there\n[Test](buttonurl://link.com)",
		}, {
			text: "oh no",
			buttons: []tg_md2html.Button{
				{
					Name:     "btn1",
					Content:  "example.com",
					SameLine: false,
				}, {
					Name:     "btn2",
					Content:  "link.com",
					SameLine: true,
				},
			},
			out: "oh no\n[btn1](buttonurl://example.com)\n[btn2](buttonurl://link.com:same)",
		}, {
			text:    "I dont even knowww \\[ stuff",
			buttons: nil,
			out:     "I dont even knowww \\\\\\[ stuff", // -> \ (e:\\) becomes \\ (e:\\\\) and [ becomes \] (e:\\[) so \\\\\\[
		}, {
			text:    "I dont even knowww \\\\[ stuff",
			buttons: nil,
			out:     "I dont even knowww \\\\\\\\\\[ stuff",
		}, {
			text: "Hello there",
			buttons: []tg_md2html.Button{
				{
					Name:     "test with ' quote",
					Content:  "link.com",
					SameLine: false,
				},
			},
			out: "Hello there\n[test with ' quote](buttonurl://link.com)",
		}, {
			text: "Hello there",
			buttons: []tg_md2html.Button{
				{
					Name:     "test with ' quote",
					Content:  "link.com%22%22%22",
					SameLine: false,
				},
			},
			out: "Hello there\n[test with ' quote](buttonurl://link.com%22%22%22)",
		}, {
			text: "Hello there",
			buttons: []tg_md2html.Button{
				{
					Name:     "test with ' quote",
					Content:  "link.com\"\"",
					SameLine: false,
				},
			},
			out: "Hello there\n[test with ' quote](buttonurl://link.com\"\")",
		},
	} {
		assert.Equal(t, test.out, tg_md2html.Reverse(test.text, test.buttons))

		one, oneb := tg_md2html.MD2HTMLButtons(test.out)
		two, twob := tg_md2html.MD2HTMLButtons(tg_md2html.Reverse(one, oneb))
		assert.Equal(t, one, two)
		assert.ElementsMatch(t, oneb, twob)
	}
}

func TestIsEscaped(t *testing.T) {
	for _, x := range []struct {
		s string
		b bool
	}{
		{
			s: "a",
			b: false,
		}, {
			s: `\a`,
			b: true,
		}, {
			s: `\\a`,
			b: false,
		}, {
			s: `\\\a`,
			b: true,
		}, {
			s: `Hey there a`,
			b: false,
		}, {
			s: `Hey there \a`,
			b: true,
		}, {
			s: `Hey there \\a`,
			b: false,
		},
	} {
		assert.Equal(t, x.b, tg_md2html.IsEscaped([]rune(x.s), len([]rune(x.s[:strings.IndexRune(x.s, 'a')]))))
	}
}

var stripMD = []struct {
	in     string
	output string
}{
	{
		in:     "this is text",
		output: "this is text",
	}, {
		in:     "_italic_",
		output: "italic",
	}, {
		in:     "*bold*",
		output: "bold",
	}, {
		in:     "`code`",
		output: "code",
	}, {
		in:     "[link](test.html)",
		output: "link",
	}, {
		in:     "_test * italics_",
		output: "test * italics",
	}, {
		in:     "i dont *even* _know_ `why *_ you would`",
		output: "i dont even know why *_ you would",
	},
}

func TestStripMD(t *testing.T) {
	for _, x := range stripMD {
		assert.Equal(t, x.output, tg_md2html.StripMD(x.in), "failed to strip all markdown")
	}
}

var v string
var bs []tg_md2html.Button
var bs2 []tg_md2html.ButtonV2

func BenchmarkMD2HTML(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v, bs = tg_md2html.MD2HTMLButtons(message)
	}
}

var message = `
This is text.
There is *bold*
There is _italic_
There is a [url](example.com)
There is a button [button](buttonurl://example.com)
There is some more *bold*
And some _italic_
Good test?
`
