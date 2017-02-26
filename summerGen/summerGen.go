package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "summerGen"
	app.Author = "Oleksiy Chechel (alex.mirrr@gmail.com)"
	app.Version = "1.6.15"
	app.Commands = []cli.Command{
		cli.Command{
			Name:   "module",
			Usage:  "Generate Summer module",
			Flags:  moduleFlags,
			Action: moduleAction,
		},
		cli.Command{
			Name:   "project",
			Usage:  "Generate Summer project structure",
			Flags:  projectFlags,
			Action: projectAction,
		},
	}
	app.Usage = "Summer projects and modules generator \n\t More information: \n\n\t\t\tsummerGen module --help \n\t\t\tsummerGen project --help\n\t"
	app.Run(os.Args)
}
