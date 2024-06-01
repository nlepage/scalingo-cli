package cmd

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/session"
	"github.com/Scalingo/cli/ui"
)

var (
	uiCommand = cli.Command{
		Name:        "ui",
		Category:    "Global",
		Description: "Start interactive Terminal User Interface",
		Usage:       "Start TUI",
		Before: func(c *cli.Context) error {
			token := os.Getenv("SCALINGO_API_TOKEN")

			currentUser, err := config.C.CurrentUser(c.Context)
			if err != nil || currentUser == nil {
				err := session.Login(c.Context, session.LoginOpts{APIToken: token})
				if err != nil {
					errorQuit(c.Context, err)
				}
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			if err := ui.Start(c.Context); err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
	}
)
