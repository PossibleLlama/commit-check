Output dry-run.gif

Require commit-check
Require git

Set Shell "bash"
Set FontSize 32
Set Width 1600
Set Height 1000

# Start app
Type "commit-check --dry-run --type-list angular"
Sleep 500ms
Enter

# Select type
Sleep 500ms
Down
Sleep 1s
Enter

# Select scope
Sleep 500ms
Down
Sleep 1s
Enter
Type "ABC-123"

# Add description
Sleep 500ms
Type "This is a new feature!"
Sleep 500ms
Enter
Sleep 1s
Enter

# Select breaking change confirmation
Sleep 500ms
Down
Sleep 1s
Enter
