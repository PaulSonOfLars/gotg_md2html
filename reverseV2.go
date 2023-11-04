package tg_md2html

import (
	"errors"
	"fmt"
	"html"
	"regexp"
	"strings"
)

var ErrNoButtonContent = errors.New("no button contents")

var languageCodeblock = regexp.MustCompile(`^(?s)<code class="language-(.*?)">(.*)</code>$`)

func ReverseV2(in string, bs []ButtonV2) (string, error) {
	return defaultConverterV2.Reverse(in, bs)
}

func (cv ConverterV2) Reverse(in string, bs []ButtonV2) (string, error) {
	text, err := cv.reverse([]rune(in), bs)
	return strings.TrimSpace(text), err
}

func (cv ConverterV2) reverse(in []rune, buttons []ButtonV2) (string, error) {
	prev := 0
	out := strings.Builder{}
	for i := 0; i < len(in); i++ {
		switch in[i] {
		case '<':
			c := getHTMLTagCloseIndex(in[i+1:])
			if c < 0 {
				// "no close tag"
				return "", fmt.Errorf("no closing '>' for opening bracket at %d", i)
			}
			closeTag := i + c + 1
			tagContent := string(in[i+1 : closeTag])
			tagFields := strings.Fields(tagContent)
			if len(tagFields) < 1 {
				return "", fmt.Errorf("no tag name for HTML tag started at %d", i)
			}
			tag := tagFields[0]

			co, cc := getClosingTag(in[closeTag+1:], tag)
			if co < 0 || cc < 0 {
				// "no closing open"
				return "", fmt.Errorf("no closing tag for HTML tag %q started at %d", tag, i)
			}
			closingOpen, closingClose := closeTag+1+co, closeTag+1+cc
			out.WriteString(html.UnescapeString(string(in[prev:i])))

			nested, err := cv.reverse(in[closeTag+1:closingOpen], nil)
			if err != nil {
				return "", err
			}

			switch tag {
			case "b", "strong":
				out.WriteString("*" + nested + "*")
			case "i", "em":
				out.WriteString("_" + nested + "_")
			case "u", "ins":
				out.WriteString("__" + nested + "__")
			case "s", "strike", "del":
				out.WriteString("~" + nested + "~")
			case "code":
				// code and pre don't look at nested values, because they're not parsed
				out.WriteString("`" + html.UnescapeString(string(in[closeTag+1:closingOpen])) + "`")
			case "pre":
				// code and pre don't look at nested values, because they're not parsed
				content := html.UnescapeString(string(in[closeTag+1 : closingOpen]))
				m := languageCodeblock.FindStringSubmatch(content)
				if len(m) > 0 {
					// This <pre> block contains a <code class...> block; handle the language.
					lang, code := m[1], m[2]
					out.WriteString("```" + lang + code + "```")
				} else {
					// This is a regular boring pre block
					out.WriteString("```" + content + "```")
				}
			case "span":
				// NOTE: All span tags are currently spoiler tags. This may change in the future.
				if len(tagFields) < 2 {
					return "", fmt.Errorf("span tag does not have enough fields %q", tagFields)
				}

				switch spanType := tagFields[1]; spanType {
				case "class=\"tg-spoiler\"":
					out.WriteString("||" + html.UnescapeString(string(in[closeTag+1:closingOpen])) + "||")
				default:
					return "", fmt.Errorf("unknown tag type %q", spanType)
				}
			case "a":
				if link.MatchString(tagContent) {
					matches := link.FindStringSubmatch(tagContent)
					out.WriteString("[" + nested + "](" + matches[1] + ")")
				} else {
					return "", fmt.Errorf("badly formatted anchor tag %q", tagContent)
				}
			case "tg-emoji":
				if customEmoji.MatchString(tagContent) {
					matches := customEmoji.FindStringSubmatch(tagContent)
					out.WriteString("![" + nested + "](tg://emoji?id=" + matches[1] + ")")
				} else {
					return "", fmt.Errorf("badly formatted anchor tag %q", tagContent)
				}
			default:
				return "", fmt.Errorf("unknown tag %q", tag)
			}

			prev = closingClose + 1
			i = closingClose

		case '\\', '_', '*', '~', '`', '[', ']', '(', ')': // these all need to be escaped to ensure we retain the same message
			out.WriteString(html.UnescapeString(string(in[prev:i])))
			out.WriteRune('\\')
			out.WriteRune(in[i])
			prev = i + 1
		}
	}
	out.WriteString(html.UnescapeString(string(in[prev:])))

	for idx, btn := range buttons {
		bText, err := cv.ButtonToMarkdown(btn)
		if err != nil {
			return "", fmt.Errorf("failed to convert button %d (%s) to markdown: %w", idx, btn.Name, err)
		}
		out.WriteString("\n" + bText)
	}

	return out.String(), nil
}

func (cv ConverterV2) ButtonToMarkdown(btn ButtonV2) (string, error) {
	sameline := ""
	if btn.SameLine {
		sameline = cv.SameLineSuffix
	}

	if prefix, ok := cv.Prefixes[btn.Type]; ok {
		return "[" + EscapeMarkdownV2([]rune(btn.Name)) + "](" + prefix + "//" + html.UnescapeString(btn.Content) + sameline + ")", nil
	}
	return "", ErrNoButtonContent
}
