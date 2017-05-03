#!/usr/bin/env bash

# Load the common demo lib
demo_dir="$(dirname "$0")"
. "$demo_dir/common.sh"

#####################################################################
# Start the good stuff

# Kill the server if it is already running
SILENCED=`killall -9 goapi 2>&1 >/dev/null`

pe "goapi server"
