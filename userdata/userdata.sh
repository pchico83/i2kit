#!/bin/sh

running=false
while [ "$running" = false ]
do
  sleep 1
  curl --unix-socket /var/run/docker.sock http://localhost/containers/json >/dev/null 2>&1
  if [ $? -eq 0 ]; then
    running=true
  fi
done

curl http://169.254.169.254/latest/user-data > script.sh
source script.sh
