package common

import (
	"net"
	"encoding/json"
)

type LBInfo struct {
	ServiceIP net.IP
	ServicePort int32
	BackendIPs []net.IP
	BackendPort int32
	Protocol string
}

func LbPayload(port int32, nodeport int32, protocol string, serviceIP net.IP, backendIPs []net.IP) (string, error) {
	lb := &LBInfo{
		ServiceIP: serviceIP,
		ServicePort: port,
		BackendPort: nodeport,
		BackendIPs: backendIPs,
		Protocol: protocol,
	}
	bytes, err := json.Marshal(lb)
	return string(bytes), err
}

func UnmarshalPayload(payload string) (*LBInfo, error) {
	lb := &LBInfo{}
	err := json.Unmarshal([]byte(payload), lb)
	return lb, err
}
