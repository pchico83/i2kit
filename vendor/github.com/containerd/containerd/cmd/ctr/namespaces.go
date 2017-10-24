package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/log"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var namespacesCommand = cli.Command{
	Name:  "namespaces",
	Usage: "manage namespaces",
	Subcommands: cli.Commands{
		namespacesCreateCommand,
		namespacesSetLabelsCommand,
		namespacesListCommand,
		namespacesRemoveCommand,
	},
}

var namespacesCreateCommand = cli.Command{
	Name:        "create",
	Usage:       "create a new namespace.",
	ArgsUsage:   "[flags] <name> [<key>=<value]",
	Description: "Create a new namespace. It must be unique.",
	Action: func(context *cli.Context) error {
		namespace, labels := objectWithLabelArgs(context)
		if namespace == "" {
			return errors.New("please specify a namespace")
		}
		client, ctx, cancel, err := newClient(context)
		if err != nil {
			return err
		}
		defer cancel()
		namespaces := client.NamespaceService()
		return namespaces.Create(ctx, namespace, labels)
	},
}

var namespacesSetLabelsCommand = cli.Command{
	Name:        "label",
	Usage:       "set and clear labels for a namespace.",
	ArgsUsage:   "[flags] <name> [<key>=<value>, ...]",
	Description: "Set and clear labels for a namespace.",
	Flags:       []cli.Flag{},
	Action: func(context *cli.Context) error {
		namespace, labels := objectWithLabelArgs(context)
		if namespace == "" {
			return errors.New("please specify a namespace")
		}
		client, ctx, cancel, err := newClient(context)
		if err != nil {
			return err
		}
		defer cancel()
		namespaces := client.NamespaceService()
		for k, v := range labels {
			if err := namespaces.SetLabel(ctx, namespace, k, v); err != nil {
				return err
			}
		}
		return nil
	},
}

var namespacesListCommand = cli.Command{
	Name:        "list",
	Aliases:     []string{"ls"},
	Usage:       "list namespaces.",
	ArgsUsage:   "[flags]",
	Description: "List namespaces.",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "print only the namespace name.",
		},
	},
	Action: func(context *cli.Context) error {
		quiet := context.Bool("quiet")
		client, ctx, cancel, err := newClient(context)
		if err != nil {
			return err
		}
		defer cancel()
		namespaces := client.NamespaceService()
		nss, err := namespaces.List(ctx)
		if err != nil {
			return err
		}

		if quiet {
			for _, ns := range nss {
				fmt.Println(ns)
			}
			return nil
		}

		tw := tabwriter.NewWriter(os.Stdout, 1, 8, 1, ' ', 0)
		fmt.Fprintln(tw, "NAME\tLABELS\t")
		for _, ns := range nss {
			labels, err := namespaces.Labels(ctx, ns)
			if err != nil {
				return err
			}

			var labelStrings []string
			for k, v := range labels {
				labelStrings = append(labelStrings, strings.Join([]string{k, v}, "="))
			}
			sort.Strings(labelStrings)

			fmt.Fprintf(tw, "%v\t%v\t\n", ns, strings.Join(labelStrings, ","))
		}
		return tw.Flush()
	},
}

var namespacesRemoveCommand = cli.Command{
	Name:        "remove",
	Aliases:     []string{"rm"},
	Usage:       "remove one or more namespaces",
	ArgsUsage:   "[flags] <name> [<name>, ...]",
	Description: "Remove one or more namespaces. For now, the namespace must be empty.",
	Action: func(context *cli.Context) error {
		var exitErr error
		client, ctx, cancel, err := newClient(context)
		if err != nil {
			return err
		}
		defer cancel()
		namespaces := client.NamespaceService()
		for _, target := range context.Args() {
			if err := namespaces.Delete(ctx, target); err != nil {
				if !errdefs.IsNotFound(err) {
					if exitErr == nil {
						exitErr = errors.Wrapf(err, "unable to delete %v", target)
					}
					log.G(ctx).WithError(err).Errorf("unable to delete %v", target)
					continue
				}

			}

			fmt.Println(target)
		}
		return exitErr
	},
}
