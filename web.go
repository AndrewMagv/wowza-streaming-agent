package main

import (
	api "github.com/andrewmagv/wowza-streaming-agent/api"

	log "github.com/Sirupsen/logrus"

	"encoding/json"
	"errors"
	"net/http"
)

const (
	DefaultInfoURI = "http://localhost:29091/info"
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

func getNodeInfo(k string) (string, error) {
	resp, err := http.Get(DefaultInfoURI)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var info = make(map[string]string)

	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return "", err
	} else if v, ok := info[k]; !ok {
		return "", ErrNodeInfoNotFound
	} else {
		return v, nil
	}
}
