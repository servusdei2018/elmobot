#!/bin/bash
source .env
export DISCORD_TOKEN
./elmo --token=$DISCORD_TOKEN
