package tree_test

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"testing"

	. "ecosia/tree"
)

func TestTemplateView(t *testing.T) {
	type args struct {
		name     string
		input    string
		expected string
		tpl      string
	}

	tests := []args{
		{
			name:     "simple",
			input:    "baobab",
			expected: "baobab",
			tpl:      "{{.}}",
		},
		{
			name:     "utf8",
			input:    "Rêveberiya Xweser a Bakur û Rojhilatê Sûriyeyê",
			expected: "Rêveberiya Xweser a Bakur û Rojhilatê Sûriyeyê",
			tpl:      "{{.}}",
		},
		{
			name:     "xss",
			input:    `Tree of <style>.xss{background-image:url("javascript:alert('sup')");}</style><a class=xss></a> XSS`,
			expected: `Tree of &lt;style&gt;.xss{background-image:url(&#34;javascript:alert(&#39;sup&#39;)&#34;);}&lt;/style&gt;&lt;a class=xss&gt;&lt;/a&gt; XSS`,
			tpl:      "{{.}}",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tpl, _ := template.New("test").Parse(test.tpl)

			view := NewTemplateView(tpl)

			var w bytes.Buffer

			if err := view(&w, test.input); err != nil {
				t.Fatalf("view failed: %s", err)
			}

			if w.String() != test.expected {
				t.Fatalf("expected '%s', got '%s'", test.expected, w.String())
			}
		})
	}

	t.Run("view writer error", func(t *testing.T) {
		tpl, _ := template.New("test").Parse("{{.}}")
		view := NewTemplateView(tpl)
		err := view(badWriter{}, "oops")
		if err == nil {
			t.Fatal("expected an error here")
		}
	})

	t.Run("bad default tpl", func(t *testing.T) {
		aux := DefaultTpl
		defer func() {
			DefaultTpl = aux
		}()
		DefaultTpl = "{{missingFunction .}}"

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("should have panicked vor an invalid default template")
			}
		}()

		view := NewTemplateView(nil)
		err := view(badWriter{}, "oops")
		if err == nil {
			t.Fatal("expected an error here")
		}
	})
}

type badWriter struct{}

func (b badWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("this is only meant to fail")
}

func ExampleNewTemplateView() {
	tpl, _ := template.New("test").Parse("<p>{{.}}</p>")

	// this will render the template by inserting the tree into it,
	// then it will send it to the writer
	view := NewTemplateView(tpl)

	var buf bytes.Buffer

	if err := view(&buf, "baobab"); err != nil {
		panic(err)
	}

	fmt.Println(buf.String())
	// Output: <p>baobab</p>
}

func BenchmarkTemplateView(b *testing.B) {
	view := NewTemplateView(nil)

	var buf bytes.Buffer

	for i := 0; i < b.N; i++ {
		_ = view(&buf, "baobab")
	}
}
