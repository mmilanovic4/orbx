package formatutil

import "testing"

func TestIndentJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"keeps order, precision and html characters",
			`{"b":1,"a":"<x>&y","id":12345678901234567890}`,
			"{\n  \"b\": 1,\n  \"a\": \"<x>&y\",\n  \"id\": 12345678901234567890\n}",
		},
		{"trailing newline", "[1,2]\n", "[\n  1,\n  2\n]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IndentJSON([]byte(tt.input))
			if err != nil {
				t.Fatalf("IndentJSON() error = %v", err)
			}
			if got != tt.expected {
				t.Errorf("IndentJSON() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestIndentJSONInvalid(t *testing.T) {
	for _, input := range []string{`{"a":}`, `{} trailing`, ``} {
		if _, err := IndentJSON([]byte(input)); err == nil {
			t.Errorf("IndentJSON(%q) expected error, got nil", input)
		}
	}
}

func TestIndentXML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"nesting, text and empty elements",
			`<a><b>1</b><c></c><d x="1"/></a>`,
			"<a>\n  <b>1</b>\n  <c/>\n  <d x=\"1\"/>\n</a>",
		},
		{
			"namespace prefixes are kept",
			`<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Body><x>1</x></soap:Body></soap:Envelope>`,
			"<soap:Envelope xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\">\n  <soap:Body>\n    <x>1</x>\n  </soap:Body>\n</soap:Envelope>",
		},
		{
			"already indented input",
			"<a>\n    <b>1</b>\n\n</a>\n",
			"<a>\n  <b>1</b>\n</a>",
		},
		{
			"declaration and comment",
			`<?xml version="1.0"?><!-- note --><a/>`,
			"<?xml version=\"1.0\"?>\n<!-- note -->\n<a/>",
		},
		{
			"escaping",
			`<a t="x &amp; &quot;y&quot;">1 &lt; 2 &amp; 3</a>`,
			`<a t="x &amp; &quot;y&quot;">1 &lt; 2 &amp; 3</a>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IndentXML([]byte(tt.input))
			if err != nil {
				t.Fatalf("IndentXML() error = %v", err)
			}
			if got != tt.expected {
				t.Errorf("IndentXML() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestIndentXMLInvalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"mismatched tags", "<a><b>x</a>"},
		{"unclosed tag", "<a><b>x</b>"},
		{"broken attribute", "<a x=1></a>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := IndentXML([]byte(tt.input)); err == nil {
				t.Errorf("IndentXML(%q) expected error, got nil", tt.input)
			}
		})
	}
}
