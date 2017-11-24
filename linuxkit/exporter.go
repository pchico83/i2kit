package linuxkit

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/pchico83/i2kit/k8"
	"github.com/pchico83/i2kit/utils"
	log "github.com/sirupsen/logrus"
)

var re utils.CustomRegex = utils.CustomRegex{regexp.MustCompile("Created AMI: (?P<ami>ami-[a-z0-9]+)")}

//Export generates a AWS ami from a linuxkit template
func Export(deployment *k8.Deployment) (string, error) {
	linuxkitPath, err := exec.Command("which", "linuxkit").Output()
	if err != nil {
		return "", fmt.Errorf("Could not find linuxkit in the system (Err: %s). Have you installed it? go get -u github.com/linuxkit/linuxkit/src/cmd/linuxkit", err)
	}
	linuxkitPathTrimmed := strings.Trim(string(linuxkitPath), "\n")

	log.Info(">>> Building specialized linux distribution...")
	template, err := read("./linuxkit/aws.yml")
	if err != nil {
		return "", err
	}
	for _, container := range deployment.Spec.Template.Spec.Containers {
		template.Services = append(
			template.Services,
			&containerYml{
				Name:  container.Name,
				Image: container.Image,
			},
		)
	}

	tempFile, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("i2kit-%s", deployment.Metadata.Name))
	defer os.Remove(tempFile.Name())
	write(template, tempFile.Name())
	outputPath := fmt.Sprintf("%s.raw", tempFile.Name())
	cmd := exec.Command(linuxkitPathTrimmed, "build", "-format", "aws", "-dir", os.TempDir(), tempFile.Name())
	_, err = utils.StreamCommand(cmd)
	if err != nil {
		return "", fmt.Errorf("Error executing 'linuxkit build' command: %s", err)
	}

	log.Info(">>> Pushing image to AWS... (this process will take a while)")
	cmd = exec.Command(linuxkitPathTrimmed, "-v", "push", "aws", "-timeout", "1200", "-bucket", "i2kit", outputPath) // TODO configurable parameter (bucket)
	stdoutStderr, err := utils.StreamCommand(cmd)
	if err != nil {
		return "", fmt.Errorf("Error executing 'linuxkit push' command: %s", err)
	}
	output := re.FindStringSubmatchMap(stdoutStderr.String())
	return output["ami"], nil
}
