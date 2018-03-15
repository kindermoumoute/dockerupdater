package main

import (
	"encoding/json"
	"net/http"

	"docker.io/go-docker"
	"github.com/prometheus/common/log"
)

type server struct {
	cli     *docker.Client
	updates chan string
}

func (s *server) catchWebhooks(w http.ResponseWriter, r *http.Request) {
	var webhook DockerHub
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&webhook)
	if err != nil {
		log.Infoln(err)
		return
	}
	if webhook.PushData.Tag == "latest" {
		s.updates <- webhook.Repository.RepoName + webhook.PushData.Tag
	}
}
