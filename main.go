package main

import (
	"os"

	"github.com/urfave/cli"
	"github.com/victorgama/gh/commands"
)

func main() {
	app := cli.NewApp()
	app.Name = "gh"
	app.Version = "0.1.0"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Victor Gama",
			Email: "hey@vito.io",
		},
	}
	app.HelpName = "gh"
	app.Usage = "Controls GitHub from the command line"
	app.UsageText = "gh [command]"
	app.Commands = []cli.Command{
		commands.NewRepo,
		commands.RmRepo,
		commands.RepoList,
		commands.Collab,
		commands.Teams,
		commands.Open,
	}
	app.Run(os.Args)
}
