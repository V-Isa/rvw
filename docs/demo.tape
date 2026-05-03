Output docs/demo.gif

Set Shell "bash"
Set FontSize 16
Set Width 1200
Set Height 760
Set TypingSpeed 20ms
Set Framerate 24
Set Theme "Catppuccin Mocha"

Hide
Type "mkdir -p .cache/demo-bin"
Enter
Type "GOCACHE=$PWD/.cache/go-build go build -o .cache/demo-bin/rvw ./cmd/rvw"
Enter
Type "cp docs/demo-output.sh .cache/demo-bin/terraform"
Enter
Type "chmod +x .cache/demo-bin/terraform"
Enter
Type "export PATH=$PWD/.cache/demo-bin:$PATH"
Enter
Type "export PS1='$ '"
Enter
Type "clear"
Enter
Type "rvw run terraform plan"
Enter
Sleep 300ms
Show

Sleep 250ms
Type@50ms "jjjjjjj"
Sleep 80ms
Type "c"
Sleep 80ms
Type "Block wildcard IAM policy."
Enter
Sleep 220ms
Type@50ms "jjjjjjjjjj"
Sleep 80ms
Type "c"
Sleep 80ms
Type "Confirm DB resize window and rollback plan."
Enter
Sleep 220ms
Type "y"
Sleep 500ms
Type "q"
Hide
Type "rm -rf .cache/demo-bin"
Enter
