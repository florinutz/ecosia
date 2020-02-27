package tree

import (
	"fmt"
	"html/template"
	"io"
)

var DefaultTpl = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Tree greeter app</title>
</head>
<body>
<p>
	{{- if . -}}
		It's nice to know that your favorite tree is a {{.}}
	{{- else -}}
		Please tell me your favorite tree
	{{- end -}}
</p>
</body>
</html>`

// ViewFunc decorates data (tree) and writes the result to a writer
type ViewFunc func(w io.Writer, tree string) error

// NewTemplateView generates a ViewFunc that injects data into a html template
func NewTemplateView(tpl *template.Template) ViewFunc {
	if tpl == nil {
		// use default tpl
		var err error
		tpl, err = template.New("tree").Parse(DefaultTpl)
		if err != nil {
			panic(fmt.Sprintf("default template is broken: %s", err))
		}
	}

	return func(w io.Writer, tree string) error {
		if err := tpl.Execute(w, tree); err != nil {
			return fmt.Errorf("couldn't send rendered tpl: %w", err)
		}

		return nil
	}
}
