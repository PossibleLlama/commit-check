Output basic.gif

Require commit-check
Require git

Set Shell "bash"
Set FontSize 32
Set Width 1600
Set Height 1000

# Start app
Type "commit-check"
Sleep 500ms
Enter

# Select type
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 1s
Enter

# Select scope
Sleep 1s
Enter

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

# Display results
# format subject, newline, body
Hide
Type "git log --pretty=format:%s%n%b"
Enter
Show
Sleep 3s
