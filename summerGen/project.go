package main

import (
	"errors"
	"fmt"
	"os"

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
		if err := os.Mkdir(name+"/vendor", 0755); err != nil {
			return err
		}
	}
	if err := os.MkdirAll(name+"/"+c.String("views"), 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(name+"/"+c.String("views-dot"), 0755); err != nil {
		return err
	}

	fo, err := os.Create(name + "/main.go")
	if err != nil {
		panic(err)
	}

	if err := mainTpl.ExecuteTemplate(fo, "main.go", map[string]interface{}{
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
	fo.Close()
	fmt.Println("Project", name, "successful created!")

	return nil
}
