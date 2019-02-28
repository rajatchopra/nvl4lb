package lb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/NVIDIA/nvl4lb/pkg/common"
)

func StartServer(addr string) {
	router := mux.NewRouter()
	s := &server{
		Server: http.Server{
			Handler: router,
			Addr:    addr,
		},
	}
	router.NotFoundHandler = http.HandlerFunc(http.NotFound)
	router.HandleFunc("/update", s.handleUpdate).Methods("POST")
	router.HandleFunc("/sync", s.handleSync).Methods("POST")
	router.HandleFunc("/delete", s.handleDelete).Methods("DELETE")

	fmt.Println(s.ListenAndServe())
}

func (s *server) requestToLbInfo(r *http.Request) (*common.LBInfo, error) {
	lbInfo := &common.LBInfo{}
	b, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(b, lbInfo); err != nil {
		return nil, fmt.Errorf("JSON unmarshal error: %v", err)
	}
	return lbInfo, nil
}

func (s *server) handleUpdate(w http.ResponseWriter, r *http.Request) {
	lbInfo, err := s.requestToLbInfo(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}
	res, err := ipvsUpdate(lbInfo)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
	} else {
		// Empty response JSON means success with no body
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(res); err != nil {
			logrus.Warningf("Error writing HTTP response: %v", err)
		}
	}
}

func (s *server) handleSync(w http.ResponseWriter, r *http.Request) {
	// TODO: refresh all virtual servers with a new set of backend real servers
}

func (s *server) handleDelete(w http.ResponseWriter, r *http.Request) {
	lbInfo, err := s.requestToLbInfo(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}
	res, err := ipvsDelete(lbInfo)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
	} else {
		// Empty response JSON means success with no body
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(res); err != nil {
			logrus.Warningf("Error writing HTTP response: %v", err)
		}
	}
}
