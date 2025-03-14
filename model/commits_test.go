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
			expected: "feat: add a new feature\n\nBREAKING CHANGE",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.commit.String()

			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestCommitIsValid(t *testing.T) {
	var tests = []struct {
		name     string
		commit   Commit
		expected bool
	}{
		{
			name: "valid commit",
			commit: Commit{
				Type:             CommitTypeFeat,
				Scope:            "",
				Description:      "add a new feature",
				IsBreakingChange: false,
			},
			expected: true,
		},
		{
			name: "invalid commit",
			commit: Commit{
				Type:             CommitTypeFeat,
				Scope:            "",
				Description:      "",
				IsBreakingChange: false,
			},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.commit.IsValid()

			assert.Equal(t, test.expected, actual)
		})
	}
}
