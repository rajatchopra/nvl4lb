package controller

import (
	"net"

	"github.com/sirupsen/logrus"

	kapi "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

var (
	serviceType   reflect.Type = reflect.TypeOf(&kapi.Service{})
)

func (c *controller) serviceHandlers() cache.ResourceEventHandlerFuncs  {
	return cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			svc, ok := obj.(*kapi.Service)
			if !ok {
				logrus.Errorf("Errorneous object type in add service event")
				return
			}
			c.addService(obj.(*kapi.Service))
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			// care about update only if old one didn't have a nodeport, but new one does
			svc, ok := obj.(*kapi.Service)
			if !ok {
				logrus.Errorf("Errorneous object type in add service event")
				return
			}
			c.updateService(obj.(*kapi.Service))
		},
		DeleteFunc: func(obj interface{}) {
			if serviceType != reflect.TypeOf(obj) {
				tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
				if !ok {
					logrus.Errorf("couldn't get object from tombstone: %+v", obj)
					return
				}
				obj = tombstone.Obj
				objType := reflect.TypeOf(obj)
				if serviceType != objType {
					logrus.Errorf("expected tombstone object resource type %v but got %v", serviceType, objType)
					return
				}
			}
			c.deleteService(obj.(*kapi.Service))
		},
	}
}
