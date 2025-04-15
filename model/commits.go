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
	dryRun           bool

	quit bool // This is used to quit the program. TODO, find a more elegant way to do this
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
		s += "\nBREAKING CHANGE"
	}
	return s
}

func (c Commit) IsCommittable() bool {
	return c.Type != "" && len(strings.TrimSpace(c.Description)) > 3 && !c.dryRun
}

func (c Commit) IsCommittableReason() []string {
	reasons := []string{}
	if c.Type == "" {
		reasons = append(reasons, "Type is empty")
	}
	if len(strings.TrimSpace(c.Description)) < 3 {
		reasons = append(reasons, "Description is too short")
	}
	if c.dryRun {
		reasons = append(reasons, "Dry run is set")
	}
	return reasons
}

func (c *Commit) Quit(q bool) {
	c.quit = q
}

func (c Commit) HasQuit() bool {
	return c.quit
}

func (c *Commit) DryRun(d bool) {
	c.dryRun = d
}

func (c Commit) IsDryRun() bool {
	return c.dryRun
}
