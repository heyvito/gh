package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
	"github.com/victorgama/gh/utils"
	"github.com/victorgama/go-octokit/octokit"
)

var collabLogger = utils.Logger.WithExtra("collab")

var collabAdd = cli.Command{
	Name:      "add",
	Usage:     "Adds a user or team to a repository",
	ArgsUsage: "[repository] [team-slug|contributor-username](:permission)",
	Action: func(c *cli.Context) error {
		if len(c.Args()) != 2 {
			return fmt.Errorf("Usage: gh collab add [repository] [team-slug|contributor-username](:permission)")
		}

		role := ""
		toAdd := strings.ToLower(c.Args()[1])
		if strings.Contains(toAdd, ":") {
			split := strings.Split(toAdd, ":")
			role = split[1]
			toAdd = split[0]

			roleAliases := map[string]string{
				"read":  "pull",
				"write": "push",
			}

			if v, present := roleAliases[role]; present {
				role = v
			}

			if role != "pull" && role != "push" && role != "admin" {
				return fmt.Errorf("incorrect role %s: valid roles are pull/read, push/write and admin", role)
			}
		}

		repoURL := utils.RepoURLFromString(c.Args()[0])
		repoURL.AutoComplete()

		collabLogger.Timing("Just a second...")
		isOrg := utils.UserIsOrg(repoURL.Username)

		if !isOrg && role != "" {
			return fmt.Errorf("cannot set permission level on a non-org repository collaborator")
		}

		var target interface{}

		if isOrg {
			// At this point, we don't know whether the collaborator is a team or user
			t, err := utils.GetTeamByName(repoURL.Username, toAdd, collabLogger, false)
			if err != nil {
				return err
			}
			if t != nil {
				target = t
			} else {
				collabLogger.Warn("No team found under %s/%s. Looking for users...", repoURL.Username, toAdd)
			}
		}

		client := utils.NewClient()
		if target == nil {
			url, err := octokit.UserURL.Expand(octokit.M{"user": toAdd})
			if err != nil {
				return err
			}

			u, resp := client.Users(url).One()
			if resp.Response != nil && resp.Response.StatusCode != 404 {
				utils.HandleClientError(resp, collabLogger)
			}
			if u == nil {
				collabLogger.Warn("No user found with handle @%s. Aborting.", toAdd)
				return fmt.Errorf("aborted")
			}
			if isOrg {
				collabLogger.Warn("WARNING: Adding user to org repository as an outside collaborator!")
			}
			target = u
		}

		if isOrg {
			if t, ok := target.(*octokit.Team); ok {
				collabLogger.Timing("Adding %s/%s to %s", repoURL.Username, t.Slug, repoURL.ToURL())
				success, resp := client.Teams().UpdateRepository(&octokit.TeamRepositoryURL, octokit.M{"id": t.ID, "owner": repoURL.Username, "repo": repoURL.RepoName}, role)
				if !success {
					utils.HandleClientError(resp, collabLogger)
				}
				collabLogger.Success("Added %s/%s to %s", repoURL.Username, t.Slug, repoURL.ToURL())
				return nil
			}
			u := target.(*octokit.User)
			users, err := utils.GetAllUsersForOrg(repoURL.Username)
			if err != nil {
				return err
			}
			userIsPresent := false
			for _, orgUsr := range users {
				if u.Login == orgUsr.Login {
					userIsPresent = true
					break
				}
			}
			if userIsPresent {
				collabLogger.Warn("WARNING: Adding org user as an outside collaborator!")
			} else {
				fmt.Println("Hey there! You're about to add an outside user to an org repository.")
				fmt.Printf("Mind checking out whether this is the @%s you're looking for?\n", toAdd)
				fmt.Println("")
				fmt.Printf("        Name: %s\n", u.Name)
				fmt.Printf("       Email: %s\n", u.Email)
				fmt.Printf("Organization: %s\n", u.Company)
				fmt.Printf("         URL: %s\n", u.Blog)
				fmt.Println("")
				if !utils.Confirm(fmt.Sprintf("Continue adding @%s? y/[n]", toAdd), false) {
					return fmt.Errorf("aborting")
				}
			}
		}
		if role == "" {
			collabLogger.Info("Defaulting user permission level to 'push'")
			role = "push"
		}
		u := target.(*octokit.User)
		collabLogger.Timing("Adding @%s to %s", u.Login, repoURL.ToURL())
		success, resp := client.Collaborators().AddCollaborator(&octokit.CollaboratorsURL, octokit.M{"owner": repoURL.Username, "repo": repoURL.RepoName, "username": u.Login}, role)
		if !success {
			utils.HandleClientError(resp, collabLogger)
		}
		collabLogger.Success("Added %s to %s", u.Login, repoURL.ToURL())
		return nil
	},
}

