package tg_md2html

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMD2HTML(t *testing.T) {
	type mdTestStruct struct {
		input  string
		output string
	}
	for _, test := range []mdTestStruct{
		{
			input:  "hello there",
			output: "hello there",
		}, {
			input:  "_hello_ there",
			output: "<i>hello</i> there",
		}, {
			input:  "hello _there_",
			output: "hello <i>there</i>",
		}, {
			input:  "_hello there_",
			output: "<i>hello there</i>",
		}, {
			input:  "_hello_ there_",
			output: "<i>hello</i> there_",
		}, {
			input:  "_hello _there_",
			output: "<i>hello _there</i>",
		}, {
			input:  "_hello _ there_",
			output: "<i>hello _ there</i>",
		}, {
			input:  "so_hello _there_",
			output: "so_hello <i>there</i>",
		}, {
			input:  "_hello you_there_",
			output: "<i>hello you_there</i>",
		}, {
			input:  "`hello` there",
			output: "<code>hello</code> there",
		}, {
			input:  "*hello* there",
			output: "<b>hello</b> there",
		}, {
			input:  "hello [there](link.com)",
			output: `hello <a href="link.com">there</a>`,
		}, {
			input:  "hello [there](buttonurl://link.com)",
			output: `hello <a href="buttonurl://link.com">there</a>`,
		}, {
			input:  "hello [there[]](link.com)",
			output: `hello <a href="link.com">there[]</a>`,
		}, {
			input:  "[hello] soo] () [there](link.com)",
			output: `<a href="link.com">hello] soo] () [there</a>`,
		}, {
			input:  "_hello_ `there` *bold* [url](link.com) _`notcode`_ *_notitalic_* [weird not italic _](morelink.co.uk)_",
			output: "<i>hello</i> <code>there</code> <b>bold</b> <a href=\"link.com\">url</a> <i>`notcode`</i> <b>_notitalic_</b> <a href=\"morelink.co.uk\">weird not italic _</a>_",
		}, {
			input:  "[hello] soo] () [there](link.com)",
			output: `<a href="link.com">hello] soo] () [there</a>`,
		}, {
			input:  "]]]]]]] )))))))  ((((([link](example.com) [link2](example2.com) [link3](example3.com) ]]]]](((())))",
			output: `]]]]]]] )))))))  (((((<a href="example.com">link</a> <a href="example2.com">link2</a> <a href="example3.com">link3</a> ]]]]](((())))`,
		}, {
			input:  "[reallybiglink\\](example3.com) [insidelink](exampleLink.com)",
			output: `<a href="exampleLink.com">reallybiglink](example3.com) [insidelink</a>`,
		}, {
			input:  "[link](example.com) [escapedlink2]\\(example2.com) \\[escapedlink3](example3.com) [notalink] (no.com) [reallybiglink\\](example3.com) [insidelink](example3.com)",
			output: `<a href="example.com">link</a> <a href="example3.com">escapedlink2](example2.com) [escapedlink3</a> <a href="example3.com">notalink] (no.com) [reallybiglink](example3.com) [insidelink</a>`,
		}, {
			input:  "hello there _friend_ how * are _ you? [link[with a sub box!]](example.com) emoji [emoji link ](example.com)",
			output: `hello there <i>friend</i> how * are _ you? <a href="example.com">link[with a sub box!]</a> emoji <a href="example.com">emoji link </a>`,
		}, {
			input:  "_hello_1",
			output: "_hello_1",
		}, {
			input:  `*\**`,
			output: "<b>*</b>",
		}, {
			input:  "hell_o [there[]](link.com/this_isfine)",
			output: `hell_o <a href="link.com/this_isfine">there[]</a>`,
		},
	} {
		assert.Equal(t, test.output, MD2HTML(test.input))
	}
}

func TestMD2HTMLButtons(t *testing.T) {
	type mdTestStruct struct {
		input  string
		output string
		btns   []Button
	}

	for _, test := range []mdTestStruct{
		{
			input:  "hello [there](buttonurl://link.com)",
			output: `hello`,
			btns: []Button{{
				Name:     "there",
				Content:  "link.com",
				SameLine: false,
			}},
		}, {
			input:  "hey there! My name is @MissRose_bot. go to [Rules](buttonurl://t.me/MissRose_bot?start=idek_12345)",
			output: "hey there! My name is @MissRose_bot. go to",
			btns: []Button{{
				Name:     "Rules",
				Content:  "t.me/MissRose_bot?start=idek_12345",
				SameLine: false,
			}},
		}, {
			input:  "no [1](buttonurl://link.com)[2](buttonurl://link.com)[3](buttonurl://link.com)",
			output: `no`,
			btns: []Button{{
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
			btns: []Button{{
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
			btns: []Button{{
				Name:     "link",
				Content:  "link.com",
				SameLine: false,
			}},
		},
	} {
		out, btns := MD2HTMLButtons(test.input)
		assert.Equal(t, test.output, out)
		assert.ElementsMatch(t, test.btns, btns)
	}
}

func TestReverse(t *testing.T) {
	for _, test := range []string{
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
	} {
		assert.Equal(t, MD2HTML(test), MD2HTML(Reverse(MD2HTML(test), nil)))
	}
}

func TestReverseBtns(t *testing.T) {
	type TestRevBtn struct {
		text    string
		buttons []Button
		out     string
	}

	for _, test := range []TestRevBtn{
		{
			text:    "Hello there",
			buttons: []Button{},
			out:     "Hello there",
		}, {
			text:    "Hello there <i>italic</i>",
			buttons: []Button{},
			out:     "Hello there _italic_",
		}, {
			text: "Hello there",
			buttons: []Button{
				{
					Name:     "Test",
					Content:  "link.com",
					SameLine: false,
				},
			},
			out: "Hello there\n[Test](buttonurl://link.com)",
		}, {
			text: "oh no",
			buttons: []Button{
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
			out:     "I dont even knowww \\\\[ stuff",
		}, {
			text:    "I dont even knowww \\\\[ stuff",
			buttons: nil,
			out:     "I dont even knowww \\\\\\\\[ stuff",
		}, {
			text: "Hello there",
			buttons: []Button{
				{
					Name:     "test with ' quote",
					Content:  "link.com",
					SameLine: false,
				},
			},
			out: "Hello there\n[test with ' quote](buttonurl://link.com)",
		}, {
			text: "Hello there",
			buttons: []Button{
				{
					Name:     "test with ' quote",
					Content:  "link.com%22%22%22",
					SameLine: false,
				},
			},
			out: "Hello there\n[test with ' quote](buttonurl://link.com%22%22%22)",
		}, {
			text: "Hello there",
			buttons: []Button{
				{
					Name:     "test with ' quote",
					Content:  "link.com\"\"",
					SameLine: false,
				},
			},
			out: "Hello there\n[test with ' quote](buttonurl://link.com\"\")",
		},
	} {
		assert.Equal(t, test.out, Reverse(test.text, test.buttons))

		one, oneb := MD2HTMLButtons(test.out)
		two, twob := MD2HTMLButtons(Reverse(one, oneb))
		assert.Equal(t, one, two)
		assert.ElementsMatch(t, oneb, twob)
	}
}
