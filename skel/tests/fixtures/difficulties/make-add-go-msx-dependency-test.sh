#!/usr/bin/env sh

echo "Making: Skel add-go-msx-dependency Test Script"

TXTAR_IGNORE="**/.git/** **/.github/**"

txtarwrap - "$fixt"/golden/plain-root.txtar \
< "$fixt"/build/add-go-msx-dependency.stub \
> "$fixt"/final/add-go-msx-dependency-test.txtar
