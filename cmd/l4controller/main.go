package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/NVIDIA/nvl4lb/pkg/controller"
)

func main() {
	c := cli.NewApp()
	c.Name = "nvl4controller"
	c.Usage = "run l4controller to start the watcher and programmer for l4 services"
	c.Version = "0.0.1"
	c.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "lbserver",
			Value: "127.0.0.1:8080",
			Usage: "IP address and port for reaching out to the lb server for publishing lb entries",
		},
		cli.StringFlag{
			Name:  "lbCidr",
			Value: "10.162.151.1/24",
			Usage: "CIDR value from which the load balancer IPs will be chosen and assigned to the serviceIP. The external load balancer will be configured to receive traffic at the chosen IP as well. The upstream ToR to the load balancer should route this CIDR to the machine(s) where external lb is running.",
		},
		cli.StringFlag{
			Name:  "backend-selector",
			Value: "",
			Usage: "Comma separated key-value pairs to serve as label selectors for nodes that will serve as backends to external lb. These are the nodes where the traffic will arrive at from the client after being load balanced by lb. The destination port will be the nodeport for that service. It will be assumed that these nodes will have nodeport loadbalancing enabled to further load balance to pod backends.",
		},
		cli.StringFlag{
			Name:  "kubeconfig",
			Value: "",
			Usage: "Absolute path to the kubeconfig file for watching resources to the kubernetes cluster",
		},
	}
	c.Action = func(c *cli.Context) error {
		return runSvcWatcher(c)
	}

	if err := c.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func runSvcWatcher(ctx *cli.Context) error {
	lb := ctx.String("lbserver")

	var err error
	var kconfig *rest.Config
	kubeconfig := ctx.String("kubeconfig")
	if kubeconfig != "" {
		kconfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		kconfig, err = rest.InClusterConfig()
		if err != nil {
			logrus.Errorf("Error in getting in-cluster-config: %v", err)
			return err
		}
	}
	kClient, err := kubernetes.NewForConfig(kconfig)
	if err != nil {
		logrus.Errorf("Error creating client from internal kubeconfig: %v", err)
		return err
	}
	err = controller.Start(lb, ctx.String("lbCidr"), ctx.String("backend-selector"), kClient)
	if err != nil {
		logrus.Errorf("Could not start controller: %v", err)
		return err
	}
	// wait forever
	select {}
}