var collabRm = cli.Command{
	Name:      "rm",
	Usage:     "Removes a user or team from a repository",
	ArgsUsage: "[repository] [team-slug|contributor-username]",
	Action: func(c *cli.Context) error {
		if len(c.Args()) != 2 {
			return fmt.Errorf("Usage: gh collab rm [repository] [team-slug|contributor-username]")
		}

		toRm := strings.ToLower(c.Args()[1])
		repoURL := utils.RepoURLFromString(c.Args()[0])
		repoURL.AutoComplete()

		collabLogger.Timing("Just a second...")
		isOrg := utils.UserIsOrg(repoURL.Username)

		var target interface{}

		if isOrg {
			// At this point, we don't know whether the collaborator is a team or user
			t, err := utils.GetTeamByName(repoURL.Username, toRm, collabLogger, false)
			if err != nil {
				return err
			}
			if t != nil {
				target = t
			} else {
				collabLogger.Warn("No team found under %s/%s. Looking for users...", repoURL.Username, toRm)
			}
		}

		client := utils.NewClient()
		if target == nil {
			url, err := octokit.UserURL.Expand(octokit.M{"user": toRm})
			if err != nil {
				return err
			}

			u, resp := client.Users(url).One()
			if resp.Response != nil && resp.Response.StatusCode != 404 {
				utils.HandleClientError(resp, collabLogger)
			}
			if u == nil {
				collabLogger.Warn("No user found with handle @%s. Aborting.", toRm)
				return fmt.Errorf("aborted")
			}
			target = u
		}

		if isOrg {
			if t, ok := target.(*octokit.Team); ok {
				collabLogger.Timing("Removing %s/%s from %s", repoURL.Username, t.Slug, repoURL.ToURL())
				success, resp := client.Teams().RemoveRepository(&octokit.TeamRepositoryURL, octokit.M{"id": t.ID, "owner": repoURL.Username, "repo": repoURL.RepoName})
				if !success {
					utils.HandleClientError(resp, collabLogger)
				}
				collabLogger.Success("Removed %s/%s from %s", repoURL.Username, t.Slug, repoURL.ToURL())
				return nil
			}
		}
		u := target.(*octokit.User)
		collabLogger.Timing("Removing @%s from %s", u.Login, repoURL.ToURL())
		success, resp := client.Collaborators().RemoveCollaborator(&octokit.CollaboratorsURL, octokit.M{"owner": repoURL.Username, "repo": repoURL.RepoName, "username": u.Login})
		if !success {
			utils.HandleClientError(resp, collabLogger)
		}
		collabLogger.Success("Removed %s from %s", u.Login, repoURL.ToURL())
		return nil
	},
}

var collabList = cli.Command{
	Name:      "list",
	Usage:     "Lists teams and/or contributors on a given repository",
	ArgsUsage: "[repository]",
	Action: func(c *cli.Context) error {
		if len(c.Args()) != 1 {
			return fmt.Errorf("expecting a repository name as argument. Aborting")
		}
		repoURL := utils.RepoURLFromString(c.Args()[0])
		repoURL.AutoComplete()

		collabLogger.Timing("Just a second...")

		if utils.UserIsOrg(repoURL.Username) {
			collabLogger.Timing("Fetching teams for %s", repoURL.ToURL())
			teams, err := utils.GetAllTeamsForRepo(&repoURL)
			if err != nil {
				if err, ok := err.(*octokit.ResponseError); ok {
					collabLogger.Error("%s", utils.FormatError(err))
					os.Exit(1)
				}
				return err
			}
			if len(teams) > 1 {
				fmt.Println("")
				color.New(color.Bold, color.Underline).Println(repoURL.ToURL())
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Team", "Permission"})
				table.SetAutoFormatHeaders(true)
				for _, team := range teams {
					table.Append([]string{fmt.Sprintf("%s (@%s)", team.Name, team.Slug), team.Permission})
				}
				table.Render()
			} else {
				collabLogger.Warn("No teams defined for %s. Falling back to contributors list...", repoURL.ToURL())
				collabs, err := utils.GetAllCollabs(&repoURL)
				if err != nil {
					if err, ok := err.(*octokit.ResponseError); ok {
						collabLogger.Error("%s", utils.FormatError(err))
						os.Exit(1)
					}
					return err
				}
				if len(collabs) > 0 {
					fmt.Println("")
					color.New(color.Bold, color.Underline).Println(repoURL.ToURL())
					table := tablewriter.NewWriter(os.Stdout)
					table.SetHeader([]string{"User", "Push?", "Pull?", "Admin?"})
					table.SetAutoFormatHeaders(true)
					for _, collab := range collabs {
						push := "No"
						pull := "No"
						admin := "No"
						if collab.Permissions.Admin {
							admin = "Yes"
						}
						if collab.Permissions.Pull {
							pull = "Yes"
						}
						if collab.Permissions.Push {
							push = "Yes"
						}
						table.Append([]string{fmt.Sprintf("@%s", collab.Login), push, pull, admin})
					}
					table.Render()
				} else {
					fmt.Println("No collaborators")
				}
			}
			fmt.Println("")
		} else {
			collabLogger.Timing("Fetching collaborators for %s", repoURL.ToURL())
			fmt.Println("")
			collabs, err := utils.GetAllCollabs(&repoURL)
			if err != nil {
				if err, ok := err.(*octokit.ResponseError); ok {
					collabLogger.Error("%s", utils.FormatError(err))
					os.Exit(1)
				}
				return err
			}
			color.New(color.Bold, color.Underline).Println(repoURL.ToURL())
			if len(collabs) > 0 {
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"User", "Push?", "Pull?", "Admin?"})
				table.SetAutoFormatHeaders(true)
				for _, collab := range collabs {
					push := "No"
					pull := "No"
					admin := "No"
					if collab.Permissions.Admin {
						admin = "Yes"
					}
					if collab.Permissions.Pull {
						pull = "Yes"
					}
					if collab.Permissions.Push {
						push = "Yes"
					}
					table.Append([]string{fmt.Sprintf("@%s", collab.Login), push, pull, admin})
				}
				table.Render()
			} else {
				fmt.Println("No collaborators")
			}
		}
		return nil
	},
}

// Collab provides collaboration-related commands
var Collab = cli.Command{
	Name:    "collab",
	Aliases: []string{"c"},
	Usage:   "Collaborator actions",
	Subcommands: []cli.Command{
		collabAdd,
		collabRm,
		collabList,
	},
}
