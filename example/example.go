package main

import (
	"fmt"

	"github.com/pchico83/i2kit/k8"
	"github.com/pchico83/i2kit/linuxkit"
)

func main() {
	deployment, err := k8.Read("./k8/templates/test.yml")
	if err != nil {
		fmt.Print(err)
	}
	mobyTemplate, err := linuxkit.GetTemplate(deployment)
	if err != nil {
		fmt.Print(err)
	}
	ami, err := linuxkit.Export(mobyTemplate, deployment.GetName())
	if err != nil {
		fmt.Print(err)
	}
	fmt.Print(ami)
}
