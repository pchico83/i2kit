package linuxkit

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"runtime"

	log "github.com/Sirupsen/logrus"
	"github.com/moby/tool/src/moby"
)

// Build image from template with Moby tool
func mobyBuild(mobyTemplate *moby.Moby, deploymentName string) error {
	moby.MobyDir = filepath.Join(os.Getenv("HOME"), ".moby")
	var buf *bytes.Buffer
	var w io.Writer
	buf = new(bytes.Buffer)
	w = buf
	err := moby.Build(*mobyTemplate, w, true, "")
	if err != nil {
		return err
	}
	image := buf.Bytes()
	log.Infof("Create outputs:")
	buildFormats := []string{"raw"} // AWS requires a RAW image
	buildHyperkit := runtime.GOOS == "darwin"
	return moby.Formats(filepath.Join(moby.MobyDir, deploymentName), image, buildFormats, 1024, buildHyperkit)
}

// Push image built by Moby to AWS S3 and then build the associated AMI
func linuxkitPush() (string, error) {
	return "ami", nil
}

//Export generates a AWS ami from a linuxkit template
func Export(mobyTemplate *moby.Moby, deploymentName string) (string, error) {
	err := mobyBuild(mobyTemplate, deploymentName)
	if err != nil {
		return "", err
	}
	ami, err := linuxkitPush()
	if err != nil {
		return "", err
	}
	return ami, nil
}
