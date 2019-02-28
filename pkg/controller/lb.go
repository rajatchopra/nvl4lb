package controller

import (
	"fmt"
	"net"
	"net/http"
	"bytes"
	"io/ioutil"

	"github.com/sirupsen/logrus"

	"github.com/NVIDIA/nvl4lb/pkg/common"
)

func (c *controller) send(action string, payload string) error {
	url := fmt.Sprintf("http://%s:%d/%s", c.lbIP, c.lbPort, action)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	logrus.Debugf("response Status:", resp.Status)
	logrus.Debugf("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	logrus.Debugf("response Body:", string(body))
	return nil
}

func (c *controller) lbUpdate(port, nodePort int32, protocol string, serviceIP net.IP) {
	payload, err := common.LbPayload(port, nodePort, protocol, serviceIP, c.backendNodes)
	if err != nil {
		logrus.Errorf("Failed to update LB: %v", err)
		return
	}
	err = c.send("update", payload)
	if err != nil {
		logrus.Errorf("Failed to update LB: %v", err)
		return
	}
}

func (c *controller) lbDelete(port, nodePort int32, protocol string, serviceIP net.IP) {
	payload, err := common.LbPayload(port, nodePort, protocol, serviceIP, nil)
	if err != nil {
		logrus.Errorf("Failed to create LB payload: %v", err)
		return
	}
	err = c.send("delete", payload)
	if err != nil {
		logrus.Errorf("Failed to delete LB: %v", err)
		return
	}
}

func (c *controller) lbUpdateAll() {
	payload, err := common.LbPayload(0, 0, "", nil, c.backendNodes)
	if err != nil {
		logrus.Errorf("Failed to create LB payload: %v", err)
		return
	}
	err = c.send("sync", payload)
	if err != nil {
		logrus.Errorf("Failed to sync LB: %v", err)
		return
	}
}
