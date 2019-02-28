package controller

import (
	"fmt"
	"net"
	"time"

	"github.com/apcera/util/iprange"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	informerfactory "k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

const (
	resyncInterval        = 12 * time.Hour
)

type Controller interface {
	Run(c kubernetes.Interface) error
}

type controller struct {
	lbIP string
	lbPort string
	backendNodes []net.IP
	selector *metav1.LabelSelector

	ipAllocator *iprange.IPRangeAllocator
	svcInformer cache.SharedIndexInformer
	nodeInformer cache.SharedIndexInformer
}

func Start(lb string, backendSelector string, kClient kubernetes.Interface) error {
	c, err := New(lb, backendSelector)
	if err != nil {
		return err
	}
	return c.Run(kClient)
}

func New(lb string, backendSelector string) (Controller, error) {
	lbIP, lbPort, err := net.SplitHostPort(lb)
	if err != nil {
		return nil, fmt.Errorf("Error in parsing lb address: %v", err)
	}

	selector, err := metav1.ParseToLabelSelector(backendSelector)
	if err != nil {
		return nil, fmt.Errorf("Error in parsing label selector: %v", err)
	}

	c := &controller{
		lbIP: lbIP,
		lbPort: lbPort,
		selector: selector,
	}
	// TODO: fix hardcoded IP range
	ipr, _ := iprange.ParseIPRange("10.10.100.1/24")
	c.ipAllocator = iprange.NewAllocator(ipr)
	return c, nil
}

func (c *controller) Run(kClient kubernetes.Interface) error {
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
