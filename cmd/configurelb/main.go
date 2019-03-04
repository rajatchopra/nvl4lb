package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/NVIDIA/nvl4lb/pkg/lb"
)

func main() {
	c := cli.NewApp()
	c.Name = "configureLB"
	c.Usage = "run configureLB to start the listener that will program an IPVS lb"
	c.Version = "0.0.1"
	c.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "addr",
			Value: "0.0.0.0:8000",
			Usage: "Address and port to listen upon",
		},
		cli.StringFlag{
			Name:  "interface",
			Value: "eth0",
			Usage: "Interface where virtual IPs should be created",
		},
	}
	c.Action = func(c *cli.Context) error {
		return runConfigureLB(c)
	}

	if err := c.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func runConfigureLB(ctx *cli.Context) error {
	logrus.SetLevel(logrus.Level(5))
	eth := ctx.String("interface")
	lb.VirtInterface = eth
	addr := ctx.String("addr")
	lb.StartServer(addr)
	return nil
}
