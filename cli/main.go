package main

import (
	"io"
	"log"
	"os"

	"github.com/pchico83/i2kit/cli/cmd"
	"github.com/spf13/cobra"
)

func main() {
	commands := &cobra.Command{
		Use:   "i2kit COMMAND [ARG...]",
		Short: "Manage i2kit applications",
	}
	piper, pipew := io.Pipe()
	defer piper.Close()
	defer pipew.Close()
	go func() {
		io.Copy(os.Stdout, piper)
	}()
	commands.AddCommand(
		cmd.Deploy(pipew),
		cmd.Destroy(pipew),
	)

	if err := commands.Execute(); err != nil {
		log.Fatal(err)
	}
}
