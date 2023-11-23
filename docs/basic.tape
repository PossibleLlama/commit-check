Output basic.gif

Require commit-check

Set Shell "bash"
Set FontSize 32
Set Width 1600
Set Height 1000

# Start app
Type "commit-check"
Sleep 500ms
Enter

# Select type
Sleep 1s
Down
Sleep 500ms
Down
Sleep 1s
Enter

# Select scope
Sleep 2s
Down
Sleep 1s
Enter

# Add description
Sleep 2s
Type "This is a new feature!"
Sleep 1s
Enter
Sleep 1s
Enter

# Select breaking change confirmation
Sleep 2s
Enter

Sleep 2s
