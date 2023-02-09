#!/usr/bin/env sh

# Copyright © 2023, Cisco Systems Inc.
# Use of this source code is governed by an MIT-style license that can be
# found in the LICENSE file or at https://opensource.org/licenses/MIT.

echo "Testing: Skel Root"

if [ -z "$FIXT" ]; then
  echo "FIXT not set, it should be the fixtures directory path in your go-msx repo"
  exit 1
fi

if [ ! -d "someservice" ]; then
  echo "You should be in a directory where you ran skel (it should have a someservice subdir in it)"
  exit 1
fi

export TXTAR_IGNORE="**/.git/** **/.github/** **/go.sum"
export TXTAR_CONTENTS="**"

txtarwrap . "$FIXT"/final/root-test.txtar \
| testscript

