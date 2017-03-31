package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/octokit/go-octokit/octokit"
)

// NormalizeRepoName tries to normalize a repository name by replacing special
// characters and stripping a .git suffix from it.
func NormalizeRepoName(in string) (string, bool) {
	dashfier := regexp.MustCompile(`[^\w\.\-]+`)
	gitStripper := regexp.MustCompilePOSIX(`(\.git)+$`)
	newIn := dashfier.ReplaceAllString(in, "-")
	newIn = gitStripper.ReplaceAllString(newIn, "")
	return newIn, newIn != in
}

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

func RepoURLFromString(s string) RepoURL {
	if strings.Contains(s, "/") {
		r := regexp.MustCompile(`(?:([^\/]+)/?)(.*)`)
		groups := r.FindStringSubmatch(s)
		return RepoURL{
			Username: groups[1],
			RepoName: groups[2],
		}
	}
	return RepoURL{
		RepoName: s,
	}
}

// RepoURL represents a username/reponame structure
type RepoURL struct {
	Username string
	RepoName string
}

// ToURL transforms the RepoURL into a username/repo string
func (r *RepoURL) ToURL() string {
	usr, hasUsr := CurrentUserName()
	if !hasUsr && r.Username == "" {
		fmt.Println("To use short-format repository names, please define your GitHub username using")
		fmt.Println("the GITHUB_USERNAME environment variable.")
		os.Exit(1)
	}
	if r.Username == "" {
		r.Username = usr
	}
	return fmt.Sprintf("%s/%s", r.Username, r.RepoName)
}

func UserIsOrg(name string) bool {
	client := NewClient()
	url, _ := octokit.UserURL.Expand(octokit.M{"user": name})
	e, _ := client.Users(url).One()
	return e.Type == "Organization"
}
