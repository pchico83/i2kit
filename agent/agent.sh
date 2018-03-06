#!/bin/sh

function log {
  echo "i2kit agent: " $1
}

function signal {
  cfn-signal -e $1 --stack $STACK --resource ASG --region $REGION
  [ $? -ne 0 ] && log "cfn-signal failed to be sent."
  return $1
}

[ -z "$COMPOSE" ] && log "'COMPOSE' must be defined." && return signal 1
[ -z "$STACK" ] && log "'STACK' must be defined." && return signal 1
[ -z "$REGION" ] && log "'REGION' must be defined." && return signal 1

if [ -n "$CONFIG" ]; then
  mkdir /root/.docker
  echo $CONFIG | base64 --decode > /root/.docker/config.json
  [ $? -ne 0 ] && log "config failed to be decoded." && return signal 1
fi

echo $COMPOSE | base64 --decode > docker-compose.yml
[ $? -ne 0 ] && log "compose failed to be decoded." && return signal 1

log "deploying compose..."
docker-compose pull --quiet
[ $? -ne 0 ] && log "compose failed to pull images." && return signal 1
docker-compose up -d
[ $? -ne 0 ] && echo log "compose failed to be deployed." && return signal 1
log "compose successfully deployed."

return signal 0
