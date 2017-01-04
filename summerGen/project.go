package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli"
)

func projectAction(c *cli.Context) error {
	name := stripSlashes(c.String("name"))
	if name == "" {
		return errors.New("Flag --name is required")
	}
	if err := os.Mkdir(name, 0755); err != nil {
		return err
	}
	if err := os.Mkdir(name+"/files", 0755); err != nil {
		return err
	}

	if c.Bool("vendor") {
		if err := writeFile(name+"/vendor/hello/hello.go", helloTpl, "hello.go", obj{"Vendor": true}); err != nil {
			return err
		}
	} else {
		if err := writeFile(name+"/hello.go", helloTpl, "hello.go", obj{}); err != nil {
			return err
		}
	}

	views := name + "/" + stripSlashes(c.String("views"))
	viewsDotHello := name + "/" + stripSlashes(c.String("views-dot")) + "/hello"

	if err := writeFiles(views+"/", []string{"howto.html", "hello.html"}, helloTpl, arr{obj{}, obj{}}); err != nil {
		return err
	}
	if err := writeFiles(viewsDotHello+"/", []string{"icons.html", "icoinfo.html", "script.js"}, helloTpl, arr{obj{}, obj{"itclass": "{{=it.class}}"}, obj{}}); err != nil {
		return err
	}

	dbName := c.String("db")
	if len(dbName) == 0 {
		dbName = name
	}
	Title := c.String("title")
	if len(Title) == 0 {
		Title = strings.Title(name)
	}

	if err := writeFile(name+"/main.go", mainTpl, "main.go", obj{
		"Title":    Title,
		"Vendor":   c.Bool("vendor"),
		"Port":     c.Int("port"),
		"Path":     c.String("dir"),
		"DBName":   dbName,
		"Views":    stripSlashes(c.String("views")),
		"ViewsDoT": stripSlashes(c.String("views-dot")),
	}); err != nil {
		return err
	}

	if c.Bool("demo") {
		fmt.Println("Install demo pages...")
		if err := modAction("./"+name+"/", obj{
			"name":       "news",
			"Name":       "News",
			"Collection": "news",
			"Title":      "Daily news",
			"Menu":       "MainMenu",
			"AddSearch":  true,
			"AddSort":    true,
			"AddTabs":    false,
			"AddFilters": false,
			"AddPages":   true,
			"Vendor":     c.Bool("vendor"),
			"GroupTo":    "",
			"SubDir":     "",
		}); err != nil {
			return err
		}
		if err := modAction("./"+name+"/", obj{
			"name":       "users",
			"Name":       "Users",
			"Collection": "Users",
			"Title":      "Users of site",
			"Menu":       "MainMenu",
			"AddSearch":  true,
			"AddTabs":    true,
			"AddFilters": false,
			"AddPages":   true,
			"Vendor":     c.Bool("vendor"),
			"GroupTo":    "",
			"SubDir":     "",
		}); err != nil {
			return err
		}
	}

	fmt.Println("Project", name, "successful created!")
	return nil
}
