package controller

import (
	"fmt"
	"net"
	"net/http"
)

func (c *controller) send(action string, payload string) {
	url := fmt.Sprintf("http://%s:%d/%s", c.lbIP, c.lbPort, action)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	logrus.Debugf("response Status:", resp.Status)
	logrus.Debugf("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	logrus.Debugf("response Body:", string(body))
}

func (c *controller) lbUpdate(port, nodePort int32, protocol kapi.Protocol) {
	payload := common.LbPayload(port, nodePort, protocol, c.backendNodes)
	c.send("update", payload)
}

func (c *controller) lbDelete(port, nodePort int32, protocol kapi.Protocol) {
	payload := common.LbPayload(port, nodePort, protocol, nil)
	c.send("delete", payload)
}
