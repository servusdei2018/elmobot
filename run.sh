#!/bin/bash
source ./.env
export DISCORD_TOKEN
nohup ./elmo --token=$DISCORD_TOKEN &
