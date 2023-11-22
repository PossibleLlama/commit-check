package model

import "strings"

type CommitType string

func (c CommitType) String() string {
	return string(c)
}

const (
	CommitTypeFix      CommitType = "fix"
	CommitTypeFeat     CommitType = "feat"
	CommitTypeDocs     CommitType = "docs"
	CommitTypeStyle    CommitType = "style"
	CommitTypeRefactor CommitType = "refactor"
	CommitTypePerf     CommitType = "perf"
	CommitTypeTest     CommitType = "test"
	CommitTypeChore    CommitType = "chore"
)

var (
	// Based on https://www.conventionalcommits.org/
	// TODO allow config to choose between these lists
	TypeConventionalCommit []CommitType = []CommitType{
		CommitTypeFix,
		CommitTypeFeat,
	}
	// Based on https://github.com/angular/angular.js/blob/master/DEVELOPERS.md#type
	TypeAngular []CommitType = []CommitType{
		CommitTypeFix,
		CommitTypeFeat,
		CommitTypeDocs,
		CommitTypeStyle,
		CommitTypeRefactor,
		CommitTypePerf,
		CommitTypeTest,
		CommitTypeChore,
	}
)

type Commit struct {
	Type             CommitType
	Scope            string
	Description      string
	IsBreakingChange bool
}

func (c Commit) String() string {
	scope := ""
	if c.Scope != "" {
		scope = "(" + c.Scope + ")"
	}
	s := string(c.Type) + scope + ": " + strings.TrimSpace(c.Description)
	if c.IsBreakingChange {
		s += "\n\nBREAKING CHANGE"
	}
	return s
}

func (c Commit) IsValid() bool {
	return c.Type != "" && len(strings.TrimSpace(c.Description)) > 3
}
