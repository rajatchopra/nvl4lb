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
