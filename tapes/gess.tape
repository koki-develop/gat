# configuration
Output ./docs/gess.gif
Set Shell "bash"
Set FontSize 32
Set Width 1800
Set Height 600

# setup
Hide
Type "go install" Enter
Type "gat --help" Enter
Sleep 2s
Ctrl+l
Show

# demo
Type "gat --force-color ./README.md | less -R" Sleep 500ms Enter
Sleep 2s
Down@100ms 30
