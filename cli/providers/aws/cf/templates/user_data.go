package templates

import (
	"encoding/base64"
	"fmt"

	"github.com/pchico83/i2kit/cli/schemas/environment"
)

func userData(containerName, encodedCompose string, e *environment.Environment) string {
	value := fmt.Sprintf(
		`#!/bin/bash

set -e
INSTANCE_ID=$(curl http://169.254.169.254/latest/meta-data/instance-id)
sudo docker run \
	--name %s \
	-e COMPOSE=%s \
	-e CONFIG=%s \
	-e INSTANCE_ID=$INSTANCE_ID \
	-e STACK=%s \
	-e REGION=%s \
	-v /var/run/docker.sock:/var/run/docker.sock \
	--log-driver=awslogs \
	--log-opt awslogs-region=%s \
	--log-opt awslogs-group=i2kit-%s \
	--log-opt tag=i2kit-$INSTANCE_ID \
	okteto/agent:1.0`,
		containerName,
		encodedCompose,
		e.B64DockerConfig(),
		containerName,
		e.Provider.Region,
		e.Provider.Region,
		containerName,
	)
	return base64.StdEncoding.EncodeToString([]byte(value))
}
