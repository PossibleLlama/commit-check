package models

var (
	// Based on https://www.conventionalcommits.org/
	// TODO allow config to choose between these lists
	TypeConventionalCommit []string = []string{
		"fix",
		"feat",
		"BREAKING CHANGE",
	}
	// Based on https://github.com/angular/angular.js/blob/master/DEVELOPERS.md#type
	TypeAngular []string = []string{
		"feat",
		"fix",
		"docs",
		"style",
		"refactor",
		"perf",
		"test",
		"chore",
		"revert",
		"BREAKING CHANGE",
	}
)
