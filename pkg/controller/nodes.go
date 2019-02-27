package controller

import (
	"net"

	"github.com/sirupsen/logrus"

	kapi "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

var (
	nodeType reflect.Type = reflect.TypeOf(&kapi.Node{})
)

func (c *controller) nodeHandlers() cache.ResourceEventHandlerFuncs  {
	return cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			node, ok := obj.(*kapi.Node)
			if !ok {
				logrus.Errorf("Errorneous object type in add node event")
				return
			}
			c.addNode(obj.(*kapi.Node))
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			svc, ok := obj.(*kapi.Node)
			if !ok {
				logrus.Errorf("Errorneous object type in update node event")
				return
			}
			c.updateNode(obj.(*kapi.Node))
		},
		DeleteFunc: func(obj interface{}) {
			if nodeType != reflect.TypeOf(obj) {
				tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
				if !ok {
					logrus.Errorf("couldn't get object from tombstone: %+v", obj)
					return
				}
				obj = tombstone.Obj
				objType := reflect.TypeOf(obj)
				if nodeType != objType {
					logrus.Errorf("expected tombstone object resource type %v but got %v", nodeType, objType)
					return
				}
			}
			c.deleteNode(obj.(*kapi.Node))
		},
	}
}

func getNodeIP(node *kapi.Node) net.IP {
	// TODO: cleanup method to obtain node's IP (internalIP likely, but could be externalIP for some data centers)
	nodeAddr := node.Status.Addresses[0].Address
	ip := net.ParseIP(nodeAddr)
	if ip == nil {
		// error, is it a resolvable string?
		ipv4, ipv6, err := net.LookupIP(nodeAddr)
		if err != nil {
			logrus.Errorf("Error in calculating Node '%s''s IP address: %v", node.Name, err)
			return nil
		}
		if ipv4 != nil {
			ip = ipv4
		} else {
			ip = ipv6
		}
	}
	return ip
}

func (c *controller) addNode(node *kapi.Node) {
	if ip := getNodeIP(node); ip != nil {
		c.backendNodes = append(c.backendNodes, ip)
		c.syncBackends()
	}
}

func (c* controller) updateNode(node *kapi.Node) {
	// ignore for now
	// TODO: check if the update is because the label is not applicable anymore
}

func (c *controller) deleteNode(node *kapi.Node) {
	if nodeIP := getNodeIP(node); nodeIP != nil {
		for i, ip := range c.backendNodes {
			if ip.String()==nodeIP.String() {
				c.backendNodes = append(c.backendNodes[:i], c.backendNodes[i+1:]...)
				c.syncBackends()
			}
		}
	}
}

func (c *controller) syncBackends() {
	// contact lb and update all backends to new list
}
