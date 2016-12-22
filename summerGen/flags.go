package main

import (
	"github.com/urfave/cli"
)

var (
	moduleFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "Module name",
		},
		cli.StringFlag{
			Name:  "title",
			Usage: "Module title",
		},
		cli.StringFlag{
			Name:  "menu",
			Usage: "Add module to menu (mainMenu, dropMenu)",
		},
		cli.BoolFlag{
			Name:  "add-search",
			Usage: "Add search mechanism",
		},
		cli.BoolFlag{
			Name:  "add-tabs",
			Usage: "Add ajax tabs block",
		},
		cli.BoolFlag{
			Name:  "add-filters",
			Usage: "Add filters block",
		},
		cli.BoolFlag{
			Name:  "vendor",
			Usage: "Use vendor path",
		},
		cli.BoolFlag{
			Name:  "separate",
			Usage: "Separate models form controllers",
		},
	}

	projectFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "Project name",
		},
		cli.StringFlag{
			Name:  "title",
			Value: "Summer App",
			Usage: "Project title",
		},
		cli.StringFlag{
			Name:  "dir",
			Usage: "Panel path on the site (http://localhost/{dir})",
		},
		cli.StringFlag{
			Name:  "db",
			Value: "summer-app",
			Usage: "MongoDB database name",
		},
		cli.IntFlag{
			Name:  "port",
			Value: 8080,
			Usage: "Application port",
		},
		cli.StringFlag{
			Name:  "views",
			Value: "templates/main",
			Usage: "Project templates path",
		},
		cli.StringFlag{
			Name:  "views-dot",
			Value: "templates/dot",
			Usage: "Project doT.js templates path",
		},
		cli.BoolFlag{
			Name:  "vendor",
			Usage: "Use vendor path for modules",
		},
		cli.BoolFlag{
			Name:  "demo",
			Usage: "Add demo modules to project (admins and news)",
		},
	}
)
