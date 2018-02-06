#!/bin/sh

set -e

[ -z "$COMPOSE" ] && echo "'COMPOSE' must be defined" && return 1
[ -z "$STACK" ] && echo "'STACK' must be defined" && return 1
[ -z "$REGION" ] && echo "'REGION' must be defined" && return 1

if [ -n "$CONFIG" ]; then
  mkdir /root/.docker
  echo $CONFIG | base64 --decode > /root/.docker/config.json
fi

echo $COMPOSE | base64 --decode > docker-compose.yml
echo "Deploying compose..."
docker-compose up -d
echo "Compose successfully deployed!"

[ "$SIGNAL" == "NO" ] && return 0

echo "Sending cfn-signal..."
cfn-signal -e 0 --stack $STACK --resource ASG --region $REGION
echo "cfn-signal successfully sent!"
