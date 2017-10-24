package main

import (
	"fmt"
	"io"
	"os"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/log"
	"github.com/urfave/cli"
)

var imagesImportCommand = cli.Command{
	Name:        "import",
	Usage:       "import an image",
	ArgsUsage:   "[flags] <ref> <in>",
	Description: `Import an image from a tar stream.`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "ref-object",
			Value: "",
			Usage: "reference object e.g. tag@digest (default: use the object specified in ref)",
		},
		labelFlag,
	},
	Action: func(context *cli.Context) error {
		var (
			ref       = context.Args().First()
			in        = context.Args().Get(1)
			refObject = context.String("ref-object")
			labels    = labelArgs(context.StringSlice("label"))
		)
		client, ctx, cancel, err := newClient(context)
		if err != nil {
			return err
		}
		defer cancel()
		var r io.ReadCloser
		if in == "-" {
			r = os.Stdin
		} else {
			r, err = os.Open(in)
			if err != nil {
				return err
			}
		}
		img, err := client.Import(ctx,
			ref,
			r,
			containerd.WithRefObject(refObject),
			containerd.WithImportLabels(labels),
		)
		if err != nil {
			return err
		}
		if err = r.Close(); err != nil {
			return err
		}

		log.G(ctx).WithField("image", ref).Debug("unpacking")

		// TODO: Show unpack status
		fmt.Printf("unpacking %s...", img.Target().Digest)
		err = img.Unpack(ctx, context.String("snapshotter"))
		fmt.Println("done")
		return err
	},
}
