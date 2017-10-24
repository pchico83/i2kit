package main

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var taskKillCommand = cli.Command{
	Name:      "kill",
	Usage:     "signal a container (default: SIGTERM)",
	ArgsUsage: "CONTAINER",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "signal, s",
			Value: "SIGTERM",
			Usage: "signal to send to the container",
		},
		cli.IntFlag{
			Name:  "pid",
			Usage: "pid to kill",
			Value: 0,
		},
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "send signal to all processes inside the container",
		},
	},
	Action: func(context *cli.Context) error {
		id := context.Args().First()
		if id == "" {
			return errors.New("container id must be provided")
		}
		signal, err := parseSignal(context.String("signal"))
		if err != nil {
			return err
		}
		var (
			pid = context.Int("pid")
			all = context.Bool("all")
		)
		if pid > 0 && all {
			return errors.New("enter a pid or all; not both")
		}
		client, ctx, cancel, err := newClient(context)
		if err != nil {
			return err
		}
		defer cancel()
		container, err := client.LoadContainer(ctx, id)
		if err != nil {
			return err
		}
		task, err := container.Task(ctx, nil)
		if err != nil {
			return err
		}
		return task.Kill(ctx, signal)
	},
}
