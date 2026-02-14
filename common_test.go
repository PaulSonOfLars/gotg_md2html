package tg_md2html_test

import "github.com/PaulSonOfLars/gotg_md2html"

func testConverter() *tg_md2html.ConverterV2 {
	return tg_md2html.NewV2(map[string]string{
		"url":  "buttonurl",
		"text": "buttontext",
	}, map[string]string{
		"primary": "primary",
		"success": "success",
		"danger":  "danger",
	})
}
