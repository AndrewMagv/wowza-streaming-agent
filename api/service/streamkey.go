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

const (
	DEFAULT_STREAM_DURATION = 3 * time.Hour
)

var (
	Advertise string

	Node string

	Host string

	rinst *redis.Client // TODO: use ClusterClient?

	exist = errors.New("exist")
)

type StreamKeyAuth struct {
	Username string `json: "username"`
	Password string `json: "password"`
}

type StreamKeyResponse struct {
	Endpoint string
	Node     string
	Key      string
	Auth     StreamKeyAuth
}

func genStreamKey(exp time.Duration) (key string) {
	var ok bool
	for err := exist; err != nil; {
		key = fmt.Sprintf("%x", sha1.Sum([]byte(time.Now().String())))
		ok, err = rinst.SetNX(key, Advertise, exp).Result()
		if !ok {
			err = exist
		}
	}
	return
}

func genUserPass(exp time.Duration) (auth StreamKeyAuth) {
	var (
		user string
		pass string = fmt.Sprintf("%x", sha1.Sum([]byte(time.Now().String())))[:8]

		ok bool
	)
	for err := exist; err != nil; {
		user = fmt.Sprintf("%x", sha1.Sum([]byte(time.Now().String())))[:8]
		ok, err = rinst.SetNX(user, pass, exp).Result()
		if !ok {
			err = exist
		} else {
			auth = StreamKeyAuth{user, pass}
		}
	}
	return
}

func StreamKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", 403)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", 403)
		return
	}

	var (
		enc = json.NewEncoder(w)

		expStr = r.Form.Get("duration")

		expire time.Duration
	)

	if expStr != "" {
		if exp, err := time.ParseDuration(expStr); err != nil {
			http.Error(w, "Bad Request", 403)
			return
		} else {
			expire = exp
		}
	} else {
		expire = DEFAULT_STREAM_DURATION
	}

	resp := &StreamKeyResponse{
		Endpoint: Host,
		Node:     Node,
		Key:      genStreamKey(expire),
		Auth:     genUserPass(expire),
	}

	// send it back to user
	w.Header().Add("Content-Type", "application/json")
	enc.Encode(resp)

	log.WithFields(log.Fields{"key": resp.Key, "node": resp.Node, "host": resp.Endpoint, "origin": Advertise}).Info("req")
}

func init() {
	// FIXME: need to use failover client
	rinst = redis.NewClient(&redis.Options{Addr: "ambassador:6379", DB: 1})
}
