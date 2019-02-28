package controller

import (
	"net"
)

func (c *controller) getNewLoadBalancerIP() (net.IP, error) {
	return c.ipAllocator.Allocate(), nil
}

func (c *controller) freeLoadBalancerIP(ip net.IP) {
	c.ipAllocator.Release(ip)
}
