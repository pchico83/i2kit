package main

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func main() {
	configBase64 := os.Getenv("CONFIG")
	if configBase64 != "" {
		decoded, err := base64.StdEncoding.DecodeString(configBase64)
		if err != nil {
			log.Fatalf("Config decode error: %s", err.Error())
		}
		if err = os.MkdirAll("/root/.docker", 0744); err != nil {
			log.Fatalf("Mkdir error: %s", err.Error())
		}
		if err = ioutil.WriteFile("/root/.docker/config.json", decoded, 0644); err != nil {
			log.Fatalf("Config write file error: %s", err.Error())
		}
	}
	composeBase64 := os.Getenv("COMPOSE")
	if composeBase64 == "" {
		//TODO: test error output
		log.Fatal("Variable 'COMPOSE' cannot be empty.")
	}
	log.Info("Deploying compose...")
	decoded, err := base64.StdEncoding.DecodeString(composeBase64)
	if err != nil {
		log.Fatalf("Compose decode error: %s", err.Error())
	}
	err = ioutil.WriteFile("docker-compose.yml", decoded, 0644)
	if err != nil {
		log.Fatalf("Compose write file error: %s", err.Error())
	}
	cmd := exec.Command("docker-compose", "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error executing 'docker-compose up -d': %s", err.Error())
	}
	log.Info("Compose deployed successfully!")
}
