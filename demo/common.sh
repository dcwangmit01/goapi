#!/usr/bin/env bash

# This command looks up the name of the binary using the glide.yaml package name
CLI=`cat $(dirname $(dirname $(readlink -f $0)))/glide.yaml | grep "^package" | awk -F/ '{print $NF}'`
