package model

import (
	"fmt"
	"strings"
)

type CommitType string

func (c CommitType) String() string {
	return string(c)
}

const (
	CommitTypeBuild    CommitType = "build"
	CommitTypeChore    CommitType = "chore"
	CommitTypeCI       CommitType = "ci"
	CommitTypeDocs     CommitType = "docs"
	CommitTypeFeat     CommitType = "feat"
	CommitTypeFix      CommitType = "fix"
	CommitTypePerf     CommitType = "perf"
	CommitTypeRefactor CommitType = "refactor"
	CommitTypeStyle    CommitType = "style"
	CommitTypeTest     CommitType = "test"
)

var (
	// Based on https://www.conventionalcommits.org/
	TypeConventionalCommit []CommitType = []CommitType{
		CommitTypeFeat,
		CommitTypeFix,
	}

	// Based on https://github.com/angular/angular/blob/main/contributing-docs/commit-message-guidelines.md
	TypeAngular []CommitType = []CommitType{
		CommitTypeBuild,
		CommitTypeCI,
		CommitTypeDocs,
		CommitTypeFeat,
		CommitTypeFix,
		CommitTypePerf,
		CommitTypeRefactor,
		CommitTypeStyle,
		CommitTypeTest,
	}
)

type Commit struct {
	Type             CommitType
	Scope            string
	Description      string
	IsBreakingChange bool
}

func (c Commit) String() string {
	if strings.TrimSpace(c.Description) == "" {
		return ""
	}
	if c.Scope != "" {
		c.Scope = "(" + c.Scope + ")"
	}
	s := fmt.Sprintf("%s%s: %s", c.Type, c.Scope, strings.TrimSpace(c.Description))
	if c.IsBreakingChange {
		s += "\n\nBREAKING CHANGE"
	}
	return s
}

func (c Commit) IsValid() bool {
	return c.Type != "" && len(strings.TrimSpace(c.Description)) > 3
}
