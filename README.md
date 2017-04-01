# gh
##### Control GitHub from your Terminal

<center>
<img src="https://www.dropbox.com/s/u7m2v5dnpa47uzf/gh.png?dl=1" />
</center>

`gh` is a terminal utility that allows you to use GitHub directly from your terminal.
No more fiddling with the browser to create a new repository or managing teams.

## Installing
Ezpz!
```
$ go install -u github.com/victorgama/gh
```

## Configuring
To use this tool, two environment variables must be set: `GITHUB_ACCESS_TOKEN` and `GITHUB_USERNAME`.

1. Generate a new [Personal Access Token](https://github.com/settings/tokens) and store it in the `GITHUB_ACCESS_TOKEN` variable
2. Ensure to set another envvar named `GITHUB_USERNAME` with your GitHub username.
3. Done!


## Usage

> **Notice**: Assume the value of `GITHUB_USERNAME` is `octocat`.

### Creating a repository
```
gh new --private --license mit --gitignore Go my-repo
```

This will create a new private repository, `octocat/my-repo`, initialized with a MIT license file and a `.gitignore`

```
gh new --private --license mit --gitignore Go github/my-repo
```

If a complete repo name is provided, `gh` will attempt to create a new repo under the provided path. When the owner name is absent, `GITHUB_USERNAME` is assumed as its value.

### Deleting a repository
```
gh rm my-repo
```

`gh rm` will ask for confirmation and then will delete a repository under `octocat/my-repo`

```
gh rm github/my-repo
```

`gh rm` follows the same path convention as `gh new`

### Listing repositories
```
gh l | gh ls | gh list
```

This will list all repositories under you account and organizations you have access

### Collaboration management

`gh` allows you to manage your collaborators and teams.

#### Listing
```
gh collab list my-repo
gh collab list octocat/my-repo
```

- For an user repository:
    - Lists all collaborators and their permissions
- For an organization repository:
    - Lists all teams that have access to this repository
    - Lists all outside collaborators for the given repository

#### Adding collaborators and teams
```
gh collab add [repository] [team-slug|contributor-username](:permission-level)
gh collab add github/secret design:write
```
Adds an user or team to an user or org repository.

- For an user repository:
    - Assumes the provided value is an username, searches for it and adds them as a collaborator
    > **Notice**: It is not possible to use a `:permission-level` for user repositories.
- For an organization repository:
    - Assumes the provided value is a team slug. Searches for it and adds it under the given `:permission-level`.
    - If `permission-level` is absent, assumes the team permission level.
    - If the target organization does not have a team with the provided name, searches for users, and adds them under the given `:permission-level`
    - Assumes `push` as the permission, if absent.
Valid values for `permission-level` are `read|pull`, `write|push`, and `admin`

> **Notice**: You can also use this command to update a user or team permission level. ✨

> **Protip**: `gh` will ask for confirmation if the operation may cause unintended results.

#### Removing collaborators
```
gh collab rm [repository] [team-slug|contributor-username]
gh collab rm github/secret octocat
```
Removes an user or team from an user or org repository.

- For an user repository:
    - Assumes the provided value is an username, and attempts to remove it from the repository
- For an organization repository:
    - Assumes the provided value is a team slug. Searches for it in the repo's organization and if found, removes it from the repository.
    - If the target organization does not have a team with the provided name, searches for users, and removes them from the repo.

### Team management

#### Listing teams
```
gh teams list [org]
gh teams list github
```
Lists all teams for a given organization. Includes its name, slug, description, privacy and permission options.

#### Listing members
```
gh teams members [org] [team-slug]
gh teams members github design
```
Lists handles for users in a given organization team

#### Adding members to a team
```
gh teams add [username](:role) [org] [team-slug]
gh teams add octocat github design
```
Adds a given `username`, under an optional `role`, into a team under `team-slug` on the `org` organization.

Please note that if `username` is not yet a member of the target organization, an invite will sent to them.
Permitted role values are `member` and `maintainer`.

#### Removing members from a team
```
gh teams rm [username] [org] [team-slug]
gh teams rm octocat github design
```
Removes a given `username` from a team under `team-slug` on the `org` organization.

### Silly utilities

#### Quickly opening a repository
```
gh open [repo] | gh o [repo]
gh open github/octicons
gh open my-repo
```

Opens a repository on GitHub.com on your default browser

## TODO
- [ ] Add tests
- [ ] Support GitHub Enterprise (?)
- [ ] Improve help topics

## License
```
MIT License

Copyright (c) 2017 Victor Gama

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
