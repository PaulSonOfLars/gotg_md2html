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
			out: "<i>hello</i>_",
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
