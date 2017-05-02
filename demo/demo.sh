#!/usr/bin/env bash

# sudo apt-get -yq install ttyrec

<<EOF
THE_TLDR_INSTRUCTIONS
    ./demo.sh  # this script
THE_TLDR_INSTRUCTIONS
EOF

# install demo-magic script in a temporary location
if [ ! -f /tmp/demo-magic.sh ]; then
    curl -fsSL https://raw.githubusercontent.com/paxtonhare/demo-magic/master/demo-magic.sh \
        -o /tmp/demo-magic.sh
fi

silent() {
    "$@" 2>&1 > /dev/null
}

# install pv and jq (needed for parsing output)
if ! silent which pv || ! silent which jq; then
    # pv required by demo-magic
    # jq required for parsing output
    # ttyrec required for recording output
    if [ `uname` = "Darwin" ]; then
	brew install pv jq ttyrec
    else
	sudo apt-get -yq install pv jq ttyrec
    fi
fi


# Configure the options

# speed at which to simulate typing. bigger num = faster
TYPE_SPEED=1000

# include the magic
. /tmp/demo-magic.sh

# custom prompt
# see http://www.tldp.org/HOWTO/Bash-Prompt-HOWTO/bash-prompt-escape-sequences.html for escape sequences
DEMO_PROMPT="${GREEN}➜ ${CYAN}\W "

# hide the evidence
clear

#####################################################################
# Start the good stuff

SLEEP_SHORT=2
SLEEP_LONG=5

# Kill the server if it is already running
SILENCED=`killall -9 goapi 2>&1 >/dev/null`

p  "# Show CLI command usage"
pe "goapi --help"
sleep $SLEEP_LONG
clear

p  "# Start the goapi server"
pe "goapi server &"
sleep $SLEEP_SHORT
clear

p  "# Obtain a JWT auth token as a client"
pe "goapi auth login --help"
sleep $SLEEP_SHORT

p  "TOKEN=$(goapi auth login -t admin password)"
TOKEN=$(goapi auth login -t admin password)
sleep $SLEEP_SHORT
