#!/usr/bin/env bash

<<EOF
THE_TLDR_INSTRUCTIONS
    # from the project root dir
    ./demo/demo.sh
    # adjust the SLEEP_SHORT and SLEEP_LONG vars below
THE_TLDR_INSTRUCTIONS
EOF

SLEEP_SHORT=2
SLEEP_LONG=5

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
DEMO_PROMPT="${CYAN}$ "

# hide the evidence
clear
