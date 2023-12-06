# Commit check

A quick wrapper around git that formats the commits in a conventional style.

The benefit to using this over other tooling such as commitizen is that this
can be compiled down to a binary and doesn't therefore require node or python
on the machine.
It also doesn't require others on the team to use the standardized format,
although you will get more out of other tooling in the space if everyone does
use the same commit structure.

## Usage

<img alt="Example usage of commit-check" src="examples/basic.gif" width="600" />

The application will prompt you for inputs at each point.
Anything added as arguments will be ignored.

``` bash
commit-check
```

You can specify what prefixes for commit types you are going to follow by
using the `-l` or `type-list` flag.
This only accepts `angular` or `conventionalcommit` currently.
Raise an issue or contribute directly if there are other types that would be
useful, as well as links to documentation showing where the convention comes
from.

``` bash
commit-check -l angular
commit-check --type-list conventionalcommit
```

You can also add `--dry-run` flag, or shorthand of `-d`.

<img alt="Example usage of commit-check with all arguments" src="examples/dry-run.gif" width="600" />

## Configuration

A configuration file can be used to set default values, and
can be located either at `/etc/commit-check/config.yaml`, or
`$HOME/.commit-check/config.yaml`.
A full example can be found below.

```yaml
plugins:
  jira:
    url: "https://company.atlassian.net"
    username: "yourname@company.com"
    apiKey: "0123"
    project:
      - "ABC"
    status:
      - "In Progress"
  clickup:
    apiKey: "0123"
    listIds:
      - "0123"
    assignee: "yourname@company.com"
```

## Plugins

commit-check can retrieve cards assigned to you for different
systems such as Jira and Clickup, and add those to the list of
available scopes.

These will not be activated unless there is a valid configuration.

### Jira

Required configuration:
- `url` - The url of the Jira server.
- `username` - The users email.
- `apiKey` - A PAT for the user.

Optional configuration:
- `projects` - Filter the cards to only those in the given projects/
  List of strings.
- `status` - Filter the card to only those in the given statuses.
  List of strings.
 
### ClickUp

Required configuration:
- `apiKey` - A PAT for the given user.
- `listIds` - A list of IDs to search for tasks in.
  You will likely need to use the API to identify which lists are
  the ones you want to use.
  List of strings.

Optional configuration:
- `assignee` - Filter the card to only those assigned this email.
