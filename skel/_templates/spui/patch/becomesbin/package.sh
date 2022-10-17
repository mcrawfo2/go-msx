#!/bin/bash -e

cd build
rm -f skyfall.tar
tar -cvf ../skyfall.tar .
mv ../skyfall.tar .
