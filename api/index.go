package api

import (
	"github.com/andrewmagv/wowza-streaming-agent/api/service"

	"net/http"
)

func init() {
	mux = http.NewServeMux()
	s = &http.Server{Handler: mux}

	mux := GetServeMux()
	mux.HandleFunc("/info", service.Info)
	mux.HandleFunc("/stream", service.StreamKey)
}
