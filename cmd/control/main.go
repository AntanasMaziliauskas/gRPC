package main

import (
	"os"

	"github.com/AntanasMaziliauskas/grpc/control"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {

	control := control.Application{}

	control.Init()

	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:  "dropnode",
			Usage: "Looks for information of specific person on all nodes that are connected",
			Action: func(c *cli.Context) error {
				return control.DropNode(c)
			},
		},
		{
			Name:  "listnodes",
			Usage: "Looks for information of specific person on all nodes that are connected",
			Action: func(c *cli.Context) error {
				return control.ListNodes(c)
			},
		},
		{
			Name:  "getonepersonbroadcast",
			Usage: "Looks for information of specific person on all nodes that are connected",
			Action: func(c *cli.Context) error {
				return control.GetOnePersonBroadcast(c)
			},
		},
		{
			Name:  "getonepersonnode",
			Usage: "Looks for information of specific person on all nodes that are connected",
			Action: func(c *cli.Context) error {
				return control.GetOnePersonNode(c)
			},
		},
		{
			Name:  "getmultipersonbroadcast",
			Usage: "Looks for information of specific person on all nodes that are connected",
			Action: func(c *cli.Context) error {
				return control.GetMultiPersonBroadcast(c)
			},
		},
		{
			Name:  "getmultipersonnode",
			Usage: "Looks for information of specific person on all nodes that are connected",
			Action: func(c *cli.Context) error {
				return control.GetMultiPersonNode(c)
			},
		},
		{
			Name:  "droponepersonbroadcast",
			Usage: "Looks for information of specific person on all nodes that are connected",
			Action: func(c *cli.Context) error {
				return control.DropOnePersonBroadcast(c)
			},
		},
		{
			Name:  "droponepersonnode",
			Usage: "Looks for information of specific person on all nodes that are connected",
			Action: func(c *cli.Context) error {
				return control.DropOnePersonNode(c)
			},
		},
		{
			Name:  "dropmultipersonbroadcast",
			Usage: "Looks for information of specific person on all nodes that are connected",
			Action: func(c *cli.Context) error {
				return control.DropMultiPersonBroadcast(c)
			},
		},
		{
			Name:  "dropmultipersonnode",
			Usage: "Looks for information of specific person on all nodes that are connected",
			Action: func(c *cli.Context) error {
				return control.DropMultiPersonNode(c)
			},
		},
		{
			Name:  "insertonepersonnode",
			Usage: "Looks for information of specific person on all nodes that are connected",
			Action: func(c *cli.Context) error {
				return control.InsertOnePersonNode(c)
			},
		},
		{
			Name:  "insertmultipersonnode",
			Usage: "Looks for information of specific person on all nodes that are connected",
			Action: func(c *cli.Context) error {
				return control.InsertMultiPersonNode(c)
			},
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "node",
			Value: "Node-001",
			Usage: "Node name",
		},
		cli.StringFlag{
			Name:  "person",
			Value: "Jonas",
			Usage: "One or a list of names you want to get info about. For a list use ','",
		},
	}
	/*flag.Usage = func() {
		fmt.Printf("Usage:\n")
		fmt.Printf("   %s\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}*/

	if err := app.Run(os.Args); err != nil {
		logrus.WithError(err).Fatal("Starting application")
	}
}
