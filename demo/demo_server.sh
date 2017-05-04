#!/usr/bin/env bash

# Load the common demo lib
demo_dir="$(dirname "$0")"
. "$demo_dir/demo_magic_common.sh"
. "$demo_dir/common.sh"

#####################################################################
# Start the good stuff

# Kill the server if it is already running
SILENCED=`killall -9 ${CLI} 2>&1 >/dev/null`

pe "${CLI} server"
