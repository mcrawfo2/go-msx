#!/usr/bin/env bash

# Copyright © 2023, Cisco Systems Inc.
# Use of this source code is governed by an MIT-style license that can be
# found in the LICENSE file or at https://opensource.org/licenses/MIT.

echo "Making: Skel Root Assets & Test"

if [ -z "$FIXT" ]; then
  echo "FIXT not set, it should be the fixtures directory path in your go-msx repo"
  exit 1
fi
if [ ! -d someservice ]; then
  echo "you should be in a directory where you ran skel and used defaults at the prompts"
  exit 1
fi

export TXTAR_IGNORE="**/.git/**"
export TXTAR_CONTENTS="**"

echo "Making: plain root archive $FIXT/before/plain-root.txtar"
txtarwrap someservice/ > "$FIXT"/before/plain-root.txtar
echo "Output: $FIXT/before/plain-root.txtar"

echo "Making: plain subroot archive (without someservice/ prefix)"
txtarwrap -strip="someservice/" someservice/ > "$FIXT"/before/plain-subroot.txtar
echo "Output: $FIXT/before/plain-subroot.txtar"

echo "Making: golden root archive"
txtarwrap -prefix="golden" "$FIXT"/before/plain-root.txtar \
> "$FIXT"/golden/golden-root.txtar
echo "Output: $FIXT/golden/golden-root.txtar"

echo "Making: test golden root against current"
txtargen golden "$FIXT"/golden/golden-root.txtar \
> "$FIXT"/final/root/root-test.txtar
echo  "Output: $FIXT/final/root-test.txtar"
