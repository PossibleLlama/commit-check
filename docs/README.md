# Commit check

A quick wrapper around git that formats the commits in a conventional style.

The benefit to using this over other tooling such as commitizen is that this
can be compiled down to a binary and doesn't therefore require node or python
on the machine.
It also doesn't require others on the team to use the standardized format,
although you will get more out of other tooling in the space if everyone does
use the same commit structure.

## Usage

The application will show you the output, and let you input each part
such as the type of commit, what the scope of changes is, and a description
of the changes.
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

The scope can either be ignored, a manual entry be made, or (if you have
setup the configuration for it) a list of items from Jira and/or Clickup
with their IDs shown.
This can let you look at a list of scopes for changes that you are working
on, and quickly select from the most likely ones, or add in a value
yourself if the one you want isnt present.

You can also add `--dry-run` flag, or shorthand of `-d`.
When the application is running, on the summary page, pressing `D` will
toggle whether dry run is enabled or not.

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

## Verification

To verify that the downloaded artifacts were built from
this repository, you can run a check using the [Github CLI](https://cli.github.com/manual/gh_attestation_verify).

If these commands fail, it indicates that they were not
built correctly, and it's advised that you download the
artifacts directly from the [releases page](https://github.com/PossibleLlama/commit-check/releases).

To verify that the downloaded artifact originates from
this repo, you can run the following against the binary.

``` bash
gh attestation verify --owner "PossibleLlama" commit-check.bin
```

The build uses [this reusable workflow](https://github.com/PossibleLlama/workflows/blob/main/.github/workflows/golang-binary.yaml),
and as such the build can be verified that it was built
from that workflow.

```bash
gh attestation verify --owner PossibleLlama --signer-workflow "PossibleLlama/workflows/.github/workflows/golang-binary.yaml" commit-check-linux-amd64.bin
```
