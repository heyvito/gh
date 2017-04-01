package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/victorgama/go-octokit/octokit"
	"github.com/victorgama/pine"
)

// NewClient creates a new Octokit client instance based on the access token
// obtained from the environment variables
func NewClient() *octokit.Client {
	return octokit.NewClient(octokit.TokenAuth{
		AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN"),
	})
}

// CurrentUserName attempts to get a GitHub username defined in environment
// variables
func CurrentUserName() (string, bool) {
	usr := os.Getenv("GITHUB_USERNAME")
	return usr, usr != ""
}

// UserIsOrg determines whether a given username is an organization
func UserIsOrg(name string) bool {
	client := NewClient()
	url, _ := octokit.UserURL.Expand(octokit.M{"user": name})
	e, _ := client.Users(url).One()
	return e.Type == "Organization"
}

// FormatError attempts to format a given octokit ResponseError
func FormatError(err *octokit.ResponseError) string {
	str := err.Message + "\n"
	for _, e := range err.Errors {
		str += fmt.Sprintf("    - %s: %s\n", e.Field, e.Message)
	}
	return str
}

// HandleClientError checks the state of an octokit Result, prints
// and exits the application if an error ocurred during the request
func HandleClientError(resp *octokit.Result, logger pine.Writer) {
	if resp.HasError() {
		err := resp.Err
		if err != nil {
			if serr, ok := err.(*octokit.ResponseError); ok {
				logger.Error("%s", FormatError(serr))
				// logger.Error("%s", string(debug.Stack()))
				os.Exit(1)
			} else {
				logger.Error("%s", err)
				// logger.Error("%s", string(debug.Stack()))
				os.Exit(1)
			}
		}
	}
}

// GetAllUserRepositories iterates all API pages and returns a list of repositories
// that belongs to the authenticated user
func GetAllUserRepositories() ([]octokit.Repository, error) {
	client := NewClient()
	result := []octokit.Repository{}

	repos, resp := client.Repositories().All(&octokit.UserRepositoriesURL, nil)
	for {
		if resp.HasError() {
			return nil, resp.Err
		}
		result = append(result, repos...)
		if resp.NextPage != nil {
			repos, resp = client.Repositories().All(resp.NextPage, nil)
		} else {
			break
		}
	}
	return result, nil
}

// GetAllCollabs returns a list of all collaborators of a given repository
func GetAllCollabs(url *RepoURL) ([]octokit.User, error) {
	client := NewClient()
	result := []octokit.User{}

	users, resp := client.Collaborators().All(&octokit.CollaboratorsURL, octokit.M{"owner": url.Username, "repo": url.RepoName})
	for {
		if resp.HasError() {
			return nil, resp.Err
		}
		result = append(result, users...)
		if resp.NextPage != nil {
			users, resp = client.Collaborators().All(resp.NextPage, nil)
		} else {
			break
		}
	}
	return result, nil
}

// GetAllTeamsForRepo returns a list of teams that have access to a given repository
func GetAllTeamsForRepo(url *RepoURL) ([]octokit.Team, error) {
	client := NewClient()
	result := []octokit.Team{}

	teams, resp := client.Teams().GetTeamsForRepository(&octokit.TeamsRepositoryURL, octokit.M{"owner": url.Username, "repo": url.RepoName})
	for {
		if resp.HasError() {
			return nil, resp.Err
		}
		result = append(result, teams...)
		if resp.NextPage != nil {
			teams, resp = client.Teams().All(resp.NextPage, nil)
		} else {
			break
		}
	}
	return result, nil
}

// GetAllTeamsForOrg returns a list of all teams for a given organization
func GetAllTeamsForOrg(org string) ([]octokit.Team, error) {
	client := NewClient()
	result := []octokit.Team{}

	teams, resp := client.Teams().GetTeamsForRepository(&octokit.OrganizationTeamsURL, octokit.M{"org": org})
	for {
		if resp.HasError() {
			return nil, resp.Err
		}
		result = append(result, teams...)
		if resp.NextPage != nil {
			teams, resp = client.Teams().All(resp.NextPage, nil)
		} else {
			break
		}
	}
	return result, nil
}

// GetTeamByName returns a Team instance belonging to a given organization under a given name
func GetTeamByName(org, team string, logger pine.Writer, exitOnError bool) (*octokit.Team, error) {
	teams, err := GetAllTeamsForOrg(org)
	if err != nil {
		if err, ok := err.(*octokit.ResponseError); ok {
			logger.Error("%s", FormatError(err))
			os.Exit(1)
		}
		return nil, err
	}
	teamName := strings.ToLower(team)
	var t *octokit.Team
	for _, rt := range teams {
		if strings.ToLower(rt.Slug) == teamName {
			t = &rt
			break
		}
	}
	if t == nil && exitOnError {
		logger.Warn("Could not find a team named '%s' on the organization '%s'", team, org)
		os.Exit(1)
	}
	return t, nil
}

// GetTeamMembers returns a list of members of a given organization team
func GetTeamMembers(org, team string, logger pine.Writer) ([]octokit.User, *octokit.Team, error) {
	t, err := GetTeamByName(org, team, logger, true)
	if err != nil {
		return nil, nil, err
	}
	client := NewClient()
	members, resp := client.Teams().GetMembers(&octokit.TeamMembersURL, octokit.M{"id": t.ID})
	HandleClientError(resp, logger)
	return members, t, nil
}

// GetAllUsersForOrg returns a list of users that belongs to an organization
func GetAllUsersForOrg(org string) ([]octokit.User, error) {
	client := NewClient()
	result := []octokit.User{}

	users, resp := client.Organization().GetOrganizationMembers(&octokit.OrganizationMembersURL, octokit.M{"org": org})
	for {
		if resp.HasError() {
			return nil, resp.Err
		}
		result = append(result, users...)
		if resp.NextPage != nil {
			users, resp = client.Organization().GetOrganizationMembers(resp.NextPage, nil)
		} else {
			break
		}
	}
	return result, nil
}
