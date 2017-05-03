#!/usr/bin/env bash

# Clean up from the prior run
SILENCED=`killall -9 goapi 2>&1 > /dev/null`
SILENCED=`if [[ -f goapi.db ]]; then mv goapi.db goapi.db.bak; fi`

func_to_fork() {
    sleep 10  # synchronize with the client script which is started at the same time
    tmux split-window -v -p 25 -t test "yes |./demo/demo_server.sh"
}

func_to_fork &

tmux new-session -n goapi -s test "tmux set-option -t test status off; yes |./demo/demo_client.sh; tmux kill-session -t test"
