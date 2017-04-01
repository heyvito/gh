package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// RepoURL represents a username/reponame structure
type RepoURL struct {
	Username string
	RepoName string
}

// NormalizeRepoName tries to normalize a repository name by replacing special
// characters and stripping a .git suffix from it.
func NormalizeRepoName(in string) (string, bool) {
	dashfier := regexp.MustCompile(`[^\w\.\-]+`)
	gitStripper := regexp.MustCompilePOSIX(`(\.git)+$`)
	newIn := dashfier.ReplaceAllString(in, "-")
	newIn = gitStripper.ReplaceAllString(newIn, "")
	return newIn, newIn != in
}

// AutoCompleteRepoName attempts to autocomplete an incomplete repository
// name using the GITHUB_USERNAME environment variable
func AutoCompleteRepoName(repo string) string {
	r := RepoURLFromString(repo)
	return r.ToURL()
}

// RepoURLFromString creates a new RepoURL struct from a given string
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

// AutoComplete attempts to autocomplete a RepoURL instance using the
// GITHUB_USERNAME environment variable
func (r *RepoURL) AutoComplete() {
	usr, hasUsr := CurrentUserName()
	if !hasUsr && r.Username == "" {
		fmt.Println("To use short-format repository names, please define your GitHub username using")
		fmt.Println("the GITHUB_USERNAME environment variable.")
		os.Exit(1)
	}
	if r.Username == "" {
		r.Username = usr
	}
}

// ToURL transforms the RepoURL into a username/repo string
func (r *RepoURL) ToURL() string {
	r.AutoComplete()
	return fmt.Sprintf("%s/%s", r.Username, r.RepoName)
}
