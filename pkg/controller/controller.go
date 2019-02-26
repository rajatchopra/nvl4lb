package controller

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	informerfactory "k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

type Controller interface {
	Run(c kubernetes.Interface)
}

type controller struct {
	lbIP string
	lbPort string
	backendNodes []net.IP
	selector *metav1.LabelSelector

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

func New(lb string, backendSelector string, kClient kubernetes.Interface) (Controller, error) {
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
	return c, nil
}

func (c *controller) Run(c kubernetes.Interface) error {
	iFactory := informerfactory.NewSharedInformerFactory(c, resyncInterval)
	c.svcInformer = iFactory.Core().V1().Services().Informer()
	c.nodeInformer = iFactory.Core().V1().Nodes().Informer()

	// create a dummy stop channel
	stopChan := make(chan struct{})

	iFactory.Start(stopChan)
	res := iFactory.WaitForCacheSync(stopChan)
	for oType, synced := range res {
		if !synced {
			return nil, fmt.Errorf("error in syncing cache for %v informer", oType)
		}
	}
	svcInformer.AddEventHandler(c.serviceHandlers())
	nodeInformer.AddEventHandler(c.nodeHandlers())
	return nil
}
