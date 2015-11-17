package service

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/redis.v3"

	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	Advertise string

	rinst *redis.Client // TODO: use ClusterClient?

	exist = errors.New("exist")
)

type StreamKeyResponse struct {
	Endpoint string
	Key      string
}

func genSteamkey() (key string) {
	var ok bool
	for err := exist; err != nil; {
		key = fmt.Sprintf("%x", sha1.Sum([]byte(time.Now().String())))
		ok, err = rinst.SetNX(key, Advertise, 0).Result()
		if !ok {
			err = exist
		}
	}
	return
}

func StreamKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", 403)
		return
	}

	var (
		enc = json.NewEncoder(w)

		streamkey = genSteamkey()
	)

	// send it back to user
	w.Header().Add("Content-Type", "application/json")
	enc.Encode(&StreamKeyResponse{
		Endpoint: Advertise,
		Key:      streamkey,
	})

	log.WithFields(log.Fields{"key": streamkey, "origin": Advertise}).Info("req")
}

func init() {
	// FIXME: need to use failover client
	rinst = redis.NewClient(&redis.Options{Addr: "localhost:6379", DB: 1})
}
