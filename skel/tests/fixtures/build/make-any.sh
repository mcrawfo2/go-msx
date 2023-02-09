#!/usr/bin/env zsh

# Copyright © 2023, Cisco Systems Inc.
# Use of this source code is governed by an MIT-style license that can be
# found in the LICENSE file or at https://opensource.org/licenses/MIT.

# This script builds a txtar self-contained test file for a skel command

# The single parameter should be the test name (e.g. generate-app)
# The test name should be one of the keys in skel/tests/fixtures/build/tests.sh
# It will be translated into a skel command to be tested using the list there

# This script requires that the FIXT environment variable is set to the
# fixtures directory in your go-msx repo

name="$1"

if [ -z "$FIXT" ]; then
  echo "FIXT not set, it should be the fixtures directory path in your go-msx repo"
  exit 1
fi

. "$FIXT"/build/tests.sh

# shellcheck disable=SC2154
torun="${tests[$name]}"
if [ -z "$torun" ]; then
  echo "$name is not a test name I recognise, sorry"
  exit 1
fi

echo "Making: Skel test $name"

echo "Making: ${name}-test txtar"
echo "Will run: skel ${torun}"

txtarwrap - "$FIXT"/before/plain-subroot.txtar \
  < "$FIXT"/build/make-any.stub \
  | testscript -e FIXT="$FIXT" -e TEST_NAME="${name}" -e TEST_CMD="${torun}"
