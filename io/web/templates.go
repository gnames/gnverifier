package web

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path"
	"strings"

	rice "github.com/GeertJohan/go.rice"
	vlib "github.com/gnames/gnlib/ent/verifier"
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
		addFuncs(tmpl)
		_, err = tmpl.Parse(s)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

func addFuncs(tmpl *template.Template) {
	tmpl.Funcs(template.FuncMap{
		"isEven": func(i int) bool {
			return i%2 == 0
		},
		"classification": func(pathStr, rankStr string) string {
			if pathStr == "" {
				return ""
			}
			paths := strings.Split(pathStr, "|")
			var ranks []string
			if rankStr != "" {
				ranks = strings.Split(rankStr, "|")
			}

			res := make([]string, len(paths))
			for i := range paths {
				path := strings.TrimSpace(paths[i])
				if len(ranks) == len(paths) {
					rank := strings.TrimSpace(ranks[i])
					if rank != "" {
						path = fmt.Sprintf("%s (%s)", path, rank)
					}
				}
				res[i] = path
			}
			return strings.Join(res, " >> ")
		},
		"matchType": func(mt vlib.MatchTypeValue, ed int) template.HTML {
			var res string
			clr := map[string]string{
				"green":  "#080",
				"yellow": "#a80",
				"red":    "#800",
			}
			switch mt {
			case vlib.Exact:
				res = fmt.Sprintf("<span style='color: %s'>%s match by canonical form</span>", clr["green"], mt)
			case vlib.NoMatch:
				res = fmt.Sprintf("<span style='color: %s'>%s</span>", clr["red"], mt)
			case vlib.Fuzzy, vlib.PartialFuzzy:
				res = fmt.Sprintf("<span style='color: %s'>%s match, edit distance: %d</span>", clr["yellow"], mt, ed)
			default:
				res = fmt.Sprintf("<span style='color: %s'>%s match</span>", clr["yellow"], mt)
			}
			return template.HTML(res)
		},
	})
}
