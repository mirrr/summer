package ttpl

import (
	"fmt"
	"github.com/gin-gonic/gin/render"
	"gopkg.in/gin-gonic/gin.v1"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"text/template"
	"time"
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

var (
	dmu      sync.Mutex
	dots     = map[string]string{}
	spliter  = regexp.MustCompile("[\\s\\/]+")
	siteVars = map[string]bool{
		"action": false,
		"title":  false,
		"login":  false,
		"module": false,
		"path":   false,
		"ajax":   false,
		"socket": false,
		"allow":  false,
		"js":     true,
		"css":    true,
	}
)

func init() {
	if os.Getenv("ENV") == "development" {
		go func() {
			for range time.Tick(time.Second) {
				func() {
					dmu.Lock()
					defer dmu.Unlock()
					dots = map[string]string{}
				}()
			}
		}()
	}
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

	site := map[string]interface{}{}
	for key, array := range siteVars {
		Key := strings.Title(key)
		if array {
			site[key] = header[Key]
		} else {
			site[key] = ""
			if val := header[Key]; len(val) > 0 {
				site[key] = val[0]
			}
		}
		delete(header, Key)
	}
	r.Template.Funcs(template.FuncMap{"site": func(name string) interface{} { return site[name] }})

	name := r.Name
	if name != "login.html" && name != "firstStart.html" && site["allow"] != "true" {
		name = "access-deny"
	}
	if len(name) > 0 {
		if err := r.Template.ExecuteTemplate(w, name, r.Data); err != nil {
			fmt.Println("Template err: ", err.Error())
		}
	} else {
		if err := r.Template.Execute(w, r.Data); err != nil {
			fmt.Println("Template err: ", err.Error())
		}
	}

	return nil
}

func dot(dotPath string) func(name string) string {
	return func(name string) string {
		dmu.Lock()
		defer dmu.Unlock()
		fullname := dotPath + "/" + name
		if _, exists := dots[fullname]; !exists {
			dots[fullname] = "<!-- Template '" + name + "' not found! -->\n"
			if dat, err := ioutil.ReadFile(fullname); err == nil {
				s := strings.Split(name, ".")
				tplName := spliter.Split(s[0], -1)
				if s[len(s)-1] == "js" { // js темплейты
					dots[fullname] = "<!-- doT.js template - " + name + " -->\n" +
						"<script type='text/javascript' id='tpl_" + strings.Join(tplName[1:], "-") + "'>\n" + string(dat) + "</script>\n"

				} else { // html темплейты
					dots[fullname] = "<!-- doT.js template - " + name + " -->\n" +
						"<script type='text/html' id='tpl_" + strings.Join(tplName[1:], "-") + "'>\n" + string(dat) + "</script>\n"
				}
			}
		}
		return dots[fullname]
	}
}

// Use ttpl render
func Use(r *gin.Engine, pathes []string, dotPath string, funcMap ...template.FuncMap) {
	t := template.New("")
	if len(funcMap) == 0 {
		funcMap = []template.FuncMap{}
	}

	funcMap[0]["dot"] = dot(dotPath)
	t = t.Funcs(funcMap[0])

	for _, path := range pathes {
		if _, err := os.Stat(path); err == nil {
			filepath.Walk(path, func(file string, info os.FileInfo, err error) error {
				if !info.IsDir() && err == nil {
					tt, err := parseFiles(t, dotPath, path, file)
					if err == nil {
						t = tt
					} else if err != nil {
						fmt.Println(err)
					}
				}
				return nil
			})
		}
	}

	r.HTMLRender = PageTemplate{"/", t}
}

// parseFiles is the helper for the method and function. If the argument
// template is nil, it is created from the first file.
func parseFiles(t *template.Template, dotPath string, path string, filenames ...string) (*template.Template, error) {
	for _, filename := range filenames {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		s := string(b)
		name := strings.Replace(filename, path, "", 1)
		shortName := strings.Split(name, ".")[0]

		// DoT.js
		dots, err := filepath.Glob(dotPath + "/" + shortName + "/*")
		if len(dots) > 0 {
			for _, dot := range dots {
				base := filepath.Base(dot)
				if !strings.Contains(base, ".") {
					subs, _ := filepath.Glob(dotPath + "/" + shortName + "/" + base + "/*")
					for _, sub := range subs {
						if !strings.Contains(filepath.Base(sub), " ") {
							s = s + `{{ dot "` + shortName + `/` + base + `/` + filepath.Base(sub) + `"}}` + "\n"
						}
					}
				} else if !strings.Contains(base, " ") {
					s = s + `{{ dot "` + shortName + `/` + base + `" }}` + "\n"
				}
			}
		}

		sOrigin := s
		header := `{{template "header" .}}`
		footer := `{{template "footer" .}}`
		if strings.Contains(s, "SUMMER-NO-HEADER") {
			header = ""
		}
		if strings.Contains(s, "SUMMER-NO-FOOTER") {
			footer = ""
		}
		if name != "layout.html" && name != "login.html" && name != "firstStart.html" {
			s = header + s + footer
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

		name = "summer-origin-" + name
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}
		_, err = tmpl.Parse(sOrigin)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
