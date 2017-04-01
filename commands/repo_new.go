package commands

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
	"github.com/victorgama/gh/utils"
	"github.com/victorgama/go-octokit/octokit"
)

var newRepoLogger = utils.Logger.WithExtra("new")

// NewRepo is responsible for creating repositories
var NewRepo = cli.Command{
	Name:      "new",
	Usage:     "Creates a new repository",
	ArgsUsage: "(--private) (--init) (--license LICENSE) (--gitignore LANGUAGE_OR_PLATFORM) [name, ...]",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "private",
			Usage: "creates a private repository",
		},
		cli.BoolFlag{
			Name:  "init",
			Usage: "initializes the repository with a commit and empty README",
		},
		cli.StringFlag{
			Name:  "license",
			Usage: "selects a license template to apply. Use a license name from https://github.com/github/choosealicense.com/tree/gh-pages/_licenses, without the extension (For example, 'mit' or 'mozilla')",
			Value: "",
		},
		cli.StringFlag{
			Name:  "gitignore",
			Usage: "selects a gitignore template to apply. Use a language or platform name from https://github.com/github/gitignore withtout the extension (For example, 'Haskell')",
			Value: "",
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
			r.AutoComplete()
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
				Name:              re.RepoName,
				Private:           private,
				GitIgnoreTemplate: c.String("gitignore"),
				LicenseTemplate:   c.String("license"),
			})
			if res.HasError() {
				if err, ok := res.Err.(*octokit.ResponseError); ok {
					newRepoLogger.Error("%s", utils.FormatError(err))
					os.Exit(1)
				} else {
					panic(err)
				}
			}
			newRepoLogger.Success("Created: %s", repo.FullName)
		}
		return nil
	},
}
