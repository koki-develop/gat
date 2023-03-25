# configuration
Output ./docs/demo.gif
Set Shell "bash"
Set FontSize 32
Set Width 1200
Set Height 600

# setup
Hide
Type "go install" Enter
Type "gat --help" Enter
Type "cd ./tapes" Enter
Sleep 2s
Ctrl+l
Show

# demo
Type "gat ./main.go" Sleep 500ms Enter
Sleep 5s
