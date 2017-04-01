package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
	"github.com/victorgama/gh/utils"
	"github.com/victorgama/go-octokit/octokit"
)

var teamsLogger = utils.Logger.WithExtra("teams")

var teamsList = cli.Command{
	Name:      "list",
	Usage:     "Lists teams for a given organization",
	ArgsUsage: "[org]",
	Action: func(c *cli.Context) error {
		if len(c.Args()) != 1 {
			return fmt.Errorf("usage: gh teams list [org]")
		}
		teamsLogger.Timing("One moment, please...")
		teams, err := utils.GetAllTeamsForOrg(c.Args()[0])
		if err != nil {
			if err, ok := err.(*octokit.ResponseError); ok {
				collabLogger.Error("%s", utils.FormatError(err))
				os.Exit(1)
			}
			return err
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Team", "Description", "Privacy", "Permission"})
		table.SetAutoFormatHeaders(true)
		for _, team := range teams {
			table.Append([]string{fmt.Sprintf("%s (@%s)", team.Name, team.Slug), team.Description, team.Privacy, team.Permission})
		}
		table.Render()
		return nil
	},
}

var teamsMembers = cli.Command{
	Name:      "members",
	Usage:     "Lists members for a given organization team",
	ArgsUsage: "[org] [team-slug]",
	Action: func(c *cli.Context) error {
		if len(c.Args()) != 2 {
			return fmt.Errorf("usage: gh teams members [org] [team-slug]")
		}
		teamsLogger.Timing("One moment, please...")
		members, _, err := utils.GetTeamMembers(c.Args()[0], c.Args()[1], teamsLogger)
		if err != nil {
			return err
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Users"})
		table.SetAutoFormatHeaders(true)
		for _, member := range members {
			table.Append([]string{fmt.Sprintf("@%s", member.Login)})
		}
		table.Render()
		return nil
	},
}

var teamAddUser = cli.Command{
	Name:      "add",
	Usage:     "Adds a user to a team",
	ArgsUsage: "[username](:role) [org] [team-slug]",
	Action: func(c *cli.Context) error {
		if len(c.Args()) != 3 {
			return fmt.Errorf("usage: gh teams add [username](:role) [org] [team-slug]")
		}
		role := "member"
		username := c.Args()[0]
		if strings.Contains(username, ":") {
			split := strings.Split(username, ":")
			role = strings.ToLower(split[1])
			username = split[0]
		}
		if role != "member" && role != "maintainer" {
			return fmt.Errorf("when defining a role, please specify either 'maintainer' or 'member'")
		}
		orgName := c.Args()[1]
		teamName := c.Args()[2]
		teamsLogger.Timing("One moment, please...")
		members, team, err := utils.GetTeamMembers(orgName, teamName, teamsLogger)
		if err != nil {
			return nil
		}

		url, err := octokit.UserURL.Expand(octokit.M{"user": username})
		if err != nil {
			return err
		}

		client := utils.NewClient()
		user, resp := client.Users(url).One()
		utils.HandleClientError(resp, teamsLogger)

		userPresent := false
		lowerCaseUser := strings.ToLower(username)
		for _, user := range members {
			if strings.ToLower(user.Login) == lowerCaseUser {
				userPresent = true
				break
			}
		}

		fmt.Println("")
		fmt.Println("Hey there! You're about to add a new user to a team.")
		fmt.Printf("Mind checking out whether this is the @%s you're looking for?\n", username)
		fmt.Println("")
		fmt.Printf("        Name: %s\n", user.Name)
		fmt.Printf("       Email: %s\n", user.Email)
		fmt.Printf("Organization: %s\n", user.Company)
		fmt.Printf("         URL: %s\n", user.Blog)
		fmt.Println("")
		if !userPresent {
			fmt.Printf("⚠️  WARNING! Continuing will invite @%s to the %s organization!\n\n", username, orgName)
		}
		if !utils.Confirm(fmt.Sprintf("Continue adding @%s? [y]/n", username), true) {
			return fmt.Errorf("aborting")
		}

		teamsLogger.Timing("One moment, please...")
		_, resp = client.Teams().AddMembership(&octokit.TeamMembershipURL, octokit.M{"id": team.ID, "username": username}, role)
		utils.HandleClientError(resp, teamsLogger)
		teamsLogger.Success("Added @%s to %s/%s", username, orgName, team.Slug)
		return nil
	},
}

var teamRmUser = cli.Command{
	Name:      "rm",
	Usage:     "Removes a user to a team",
	ArgsUsage: "[username] [org] [team-slug]",
	Action: func(c *cli.Context) error {
		if len(c.Args()) != 3 {
			return fmt.Errorf("usage: gh teams rm [username] [org] [team-slug]")
		}
		username := c.Args()[0]
		orgName := c.Args()[1]
		teamName := c.Args()[2]
		teamsLogger.Timing("One moment, please...")
		members, team, err := utils.GetTeamMembers(orgName, teamName, teamsLogger)
		if err != nil {
			return nil
		}

		url, err := octokit.UserURL.Expand(octokit.M{"user": username})
		if err != nil {
			return err
		}

		client := utils.NewClient()
		_, resp := client.Users(url).One()
		utils.HandleClientError(resp, teamsLogger)

		userPresent := false
		lowerCaseUser := strings.ToLower(username)
		for _, user := range members {
			if strings.ToLower(user.Login) == lowerCaseUser {
				userPresent = true
				break
			}
		}

		if !userPresent {
			teamsLogger.Warn("@%s does not belong to %s/%s", orgName, team.Slug)
			os.Exit(1)
		}

		fmt.Println("")
		fmt.Println("Hey there! You're about to remove an user from a team.")
		if !utils.Confirm(fmt.Sprintf("Continue removing @%s from %s/%s? y/[n]", username, orgName, teamName), false) {
			return fmt.Errorf("aborting")
		}

		teamsLogger.Timing("One moment, please...")
		_, resp = client.Teams().RemoveMembership(&octokit.TeamMembershipURL, octokit.M{"id": team.ID, "username": username})
		utils.HandleClientError(resp, teamsLogger)
		teamsLogger.Success("Removed @%s from %s/%s", username, orgName, team.Slug)
		return nil
	},
}

// Teams exposes team-related commands
var Teams = cli.Command{
	Name:    "teams",
	Aliases: []string{"t"},
	Usage:   "Manages organization teams",
	Subcommands: []cli.Command{
		teamsList,
		teamsMembers,
		teamAddUser,
		teamRmUser,
	},
}
