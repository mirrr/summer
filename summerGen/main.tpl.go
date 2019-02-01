package main

import (
	"text/template"
)

var mainTpl = template.Must(template.New("main.go").Parse(`package main

import (
	{{if .Vendor}}"hello"{{end}}
	"fmt"
	"gopkg.in/night-codes/summer.v2"
)

type (
	obj map[string]interface{}
	arr []interface{}
)

var (
	panel = summer.Create(summer.Settings{
		Title:       "{{.Title}}",
		Port:        {{.Port}},
		DefaultPage: "hello",
		Path:        "{{.Path}}", // application path
		DBName:      "{{.DBName}}",
		Views:       "{{.Views}}",
		ViewsDoT:    "{{.ViewsDoT}}", // doT.js templates
		FirstStart: func() { // some DB migrations etc.
			fmt.Println("Application is running for the first time!")
		},
		Debug: summer.Env("production", "") == "", // set env. var "production" for Debug:false
		JS:    []string{},                         // add custom JS files to template
		CSS:   []string{},                         // add custom CSS files to template
	}){{if .Vendor}}
	helloModule  = hello.New(panel){{end}}
)

func main() {
	fmt.Println("Application started at http://localhost:{{.Port}}/{{.Path}}")
	summer.Wait()
}
`))
