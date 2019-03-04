package controller

import (
	"fmt"
	"net"
	"time"

	"github.com/apcera/util/iprange"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	informerfactory "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

const (
	resyncInterval = 12 * time.Hour
)

type Controller interface {
	Run(c kubernetes.Interface) error
}

type controller struct {
	lbIP         string
	lbPort       string
	backendNodes []net.IP
	selector     *metav1.LabelSelector

	ipAllocator  *iprange.IPRangeAllocator
	svcInformer  cache.SharedIndexInformer
	nodeInformer cache.SharedIndexInformer
	kClient      kubernetes.Interface
}

func Start(lb string, cidr string, backendSelector string, staticBackends []string, kClient kubernetes.Interface) error {
	c, err := New(lb, cidr, backendSelector, staticBackends)
	if err != nil {
		return err
	}
	return c.Run(kClient)
}

func New(lb string, cidr, backendSelector string, staticBackends []string) (Controller, error) {
	lbIP, lbPort, err := net.SplitHostPort(lb)
	if err != nil {
		return nil, fmt.Errorf("Error in parsing lb address: %v", err)
	}

	selector, err := metav1.ParseToLabelSelector(backendSelector)
	if err != nil {
		return nil, fmt.Errorf("Error in parsing label selector: %v", err)
	}

	var backendIPs []net.IP
	for _, sIP := range staticBackends {
		ip := net.ParseIP(sIP)
		if ip != nil {
			backendIPs = append(backendIPs, ip)
		}
	}
	c := &controller{
		lbIP:         lbIP,
		lbPort:       lbPort,
		selector:     selector,
		backendNodes: backendIPs,
	}
	// TODO: fix hardcoded IP range
	ipr, _ := iprange.ParseIPRange(cidr)
	c.ipAllocator = iprange.NewAllocator(ipr)
	return c, nil
}

func (c *controller) Run(kClient kubernetes.Interface) error {
	c.kClient = kClient
	iFactory := informerfactory.NewSharedInformerFactory(kClient, resyncInterval)
	c.svcInformer = iFactory.Core().V1().Services().Informer()
	c.nodeInformer = iFactory.Core().V1().Nodes().Informer()

	// create a dummy stop channel
	stopChan := make(chan struct{})

	iFactory.Start(stopChan)
	res := iFactory.WaitForCacheSync(stopChan)
	for oType, synced := range res {
		if !synced {
			return fmt.Errorf("error in syncing cache for %v informer", oType)
		}
	}
	c.svcInformer.AddEventHandler(c.serviceHandlers())
	c.nodeInformer.AddEventHandler(c.nodeHandlers())
	return nil
}
