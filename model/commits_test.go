package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommitString(t *testing.T) {
	var tests = []struct {
		name     string
		commit   Commit
		expected string
	}{
		{
			name: "empty",
			commit: Commit{
				Type:             CommitTypeFeat,
				Scope:            "",
				Description:      "",
				IsBreakingChange: false,
			},
			expected: "",
		},
		{
			name: "no scope",
			commit: Commit{
				Type:             CommitTypeFeat,
				Scope:            "",
				Description:      "add a new feature",
				IsBreakingChange: false,
			},
			expected: "feat: add a new feature",
		},
		{
			name: "with scope",
			commit: Commit{
				Type:             CommitTypeFeat,
				Scope:            "scope",
				Description:      "add a new feature",
				IsBreakingChange: false,
			},
			expected: "feat(scope): add a new feature",
		},
		{
			name: "breaking change",
			commit: Commit{
				Type:             CommitTypeFeat,
				Scope:            "",
				Description:      "add a new feature",
				IsBreakingChange: true,
			},
			expected: "feat: add a new feature\nBREAKING CHANGE",
		},
		{
			name: "quit and dryrun make no changes, both true",
			commit: Commit{
				Type:             CommitTypeFeat,
				Scope:            "",
				Description:      "add a new feature",
				IsBreakingChange: true,
				quit:             true,
				dryRun:           true,
			},
			expected: "feat: add a new feature\nBREAKING CHANGE",
		},
		{
			name: "quit and dryrun make no changes, both false",
			commit: Commit{
				Type:             CommitTypeFeat,
				Scope:            "",
				Description:      "add a new feature",
				IsBreakingChange: true,
				quit:             false,
				dryRun:           false,
			},
			expected: "feat: add a new feature\nBREAKING CHANGE",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.commit.String()

			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestCommitIsCommittable(t *testing.T) {
	var tests = []struct {
		name     string
		commit   Commit
		expected bool
	}{
		{
			name: "valid commit",
			commit: Commit{
				Type:        CommitTypeFeat,
				Description: "add a new feature",
			},
			expected: true,
		},
		{
			name: "invalid commit",
			commit: Commit{
				Type: CommitTypeFeat,
			},
			expected: false,
		},
		{
			name: "valid commit with quit set",
			commit: Commit{
				Type:        CommitTypeFeat,
				Description: "add a new feature",
				quit:        true,
			},
			expected: true, // quit does not affect committability, the separate HasQuit() method checks for this
		},
		{
			name: "valid commit with dry run",
			commit: Commit{
				Type:        CommitTypeFeat,
				Description: "add a new feature",
				dryRun:      true,
			},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.commit.IsCommittable()

			assert.Equal(t, test.expected, actual)
		})
	}
}
func TestCommitQuit(t *testing.T) {
	var tests = []struct {
		name     string
		initial  bool
		input    bool
		expected bool
	}{
		{
			name:     "Set quit to true from false",
			initial:  false,
			input:    true,
			expected: true,
		},
		{
			name:     "Set quit to false from true",
			initial:  true,
			input:    false,
			expected: false,
		},
		{
			name:     "Set quit to true from true",
			initial:  true,
			input:    true,
			expected: true,
		},
		{
			name:     "Set quit to false from false",
			initial:  false,
			input:    false,
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			commit := Commit{quit: test.initial}
			commit.Quit(test.input)

			assert.Equal(t, test.expected, commit.quit)
		})
	}
}

func TestCommitDryRun(t *testing.T) {
	var tests = []struct {
		name     string
		initial  bool
		input    bool
		expected bool
	}{
		{
			name:     "Set dryRun to true from false",
			initial:  false,
			input:    true,
			expected: true,
		},
		{
			name:     "Set dryRun to false from true",
			initial:  true,
			input:    false,
			expected: false,
		},
		{
			name:     "Set dryRun to true from true",
			initial:  true,
			input:    true,
			expected: true,
		},
		{
			name:     "Set dryRun to false from false",
			initial:  false,
			input:    false,
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			commit := Commit{dryRun: test.initial}
			commit.DryRun(test.input)

			assert.Equal(t, test.expected, commit.dryRun)
		})
	}
}
