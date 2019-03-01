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
			Name:  "listpersonsbroadcast",
			Usage: "Returns information about every person from every Node that is connected to the server",
			Action: func(c *cli.Context) error {
				return control.ListPersonsBroadcast(c)
			},
		},
		{
			Name:  "listpersonsnode",
			Usage: "Return information of all persons from specified Node",
			Action: func(c *cli.Context) error {
				return control.ListPersonsNode(c)
			},
		},
		{
			Name:  "dropnode",
			Usage: "Deletes Node from the server",
			Action: func(c *cli.Context) error {
				return control.DropNode(c)
			},
		},
		{
			Name:  "listnodes",
			Usage: "Lists all connected Nodes that are connected to the server",
			Action: func(c *cli.Context) error {
				return control.ListNodes(c)
			},
		},
		{
			Name:  "getonepersonbroadcast",
			Usage: "Returns information of a specific person on all nodes that are connected",
			Action: func(c *cli.Context) error {
				return control.GetOnePersonBroadcast(c)
			},
		},
		{
			Name:  "getonepersonnode",
			Usage: "Returns information about the person from specified Node",
			Action: func(c *cli.Context) error {
				return control.GetOnePersonNode(c)
			},
		},
		{
			Name:  "getmultipersonbroadcast",
			Usage: "Returns information about multiple persons looking through every Node connected to the server",
			Action: func(c *cli.Context) error {
				return control.GetMultiPersonBroadcast(c)
			},
		},
		{
			Name:  "getmultipersonnode",
			Usage: "Returns information about multiple persons from specified Node",
			Action: func(c *cli.Context) error {
				return control.GetMultiPersonNode(c)
			},
		},
		{
			Name:  "droponepersonbroadcast",
			Usage: "Deletes information about the person from Node that has it",
			Action: func(c *cli.Context) error {
				return control.DropOnePersonBroadcast(c)
			},
		},
		{
			Name:  "droponepersonnode",
			Usage: "Deletes information about the person from specified Node",
			Action: func(c *cli.Context) error {
				return control.DropOnePersonNode(c)
			},
		},
		{
			Name:  "dropmultipersonbroadcast",
			Usage: "Deletes information about multiple persons from Node that has it",
			Action: func(c *cli.Context) error {
				return control.DropMultiPersonBroadcast(c)
			},
		},
		{
			Name:  "dropmultipersonnode",
			Usage: "Deletes information about multiple persons from specified Node",
			Action: func(c *cli.Context) error {
				return control.DropMultiPersonNode(c)
			},
		},
		{
			Name:  "insertonepersonnode",
			Usage: "Adds person to specified Node",
			Action: func(c *cli.Context) error {
				return control.InsertOnePersonNode(c)
			},
		},
		{
			Name:  "insertmultipersonnode",
			Usage: "Adds multiple persons to specified Node",
			Action: func(c *cli.Context) error {
				return control.InsertMultiPersonNode(c)
			},
		},
		{
			Name:  "moveoneperson",
			Usage: "Moves person to specified Node",
			Action: func(c *cli.Context) error {
				return control.MoveOnePerson(c)
			},
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "node",
			Value: "Node-000",
			Usage: "Name of the Node",
		},
		cli.StringFlag{
			Name:  "person",
			Value: "Jonas",
			Usage: "One or a list of names. For a list use ','. If you are inserting person use '.' to specify age and profession",
		},
	}
	if err := app.Run(os.Args); err != nil {
		logrus.WithError(err).Fatal("Starting application")
	}
}
