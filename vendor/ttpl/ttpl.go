package ttpl

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"
)

// PageTemplate struct for gin
type PageTemplate struct {
	TemplatePath string
	templates    *template.Template
}

// PageRender struct for gin
type PageRender struct {
	Template *template.Template
	Data     interface{}
	Name     string
}

func (r PageTemplate) Instance(name string, data interface{}) render.Render {
	return PageRender{
		Template: r.templates,
		Name:     name,
		Data:     data,
	}
}

func (r PageRender) Render(w http.ResponseWriter) error {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{"text/html; charset=utf-8"}
	}

	if len(r.Name) > 0 {
		if err := r.Template.ExecuteTemplate(w, r.Name, r.Data); err != nil {
			fmt.Println("Template err: ", err.Error())
		}
	} else {
		if err := r.Template.Execute(w, r.Data); err != nil {
			fmt.Println("Template err: ", err.Error())
		}
	}

	return nil
}

// Use ttpl render
func Use(r *gin.Engine, patterns []string, dotPath string, funcMap ...template.FuncMap) {
	t := template.New("")
	if len(funcMap) > 0 {
		t = t.Funcs(funcMap[0])
	}
	for _, pattern := range patterns {
		filenames, err := filepath.Glob(pattern)
		if len(filenames) > 0 && err == nil {
			tt, err := parseFiles(t, dotPath, filenames...)
			if err == nil {
				t = tt
			} else if err != nil {
				fmt.Println(err)
			}
		} else if err != nil {
			fmt.Println(err)
		}
	}

	r.HTMLRender = PageTemplate{"/", t}
}

// parseFiles is the helper for the method and function. If the argument
// template is nil, it is created from the first file.
func parseFiles(t *template.Template, dotPath string, filenames ...string) (*template.Template, error) {
	for _, filename := range filenames {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		s := string(b)
		name := filepath.Base(filename)
		shortName := strings.Split(name, ".")[0]

		// DoT.js
		dots, err := filepath.Glob(dotPath + "/" + shortName + "/*")
		if len(dots) > 0 {
			for _, dot := range dots {
				s = s + `{{ dot "` + shortName + `/` + filepath.Base(dot) + `" }}` + "\n"
			}
		}

		if name != "layout.html" && name != "login.html" && name != "firstStart.html" {
			s = `{{template "header" .}}` + s + `{{template "footer" .}}`
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
