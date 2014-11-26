package cmd

import (
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/apps"
	"github.com/codegangsta/cli"
)

// import (
// 	"github.com/Scalingo/cli/appdetect"
// 	"github.com/Scalingo/cli/apps"
// 	"github.com/Scalingo/cli/auth"
// 	"github.com/codegangsta/cli"
// )

var (
	ScaleCommand = cli.Command{
		Name:      "scale",
		ShortName: "s",
		Usage:     "Scale your application instantly",
		Description: `Scale your application processes.
   Example
     'scalingo --app my-app scale web:2 worker:1'
     'scalingo --app my-app scale web:1 worker:0'`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c.GlobalString("app"))
			if len(c.Args()) == 0 {
				cli.ShowCommandHelp(c, "scale")
			} else if err := apps.Scale(currentApp, c.Args()); err != nil {
				errorQuit(err)
			}
		},
	}
)
