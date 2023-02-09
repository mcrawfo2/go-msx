#!/usr/bin/env sh

# Copyright © 2023, Cisco Systems Inc.
# Use of this source code is governed by an MIT-style license that can be
# found in the LICENSE file or at https://opensource.org/licenses/MIT.

echo "Making: Skel completions Test Script"

if [ -z "$FIXT" ]; then
  echo "FIXT not set, it should be the fixtures directory path in your go-msx repo"
  exit 1
fi

tmp=$(mktemp -d)
cd "$tmp" || exit 1

echo "Asking skel for fish, bash, zsh and powershell completions"
skel completion fish > fish.completion
skel completion bash > bash.completion
skel completion zsh > zsh.completion
skel completion powershell > powershell.completion

echo "Wrapping completions into a txtar with the stub script"
txtarwrap -prefix=golden - fish.completion bash.completion zsh.completion powershell.completion \
< "$FIXT"/build/completions.stub \
> "$FIXT"/final/completions-test.txtar
echo "Output: $FIXT/final/completions-test.txtar"

rm -rf "$tmp"