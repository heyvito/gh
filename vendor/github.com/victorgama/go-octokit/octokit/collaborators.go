package octokit

// CollaboratorsURL is the template for accessing the collaborators
// of a particular repository.
var (
	CollaboratorsURL = Hyperlink(
		"repos/{owner}/{repo}/collaborators{/username}")
)

// Collaborators creates a CollaboratorsService with a base url
func (c *Client) Collaborators() (repos *CollaboratorsService) {
	repos = &CollaboratorsService{client: c}
	return
}

// CollaboratorsService is a service providing access to a repositories'
// collaborators
type CollaboratorsService struct {
	client *Client
}

// All lists all the collaborating users on the given repository
//
// https://developer.github.com/v3/repos/collaborators/#list
func (r *CollaboratorsService) All(uri *Hyperlink, params M) (users []User,
	result *Result) {
	if uri == nil {
		uri = &CollaboratorsURL
	}
	url, err := uri.Expand(params)
	if err != nil {
		return nil, &Result{Err: err}
	}
	result = r.client.get(url, &users)
	return
}

// IsCollaborator checks if a user is a collaborator for a repo
//
// https://developer.github.com/v3/repos/collaborators/#check-if-a-user-is-a-collaborator
func (r *CollaboratorsService) IsCollaborator(uri *Hyperlink,
	params M) (collabStatus bool, result *Result) {
	if uri == nil {
		uri = &CollaboratorsURL
	}
	url, err := uri.Expand(params)
	if err != nil {
		return false, &Result{Err: err}
	}
	result = r.client.get(url, nil)
	collabStatus = false
	if result.Err == nil && result.Response.Response.StatusCode == 204 {
		collabStatus = true
	}
	return
}

func (t *CollaboratorsService) AddCollaborator(uri *Hyperlink, uriParams M,
	permission string) (success bool, result *Result) {
	url, err := ExpandWithDefault(uri, &CollaboratorsURL, uriParams)
	if err != nil {
		return false, &Result{Err: err}
	}
	if permission == "" {
		result = t.client.put(url, nil, nil)
	} else {
		result = t.client.put(url, M{"permission": permission}, nil)
	}
	success = (result.Response.StatusCode == 204)
	return
}

func (t *CollaboratorsService) RemoveCollaborator(uri *Hyperlink, uriParams M) (
	success bool, result *Result) {
	url, err := ExpandWithDefault(uri, &CollaboratorsURL, uriParams)
	if err != nil {
		return false, &Result{Err: err}
	}
	result = t.client.delete(url, nil, nil)
	success = (result.Response.StatusCode == 204)
	return
}
