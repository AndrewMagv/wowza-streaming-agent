package main

import (
	api "github.com/andrewmagv/wowza-streaming-agent/api"

	log "github.com/Sirupsen/logrus"

	"encoding/json"
	"errors"
	"net/http"
)

const (
	DefaultInfoURI = "http://ambassador:29091/info"
)

var (
	ErrNodeInfoNotFound = errors.New("key not in info")
)

func runAPIEndpoint(addr string, stop chan<- struct{}) {
	defer close(stop)

	server := api.GetServer()

	server.Addr = addr
	log.Error(server.ListenAndServe())
}

func getNodeInfo() (map[string]string, error) {
	resp, err := http.Get(DefaultInfoURI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var info = make(map[string]string)

	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return nil, err
	} else {
		return info, nil
	}
}
