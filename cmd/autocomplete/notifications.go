package autocomplete

import (
	"fmt"

	"github.com/Scalingo/cli/config"
	"github.com/urfave/cli"
	"gopkg.in/errgo.v1"
)

func NotificationsRemoveAutoComplete(c *cli.Context) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	client, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	resources, err := client.NotificationsList(appName)
	if err == nil {
		for _, resource := range resources {
			fmt.Println(resource.ID)
		}
	}

	return nil
}
