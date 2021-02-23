package web

import (
	"html/template"
	"io"
	"os"
	"path"

	rice "github.com/GeertJohan/go.rice"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// echoTempl implements echo.Renderer interface.
type echoTempl struct {
	templates *template.Template
}

// Render implements echo.Renderer interface.
func (t *echoTempl) Render(
	w io.Writer,
	name string,
	data interface{},
	c echo.Context,
) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplate() (*echoTempl, error) {
	box, err := rice.FindBox("./templates")
	if err != nil {
		return nil, errors.Wrap(err, "rice.FindBox")
	}
	t, err := parseFiles(box, nil)
	if err != nil {
		return nil, errors.Wrap(err, "parseFile")
	}
	return &echoTempl{t}, nil
}

func parseFiles(box *rice.Box, t *template.Template) (*template.Template, error) {
	filenames := []string{}
	err := box.Walk(".", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		filenames = append(filenames, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(filenames) == 0 {
		return nil, errors.New("no files found")
	}
	for _, filename := range filenames {
		name := path.Base(filename)
		s, err := box.String(name)
		if err != nil {
			return nil, errors.Wrap(err, "box.String")
		}
		var tmpl *template.Template
		if t == nil {
			t = template.New(name)
		}
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}
		_, err = tmpl.Parse(s)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
