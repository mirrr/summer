package ttpl

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"net/http"
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
func Use(r *gin.Engine, patterns []string, funcMap ...template.FuncMap) {
	t := &template.Template{}
	if len(funcMap) > 0 {
		t = template.New("").Funcs(funcMap[0])
	} else {
		t = template.New("")
	}
	for _, pattern := range patterns {
		t = template.Must(t.ParseGlob(pattern))
	}

	r.HTMLRender = PageTemplate{"/", t}
}
