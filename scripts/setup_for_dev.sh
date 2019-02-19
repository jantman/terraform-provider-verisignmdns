#!/bin/bash

# if you want INVALID credentials to get a 401:
# export VERISIGN_MDNS_API_TOKEN=foobarbaz
export VERISIGN_MDNS_API_TOKEN=324f3919c575bae096c4df3e638a83d6
export VERISIGN_MDNS_API_URL="http://127.0.0.1:5000/"
export VERISIGN_MDNS_DEBUG="true"
export VERISIGN_MDNS_TIMEOUT=5
export VERISIGN_ACCOUNT_ID=9999999
export VERISIGN_ZONE_NAME=example.com
export TF_LOG=DEBUG
