package lb

import (
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"

	"github.com/NVIDIA/nvl4lb/pkg/common"
)

func ipvsUpdate(lbInfo *common.LBInfo) ([]byte, error) {
	// ignore any errors
	ipvsDelete(lbInfo)

	var flag string
	switch lbInfo.Protocol {
	case "TCP":
		flag = "-t"
	case "UDP":
		flag = "-u"
	}
	// Create service
	svc := lbInfo.ServiceIP.String() + ":" + fmt.Sprintf("%d", lbInfo.ServicePort)

	cmd := exec.Command("ipvsadm", "-A", flag, svc)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Errorf("ipvsadm fail. cmd: %v", svc)
		return out, fmt.Errorf("Error creating virtual service: %v (%s)", err, string(out))
	}

	// Create backends for the service
	for _, backend := range lbInfo.BackendIPs {
		cmd = exec.Command("ipvsadm", "-a", flag, svc, "-r", backend.String()+":"+fmt.Sprintf("%d", lbInfo.BackendPort), "-m")
		out, err := cmd.CombinedOutput()
		if err != nil {
			// delete the service
			ipvsDelete(lbInfo)
			return out, fmt.Errorf("Error adding real servers to virtual service: %v (%s)", err, string(out))
		}
	}

	// setup additional (virtual) IP on interface
	cmd = exec.Command("ip", "addr", "add", lbInfo.ServiceIP.String(), "dev", VirtInterface)
	out, err := cmd.CombinedOutput()
	// add BGP publish route for the service addr on virt interface
	return out, err
}

func ipvsDelete(lbInfo *common.LBInfo) ([]byte, error) {
	var flag string
	switch lbInfo.Protocol {
	case "TCP":
		flag = "-t"
	case "UDP":
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
