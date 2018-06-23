package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/urfave/cli"
	"gopkg.in/night-codes/types.v1"
)

func moduleAction(c *cli.Context) error {
	name := stripSlashes(c.String("name"))
	title := c.String("title")
	menu := c.String("menu")

	if name == "" {
		return errors.New("Flag --name is required")
	}
	if title == "" {
		title = strings.Title(name)
	}
	if menu != "MainMenu" && menu != "DropMenu" && menu != "" {
		fmt.Println("[Warning] Menu name is wrong. Module wasn't added to menu!")
		menu = ""
	}
	s := obj{
		"name":       name,
		"Name":       strings.Title(name),
		"Collection": c.String("collection"),
		"Title":      title,
		"Menu":       menu,
		"AddSearch":  c.Bool("add-search"),
		"AddSort":    c.Bool("add-sort"),
		"AddTabs":    c.Bool("add-tabs"),
		"AddFilters": c.Bool("add-filters"),
		"AddPages":   c.Bool("add-pages"),
		"Vendor":     c.Bool("vendor"),
		"GroupTo":    c.String("group"),
		"SubDir":     stripSlashes(c.String("subdir")),
	}

	return modAction("./", s)
}

func modAction(path string, s obj) error {
	var main string

	if mainBytes, err := ioutil.ReadFile(path + "main.go"); err == nil {
		main = string(mainBytes)
	}

	if !strings.Contains(main, "gopkg.in/night-codes/summer.v1") || !strings.Contains(main, "summer.Create(") {
		return errors.New("Current directory is not Summer project")
	}

	o := 1
	startStr := "summer.Create("
	cr := strings.Index(main, startStr) + len(startStr)
	cr0 := cr
	for {
		a := strings.Index(main[cr:], "(")
		b := strings.Index(main[cr:], ")")
		if b != -1 {
			if a < b {
				cr += a + 1
				o++
			} else {
				o--
				cr += b + 1
				if o < 1 {
					break
				}
			}
			continue
		}
		return errors.New("Unsupported main.go")
	}
	sett := main[cr0 : cr-1]

	fmt.Println("SummerGen generates ", s["name"], "module...\n ")
	viewsPath := "templates/main"
	viewsDotPath := "templates/dpT.js"
	vendor := types.Bool(s["Vendor"])
	name := types.String(s["name"])
	subdir := types.String(s["SubDir"])
	app := ""

	fullName := name + strings.Title(subdir)
	s["FullName"] = fullName
	if vendor {
		fmt.Println("[Warning] Vendor mode!")

		cr01 := regexp.MustCompile("import\\s*\\({1}").FindAllStringSubmatchIndex(main, -1)
		if len(cr01) < 1 || len(cr01[0]) < 2 {
			return errors.New("Unsupported main.go")
		}
		cr1 := cr01[0][1]
		app = path + "vendor/" + fullName + "/" + fullName + ".go"
		main = main[:cr1] + "\n\t\"" + fullName + "\"\n" + main[cr1+1:cr] + "\n\t" + fullName + "Module = " + fullName + ".New(panel)" + main[cr:]
		fmt.Println("Writting", path+"main.go....")
		if err := ioutil.WriteFile(path+"main.go", []byte(main), 0755); err != nil {
			return err
		}
	} else {
		app = path + fullName + ".go"
	}

	v1 := regexp.MustCompile("Views\\s*\\:\\s*[\"`]{1}(\\S*)[\"`]{1}").FindAllStringSubmatch(sett, -1)
	if len(v1) > 0 && len(v1[0]) > 1 {
		viewsPath = v1[0][1]
	}

	v2 := regexp.MustCompile("ViewsDoT\\s*\\:\\s*[\"`]{1}(\\S*)[\"`]{1}").FindAllStringSubmatch(sett, -1)
	if len(v2) > 0 && len(v2[0]) > 1 {
		viewsDotPath = v2[0][1]
	}
	if len(subdir) > 0 {
		subdir = "/" + subdir
	}
	fmt.Println("1) Creating", app)
	if err := writeFile(app, moduleTpl, "module.go", s); err != nil {
		return err
	}

	fmt.Println("2) Creating", path+viewsPath+"/"+name+subdir+".html")
	if err := writeFile(path+viewsPath+"/"+name+subdir+".html", moduleTpl, "module.html", s); err != nil {
		return err
	}
	fmt.Println("3) Creating", path+viewsDotPath+"/"+name+subdir+"/*.* files")
	if err := writeFiles(path+viewsDotPath+"/"+name+subdir+"/", []string{"script.js", "item.html", "noitems.html", "form-add.html", "form-edit.html"}, moduleTpl, arr{s, s, s, obj{}, obj{}}); err != nil {
		return err
	}
	fmt.Println("\nModule", name, "successful created!")
	return nil
}
