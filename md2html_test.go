package tg_md2html

import (
	"testing"
	"github.com/magiconair/properties/assert"
)

type mdTestStruct struct {
	input  string
	output string
}

func TestMD2HTML(t *testing.T) {
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
			input: "[link](example.com) [escapedlink2]\\(example2.com) \\[escapedlink3](example3.com) [notalink] (no.com) [reallybiglink\\](example3.com) [insidelink](example3.com)",
			output: `<a href="example.com">link</a> <a href="example3.com">escapedlink2](example2.com) [escapedlink3</a> <a href="example3.com">notalink] (no.com) [reallybiglink](example3.com) [insidelink</a>`,
		},
	} {
		assert.Equal(t, MD2HTML(test.input), test.output)
	}
}
