package lb

import (
	"fmt"
	"os/exec"

	"github.com/NVIDIA/nvl4lb/pkg/common"
)

func ipvsUpdate(lbInfo *common.LBInfo) ([]byte, error) {
	// ignore any errors
	ipvsDelete(lbInfo)

	var flag string
	switch lbInfo.Protocol {
	case "tcp":
		flag = "-t"
	case "udp":
		flag = "-u"
	}
	// Create service
	svc := lbInfo.ServiceIP.String() + ":" + fmt.Sprintf("%d", lbInfo.ServicePort)

	cmd := exec.Command("ipvsadm", "-A", flag, svc)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return out, fmt.Errorf("Error creating virtual service: %v", err)
	}

	// Create backends for the service
	for _, backend := range lbInfo.BackendIPs {
		cmd = exec.Command("ipvsadm", "-a", flag, svc, "-r", backend.String()+":"+fmt.Sprintf("%d", lbInfo.BackendPort), "-m")
		out, err := cmd.CombinedOutput()
		if err != nil {
			// delete the service
			ipvsDelete(lbInfo)
			return out, err
		}
	}

	// setup additional (virtual) IP on interface
	cmd = exec.Command("ip", "addr", "add", lbInfo.ServiceIP.String(), "dev", VirtInterface)
	return []byte{}, nil
}

func ipvsDelete(lbInfo *common.LBInfo) ([]byte, error) {
	var flag string
	switch lbInfo.Protocol {
	case "tcp":
		flag = "-t"
	case "udp":
		flag = "-u"
	}
	svc := lbInfo.ServiceIP.String() + ":" + fmt.Sprintf("%d", lbInfo.ServicePort)
	cmd := exec.Command("ipvsadm", "-D", flag, svc)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return out, fmt.Errorf("Error deleting virtual service")
	}
	return out, nil
}
