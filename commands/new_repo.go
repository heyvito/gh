package commands

import (
	"fmt"
	"os"

	"github.com/octokit/go-octokit/octokit"
	"github.com/urfave/cli"
	"github.com/victorgama/gh/utils"
)

var newRepoLogger = utils.Logger.WithExtra("new")

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
		if len(c.Args()) < 1 {
			return fmt.Errorf("the 'new' command requires at least one repository name")
		}
		client := utils.NewClient()
		private := c.Bool("private")

		repos := []utils.RepoURL{}
		usersOrgs := map[string]bool{}

		// This first run just ensures all repositories are valid ones.
		for _, re := range c.Args() {
			r := utils.RepoURLFromString(re)
			r.ToURL()
			if newName, changed := utils.NormalizeRepoName(r.RepoName); changed {
				fmt.Printf("Your repository will be created as %s/%s\n", r.Username, newName)
				if !utils.Confirm("Seems okay? [y]/n", true) {
					return fmt.Errorf("aborted")
				}
				r.RepoName = newName
			}
			repos = append(repos, r)
		}

		for _, re := range repos {
			newRepoLogger.Timing("Creating %s", re.ToURL())
			uri := octokit.UserRepositoriesURL
			params := octokit.M{}

			userIsOrg := false

			if isOrg, present := usersOrgs[re.Username]; !present {
				usersOrgs[re.Username] = utils.UserIsOrg(re.Username)
				userIsOrg = usersOrgs[re.Username]
			} else {
				userIsOrg = isOrg
			}

			if userIsOrg {
				uri = octokit.OrganizationReposURL
				params = octokit.M{"org": re.Username}
			}

			repo, res := client.Repositories().Create(&uri, params, octokit.Repository{
				Name:    re.RepoName,
				Private: private,
			})
			if res.HasError() {
				if err, ok := res.Err.(*octokit.ResponseError); ok {
					newRepoLogger.Error("%s", utils.FormatError(err))
				}
				os.Exit(1)
			}
			newRepoLogger.Success("Created: %s", repo.FullName)
		}
		return nil
	},
}
