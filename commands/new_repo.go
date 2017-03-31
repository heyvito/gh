package commands

import (
	"fmt"

	"github.com/octokit/go-octokit/octokit"
	"github.com/urfave/cli"
	"github.com/victorgama/gh/utils"
)

var NewRepo = cli.Command{
	Name:  "new",
	Usage: "Creates a new repository",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "private",
			Usage: "creates a private repository",
		},
		cli.BoolFlag{
			Name:  "org",
			Usage: "Indicates that the new repository will be created in a organization",
		},
	},
	Action: func(c *cli.Context) error {
		if len(c.Args()) != 1 {
			return fmt.Errorf("the 'new' command requires a repository name argument")
		}
		r := utils.RepoURLFromString(c.Args()[0])
		client := utils.NewClient()
		r.ToURL()
		private := c.Bool("private")
		if newName, changed := utils.NormalizeRepoName(r.RepoName); changed {
			fmt.Printf("Your repository will be created as %s/%s\n", r.Username, newName)
			if !utils.Confirm("Seems okay? [y]/n", true) {
				return fmt.Errorf("aborted")
			}
		}
		uri := octokit.UserRepositoriesURL
		params := octokit.M{}
		if utils.UserIsOrg(r.Username) {
			uri = octokit.OrganizationReposURL
			params = octokit.M{"org": r.Username}
		}

		_, res := client.Repositories().Create(&uri, params, octokit.Repository{
			Name:    r.RepoName,
			Private: private,
		})
		if res.HasError() {
			return res.Err
		}
		return nil
	},
}
