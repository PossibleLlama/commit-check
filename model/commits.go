package models

type CommitType string

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
	if c.Description == "" {
		return ""
	}
	prefix := ""
	if c.Type != "" {
		prefix += string(c.Type)
	} else {
		prefix += string(CommitTypeFix)
	}
	if c.Scope != "" {
		prefix += "(" + c.Scope + ")"
	}
	s := prefix + ": " + c.Description
	if c.IsBreakingChange {
		s += "\nBREAKING CHANGE"
	}
	return s
}
