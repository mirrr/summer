package ttpl

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

type (
	// PageTemplate struct for gin
	PageTemplate struct {
		TemplatePath string
		templates    *template.Template
		funcMap      template.FuncMap
	}

	// PageRender struct for gin
	PageRender struct {
		Template *template.Template
		funcMap  template.FuncMap
		Data     interface{}
		Name     string
	}
)

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
		"lang":   false,
		"par_0":  false,
		"par_1":  false,
		"par_2":  false,
		"par_3":  false,
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
		funcMap:  r.funcMap,
		Name:     name,
		Data:     data,
	}
}

func (r PageRender) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{"text/html; charset=utf-8"}
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
	r.Template.Funcs(template.FuncMap{
		"site": func(name string) interface{} { return site[name] },
		"o": func(key string) interface{} {
			if lang, ok := site["lang"]; ok {
				if translate, ok1 := r.funcMap["translate"]; ok1 {
					if translate2, ok2 := translate.(func(module, key, lang string) string); ok2 {
						return translate2(site["module"].(string), key, lang.(string))
					}
				}
			}
			return ""
		},
	})

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

	r.HTMLRender = PageTemplate{"/", t, funcMap[0]}
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
