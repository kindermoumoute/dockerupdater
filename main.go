package main

import (
	"log"
	"net/http"

	"os"

	"docker.io/go-docker"
)

var (
	url      = "https://registry-1.docker.io/"
	username = "" // anonymous
	password = "" // anonymous
)

func init() {
	username = os.Getenv("USERNAME")
	password = os.Getenv("PASSWORD")

}

func main() {

	cli, err := docker.NewEnvClient()
	if err != nil {
		panic(err)
	}
	myServer := &server{
		cli,
		make(chan string),
	}
	go myServer.updateContainer()
	http.HandleFunc("/", myServer.catchWebhooks) // set router
	err = http.ListenAndServe(":9090", nil)      // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
