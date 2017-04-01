package commands

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
	"github.com/victorgama/gh/utils"
	"github.com/victorgama/go-octokit/octokit"
)

var listRepoLogger = utils.Logger.WithExtra("list")

// RepoList exposes a command responsible for listing repositories
var RepoList = cli.Command{
	Name:    "list",
	Aliases: []string{"l", "ls"},
	Usage:   "Lists your repositories",
	Action: func(c *cli.Context) error {
		listRepoLogger.Timing("Fetching repositories...")
		allRepos, err := utils.GetAllUserRepositories()
		if err != nil {
			if err, ok := err.(*octokit.ResponseError); ok {
				listRepoLogger.Error("%s", utils.FormatError(err))
				os.Exit(1)
			}
			return err
		}
		orgs := []string{}
		orgsRepos := map[string][]string{}
		repos := map[string]octokit.Repository{}

		for _, repo := range allRepos {
			orgName := repo.Owner.Login

			if _, present := orgsRepos[orgName]; !present {
				orgs = append(orgs, orgName)
				orgsRepos[orgName] = []string{}
			}
			orgsRepos[orgName] = append(orgsRepos[orgName], repo.Name)
			repos[repo.FullName] = repo
		}

		// Okay, first we want our repos, then, orgs in alphabetical order

		if username, present := utils.CurrentUserName(); present {
			// At this point, we depend on envvars. If absent, just move
			// on.
			username = strings.ToLower(username)
			index := -1
			for idx, n := range orgs {
				if strings.ToLower(n) == username {
					index = idx
					break
				}
			}
			if index > -1 {
				currentUser := orgs[index]
				orgs = append(orgs[:index], orgs[index+1:]...)
				sortedRepos := orgsRepos[currentUser]
				sort.Sort(utils.Alphabetic(sortedRepos))
				printRepositories(currentUser, sortedRepos, &repos, "You")
			}
		}

		sort.Sort(utils.Alphabetic(orgs))

		for _, orgName := range orgs {
			sortedRepos := orgsRepos[orgName]
			sort.Sort(utils.Alphabetic(sortedRepos))
			printRepositories(orgName, sortedRepos, &repos, "")
		}
		return nil
	},
}

func printRepositories(orgName string, names []string, repositories *map[string]octokit.Repository, alternativeName string) {
	fmt.Println("")
	name := orgName
	if alternativeName != "" {
		name = alternativeName
	}
	color.New(color.Bold, color.Underline).Printf("%s", name)
	color.New(color.Reset).Println("")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"⎇", "🔒", "Name", "URL"})
	table.SetAutoFormatHeaders(true)
	for _, name := range names {
		repo := (*repositories)[fmt.Sprintf("%s/%s", orgName, name)]
		fork := "  "
		access := "  "
		if repo.Fork {
			fork = "⎇"
		}
		if repo.Private {
			access = "🔒"
		}
		table.Append([]string{fork, access, repo.Name, fmt.Sprintf("https://github.com/%s", repo.FullName)})
	}
	table.Render()
}
