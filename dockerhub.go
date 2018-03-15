package main

import (
	"context"

	"strconv"

	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"docker.io/go-docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/prometheus/common/log"
)

type DockerHub struct {
	CallbackURL string `json:"callback_url"`
	PushData    struct {
		Images   []string `json:"images"`
		PushedAt int      `json:"pushed_at"`
		Pusher   string   `json:"pusher"`
		Tag      string   `json:"tag"`
	} `json:"push_data"`
	Repository struct {
		CommentCount    string `json:"comment_count"`
		DateCreated     int    `json:"date_created"`
		Description     string `json:"description"`
		Dockerfile      string `json:"dockerfile"`
		FullDescription string `json:"full_description"`
		IsOfficial      bool   `json:"is_official"`
		IsPrivate       bool   `json:"is_private"`
		IsTrusted       bool   `json:"is_trusted"`
		Name            string `json:"name"`
		Namespace       string `json:"namespace"`
		Owner           string `json:"owner"`
		RepoName        string `json:"repo_name"`
		RepoURL         string `json:"repo_url"`
		StarCount       int    `json:"star_count"`
		Status          string `json:"status"`
	} `json:"repository"`
}

func (s *server) updateContainer() {
	for imageName := range s.updates {

		creds := types.AuthConfig{
			Username:      username,
			Password:      password,
			ServerAddress: url[8:],
		}
		auth, err := s.cli.RegistryLogin(context.Background(), creds)
		if err != nil {
			log.Infoln("Wrong auth", err)
			continue
		}

		s.cli.ImagePull(context.Background(), imageName, types.ImagePullOptions{
			RegistryAuth: auth.IdentityToken,
		})

		containers, err := s.cli.ContainerList(context.Background(), types.ContainerListOptions{})
		if err != nil {
			panic(err)
		}
		var myContainer types.Container
		for _, container := range containers {
			if container.Image == imageName {
				myContainer = container
				break
			}
		}

		// copy config
		nmap := make(nat.PortMap)
		for _, p := range myContainer.Ports {
			nmap[nat.Port(strconv.Itoa(int(p.PrivatePort))+"/"+p.Type)] = []nat.PortBinding{{HostIP: p.IP, HostPort: strconv.Itoa(int(p.PublicPort))}}
		}
		//TODO: volumes

		// kill
		s.cli.ContainerKill(context.Background(), myContainer.ID, "KILL")

		// re-run
		s.cli.ContainerCreate(context.Background(),
			&container.Config{Image: imageName, Env: []string{"EXAMPLE=PLOP"}},
			&container.HostConfig{PortBindings: nmap},
			&network.NetworkingConfig{EndpointsConfig: myContainer.NetworkSettings.Networks}, "")

	}

}
