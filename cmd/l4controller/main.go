package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/NVIDIA/nvl4lb/pkg/controller"
)

func main() {
	c := cli.NewApp()
	c.Name = "nvl4controller"
	c.Usage = "run l4controller to start the watcher and programmer for l4 services"
	c.Version = config.Version
	c.Flags = append([]cli.Flag{
		cli.StringFlag{
			Name: "lbserver",
			Value: "127.0.0.1:8080",
			Usage: "IP address and port for reaching out to the lb server for publishing lb entries",
		},
		cli.StringFlag{
			"backend-selector",
			Value: "type=infra",
			Usage: "Comma separated key-value pairs to serve as label selectors for nodes that will serve as backends to external lb. These are the nodes where the traffic will arrive at from the client after being load balanced by lb. The destination port will be the nodeport for that service. It will be assumed that these nodes will have nodeport loadbalancing enabled to further load balance to pod backends.",
	})
	c.Action = func(c *cli.Context) error {
		return runSvcWatcher(c)
	}

	if err := c.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func runSvcWatcher(ctx *cli.Context) error {
	lb := ctx.String("lbserver")
	kconfig, err := rest.InClusterConfig()
	if err != nil {
		logrus.Errorf("Error in getting in-cluster-config: %v", err)
		return err
	}
	kClient, err := kubernetes.NewForConfig(kconfig)
	if err != nil {
		logrus.Errorf("Error creating client from internal kubeconfig: %v", err)
		return err
	}
	err = controller.Start(lb, ctx.String("backend-selector"), kClient)
	if err != nil {
		logrus.Errorf("Could not start controller: %v", err)
		return err
	}
	// wait forever
	select {}
}
