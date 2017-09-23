package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
)

// CustomRegex allows to retrieve regex groups easier
type CustomRegex struct {
	*regexp.Regexp
}

// FindStringSubmatchMap returns regex matches by group
func (r *CustomRegex) FindStringSubmatchMap(s string) map[string]string {
	captures := make(map[string]string)

	match := r.FindStringSubmatch(s)
	if match == nil {
		return captures
	}

	for i, name := range r.SubexpNames() {
		//
		if i == 0 {
			continue
		}
		captures[name] = match[i]

	}
	return captures
}

// StreamCommand streams command output and then returns it in a combined out+err buffer
// Thanks to Nathan LeClaire's blog post https://nathanleclaire.com/blog/2014/12/29/shelled-out-commands-in-golang/ :)
func StreamCommand(cmd *exec.Cmd) (*bytes.Buffer, error) {
	cmdOutReader, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("Error creating StdoutPipe for Cmd: %s", err)
	}
	cmdErrReader, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("Error creating StdoutPipe for Cmd: %s", err)
	}

	var outErr bytes.Buffer
	scannerOut := bufio.NewScanner(cmdOutReader)
	go func() {
		for scannerOut.Scan() {
			scannedBytes := scannerOut.Bytes()
			outErr.Write(scannedBytes)
			fmt.Printf("linuxkit out | %s\n", string(scannedBytes))
		}
	}()
	scannerErr := bufio.NewScanner(cmdErrReader)
	go func() {
		for scannerErr.Scan() {
			scannedBytes := scannerErr.Bytes()
			outErr.Write(scannedBytes)
			fmt.Printf("linuxkit err | %s\n", string(scannedBytes))
		}
	}()
	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("Error starting Cmd: %s", err)
	}
	err = cmd.Wait()
	if err != nil {
		return nil, fmt.Errorf("Error waiting for Cmd: %s", err)
	}
	return &outErr, nil
}
