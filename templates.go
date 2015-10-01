package templates

import (
	"fmt"
	"text/template"
	"path/filepath"
	"github.com/oxtoacart/bpool"
	"net/http"
)


type Templates struct {
	buffer *bpool.BufferPool
	templates map[string]*template.Template
}

func New(dir string) *Templates {

	t := make_templates(dir)


	return &Templates{
		bpool.NewBufferPool(512),
		t,
	}
}

// Load templates on program initialisation
func make_templates(templatesDir string) map[string]*template.Template {
	t := make(map[string]*template.Template)

	layouts, err := filepath.Glob(templatesDir + "layouts/*.html")
	if err != nil {
		panic("Layouts missing")
	}

	includes, err := filepath.Glob(templatesDir + "includes/*.html")
	if err != nil {
		panic("Includes missing")
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, layout := range layouts {
		files := append(includes, layout)
		t[filepath.Base(layout)] = template.Must(template.ParseFiles(files...))
	}

	return t
}


// renderTemplate is a wrapper around template.ExecuteTemplate.
func (t Templates) Render(w http.ResponseWriter, name string, layout bool, data map[string]interface{}) error {
	// Ensure the template exists in the map.
	tmpl, ok := t.templates[name]
	if !ok {
		return fmt.Errorf("The template %s does not exist.", name)
	}

	// Create a buffer to temporarily write to and check if any errors were encounted.
	buf := t.buffer.Get()
	defer t.buffer.Put(buf)

	if layout {
		err := tmpl.ExecuteTemplate(buf, "base", data)
		if err != nil {
			panic("Error generating template: " + name)
		}
	} else {
		err := tmpl.ExecuteTemplate(buf, "content", data)
		if err != nil {
			panic("Error generating template: " + name)
		}
	}
	buf.WriteTo(w)

	return nil
}
