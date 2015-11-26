package service

import (
	_ "github.com/Sirupsen/logrus"

	"encoding/json"
	"net/http"
	"os"
	"time"
)

var (
	VERSION = os.Getenv("VERSION")

	BUILD = os.Getenv("BUILD")
)

type NodeInfo struct {
	Version   string `json:"version"`
	Build     string `json:"build"`
	Timestamp string `json:"current_time"`
}

func Info(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "method not allowed", 403)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(NodeInfo{
		Version:   VERSION,
		Build:     BUILD,
		Timestamp: time.Now().Format(time.RFC3339),
	})
}
