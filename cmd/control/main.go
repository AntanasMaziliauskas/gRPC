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
			Name:  "getonepersonbroadcast",
			Usage: "Looks for information of specific person on all nodes that are connected",
			Action: func(c *cli.Context) error {
				return control.GetOnePersonBroadcast()
			},
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
