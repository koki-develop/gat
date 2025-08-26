#!/bin/bash

set -euo pipefail

rm -rf completions
mkdir completions

go build .

./gat completion bash > completions/gat.bash
./gat completion zsh  > completions/gat.zsh
./gat completion fish > completions/gat.fish
