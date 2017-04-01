package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli"
	"github.com/victorgama/gh/utils"
	"github.com/victorgama/go-octokit/octokit"
)

var rmRepoLogger = utils.Logger.WithExtra("rm")

// RmRepo exposes a command responsible for deleting a repository
var RmRepo = cli.Command{
	Name:  "rm",
	Usage: "Destroys a repository",
	Action: func(c *cli.Context) error {
		if len(c.Args()) != 1 {
			return fmt.Errorf("the 'rm' command requires one repository name")
		}
		client := utils.NewClient()
		rmRepoLogger.Timing("One moment, please...")
		r := utils.RepoURLFromString(c.Args()[0])
		r.AutoComplete()

		repo, resp := client.Repositories().One(&octokit.RepositoryURL, octokit.M{"owner": r.Username, "repo": r.RepoName})
		utils.HandleClientError(resp, rmRepoLogger)

		fmt.Println("Hey! You're about to perform a really dangerous action.")
		fmt.Printf("To confirm you really want to delete %s, please enter its name again:\n", repo.FullName)
		fmt.Print("What is its name again? ")
		reader := bufio.NewReader(os.Stdin)
		s, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		s = strings.TrimSpace(s)
		s = strings.ToLower(s)
		if s != strings.ToLower(repo.Name) {
			fmt.Println("Nope. That's not its name. Aborting.")
			return fmt.Errorf("aborted")
		}

		rmRepoLogger.Info("Removing %s...", repo.FullName)
		success, resp := client.Repositories().Delete(&octokit.RepositoryURL, octokit.M{"owner": r.Username, "repo": r.RepoName})
		if !success {
			utils.HandleClientError(resp, rmRepoLogger)
		}
		rmRepoLogger.Success("Removed %s", repo.FullName)
		return nil
	},
}
