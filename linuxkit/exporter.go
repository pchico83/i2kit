package linuxkit

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/moby/tool/src/moby"
	"github.com/pchico83/i2kit/utils"
)

var re utils.CustomRegex = utils.CustomRegex{regexp.MustCompile("Created AMI: (?P<ami>ami-[a-z0-9]+)")}

// Build image from template with Moby tool
func mobyBuild(mobyTemplate *moby.Moby, deploymentPath string) error {
	log.Info(">>> Moby build started...")
	var buf *bytes.Buffer
	var w io.Writer
	buf = new(bytes.Buffer)
	w = buf
	err := moby.Build(*mobyTemplate, w, true, "")
	if err != nil {
		return err
	}
	image := buf.Bytes()
	log.Infof("Create outputs in %s:", deploymentPath)
	buildFormats := []string{"raw"} // AWS requires a RAW image
	//buildHyperkit := runtime.GOOS == "darwin"
	return moby.Formats(deploymentPath, image, buildFormats, 1024) //, buildHyperkit)
}

// Push image built by Moby to AWS S3 and then build the associated AMI
// (Needs linuxkit installed)
func linuxkitPush(deploymentPath string) (string, error) {
	log.Info(">>> Linuxkit run started... (this process will take a while)")
	linuxkitPath, err := exec.Command("which", "linuxkit").Output()
	if err != nil {
		return "", fmt.Errorf("Could not find linuxkit in the system (Err: %s). Have you installed it? go get -u github.com/linuxkit/linuxkit/src/cmd/linuxkit", err)
	}
	linuxkitPathTrimmed := strings.Trim(string(linuxkitPath), "\n")
	deploymentPathRaw := fmt.Sprintf("%s.raw", deploymentPath)
	cmd := exec.Command(string(linuxkitPathTrimmed), "-v", "push", "aws", "-timeout", "1200", "-bucket", "i2kit", deploymentPathRaw) // TODO configurable parameter (bucket)
	stdoutStderr, err := utils.StreamCommand(cmd)
	if err != nil {
		return "", fmt.Errorf("Error executing linuxkit command: %s", err)
	}
	output := re.FindStringSubmatchMap(stdoutStderr.String())
	return output["ami"], nil
}

//Export generates a AWS ami from a linuxkit template
func Export(mobyTemplate *moby.Moby, deploymentName string) (string, error) {
	moby.MobyDir = filepath.Join(os.Getenv("HOME"), ".moby")
	deploymentPath := filepath.Join(moby.MobyDir, deploymentName)
	err := mobyBuild(mobyTemplate, deploymentPath)
	if err != nil {
		return "", err
	}
	ami, err := linuxkitPush(deploymentPath)
	if err != nil {
		return "", err
	}
	return ami, nil
}
