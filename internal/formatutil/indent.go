package formatutil

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"strings"
)

// IndentJSON pretty prints JSON exactly as written: key order, number
// precision and characters like < > & are kept, which a decode and
// re-encode round trip would change.
func IndentJSON(data []byte) (string, error) {
	var buf bytes.Buffer
	if err := json.Indent(&buf, bytes.TrimSpace(data), "", "  "); err != nil {
		return "", err
	}
	return buf.String(), nil
}

var (
	xmlTextEscaper = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;")
	xmlAttrEscaper = strings.NewReplacer("&", "&amp;", "<", "&lt;", `"`, "&quot;", "\t", "&#x9;", "\n", "&#xA;", "\r", "&#xD;")
)

// IndentXML pretty prints an XML document. Whitespace between tags is
// replaced by indentation, elements holding only text stay on one line and
// namespace prefixes are kept as written. xml.Encoder is not used because
// it rewrites prefixed names into generated namespace declarations.
func IndentXML(data []byte) (string, error) {
	// RawToken keeps prefixes intact but does not check that tags match,
	// so the document is validated with Token first
	dec := xml.NewDecoder(bytes.NewReader(data))
	for {
		_, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
	}

	var tokens []xml.Token
	dec = xml.NewDecoder(bytes.NewReader(data))
	for {
		tok, err := dec.RawToken()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		if text, ok := tok.(xml.CharData); ok && len(bytes.TrimSpace(text)) == 0 {
			continue
		}
		tokens = append(tokens, xml.CopyToken(tok))
	}

	var b strings.Builder
	depth := 0
	newline := func() {
		if b.Len() > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(strings.Repeat("  ", depth))
	}

	for i := 0; i < len(tokens); i++ {
		switch t := tokens[i].(type) {
		case xml.StartElement:
			newline()
			name := xmlName(t.Name)
			b.WriteString("<" + name)
			for _, attr := range t.Attr {
				b.WriteString(" " + xmlName(attr.Name) + `="` + xmlAttrEscaper.Replace(attr.Value) + `"`)
			}

			// <a></a> becomes <a/> and <a>text</a> stays on one line
			if _, ok := tokenAt(tokens, i+1).(xml.EndElement); ok {
				b.WriteString("/>")
				i++
				continue
			}
			if text, ok := tokenAt(tokens, i+1).(xml.CharData); ok {
				if _, ok := tokenAt(tokens, i+2).(xml.EndElement); ok {
					b.WriteString(">" + xmlTextEscaper.Replace(string(text)) + "</" + name + ">")
					i += 2
					continue
				}
			}

			b.WriteString(">")
			depth++
		case xml.EndElement:
			depth--
			newline()
			b.WriteString("</" + xmlName(t.Name) + ">")
		case xml.CharData:
			newline()
			b.WriteString(xmlTextEscaper.Replace(strings.TrimSpace(string(t))))
		case xml.Comment:
			newline()
			b.WriteString("<!--" + string(t) + "-->")
		case xml.ProcInst:
			newline()
			b.WriteString("<?" + t.Target)
			if len(t.Inst) > 0 {
				b.WriteString(" " + string(t.Inst))
			}
			b.WriteString("?>")
		case xml.Directive:
			newline()
			b.WriteString("<!" + string(t) + ">")
		}
	}

	return b.String(), nil
}

// xmlName formats a name read by RawToken, where Space holds the prefix.
func xmlName(n xml.Name) string {
	if n.Space == "" {
		return n.Local
	}
	return n.Space + ":" + n.Local
}

func tokenAt(tokens []xml.Token, i int) xml.Token {
	if i < len(tokens) {
		return tokens[i]
	}
	return nil
}
