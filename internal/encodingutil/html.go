package encodingutil

import "html"

func EncodeHTML(input string) string {
	return html.EscapeString(input)
}

func DecodeHTML(input string) string {
	return html.UnescapeString(input)
}
