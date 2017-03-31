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
			Name:  `Victor "Vito" Gama`,
			Email: "hey@vito.io",
		},
	}
	app.HelpName = "gh"
	app.Usage = "Controls GitHub account through the command line"
	app.Commands = []cli.Command{
		commands.NewRepo,
	}
	// fmt.Println(utils.NormalizeRepoName("$hello!"))
	// fmt.Println(utils.NormalizeRepoName("$hello!.git"))
	// r := regexp.MustCompile(`(?:([^\/]+)/?)(.*)`)
	// fmt.Println(r.FindStringSubmatch("d3estudio/bloomzinho"))
	app.Run(os.Args)
}
