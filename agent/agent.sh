#!/bin/sh

function log {
  echo i2kit agent: $1
}

function signal {
  log "sending cfn-signal $1"
  cfn-signal -e $1 --stack $STACK --resource ASG --region $REGION
  [ $? -ne 0 ] && log "cfn-signal failed to be sent."
}

[ -z "$COMPOSE" ] && log "'COMPOSE' must be defined." && signal 1 && return 1
[ -z "$STACK" ] && log "'STACK' must be defined." && signal 1 && return 1
[ -z "$REGION" ] && log "'REGION' must be defined." && signal 1 && return 1

if [ -n "$CONFIG" ]; then
  mkdir /root/.docker
  echo $CONFIG | base64 --decode > /root/.docker/config.json
  [ $? -ne 0 ] && log "config failed to be decoded." && signal 1 && return 1
fi

echo $COMPOSE | base64 --decode > docker-compose.yml
[ $? -ne 0 ] && log "compose failed to be decoded." && signal 1 && return 1

log "deploying compose..."
docker-compose pull --quiet
[ $? -ne 0 ] && log "compose failed to pull images." && signal 1 && return 1
docker-compose up -d
[ $? -ne 0 ] && echo log "compose failed to be deployed." && signal 1 && return 1
log "compose successfully deployed."

signal 0
