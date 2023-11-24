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
