package main

import (
	"errors"
	"fmt"
	"os"
	"text/template"

	"github.com/urfave/cli"
)

func projectAction(c *cli.Context) error {
	name := c.String("name")
	if name == "" {
		return errors.New("Flag --name is required")
	}
	if err := os.Mkdir(name, 0755); err != nil {
		return err
	}

	if c.Bool("vendor") {
		if err := os.MkdirAll(name+"/vendor/hello", 0755); err != nil {
			return err
		}
		if err := write(name+"/vendor/hello/hello.go", helloTpl, "hello.go", obj{"Vendor": true}); err != nil {
			return err
		}
	} else {
		if err := write(name+"/hello.go", helloTpl, "hello.go", obj{}); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(name+"/"+c.String("views"), 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(name+"/"+c.String("views-dot")+"/hello", 0755); err != nil {
		return err
	}

	if err := write(name+"/"+c.String("views")+"/howto.html", helloTpl, "howto.html", obj{}); err != nil {
		return err
	}
	if err := write(name+"/"+c.String("views")+"/hello.html", helloTpl, "hello.html", obj{}); err != nil {
		return err
	}

	if err := write(name+"/"+c.String("views-dot")+"/hello/icons.html", helloTpl, "icons.html", obj{}); err != nil {
		return err
	}
	if err := write(name+"/"+c.String("views-dot")+"/hello/icoinfo.html", helloTpl, "icoinfo.html", obj{"itclass": "{{=it.class}}"}); err != nil {
		return err
	}
	if err := write(name+"/"+c.String("views-dot")+"/hello/script.js", helloTpl, "script.js", obj{}); err != nil {
		return err
	}

	if err := write(name+"/main.go", mainTpl, "main.go", obj{
		"Demo":     c.Bool("demo"),
		"Title":    c.String("title"),
		"Vendor":   c.Bool("vendor"),
		"Port":     c.Int("port"),
		"Path":     c.String("dir"),
		"DBName":   c.String("db"),
		"Views":    c.String("views"),
		"ViewsDoT": c.String("views-dot"),
	}); err != nil {
		return err
	}

	fmt.Println("Project", name, "successful created!")
	return nil
}

func write(filename string, t *template.Template, tname string, data interface{}) error {
	fo, err := os.Create(filename)
	if err != nil {
		return err
	}
	if err := t.ExecuteTemplate(fo, tname, data); err != nil {
		return err
	}
	fo.Close()
	return nil
}
