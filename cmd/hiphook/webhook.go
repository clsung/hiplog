package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mozillazg/request"
)

var (
	command    = flag.String("command", "query", "command to run on webhook")
	webhookURL = flag.String("webhook_url", "", "webhook url")
)

type CreateWebHook struct {
	URL   string `json:"url"`
	Event string `json:"event"`
	Name  string `json:"name"`
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	hipToken := os.Getenv("HIPLOG_TOKEN")
	if hipToken == "" {
		log.Fatal("Need to set HIPLOG_TOKEN")
	}
	if *webhookURL == "" {
		log.Fatal("Need to specify webhook_url")
	}
	roomID := flag.Arg(0)
	// TODO: make it flag
	c := new(http.Client)
	URL := fmt.Sprintf("https://api.hipchat.com/v2/room/%s/webhook", roomID)
	req := request.NewRequest(c)
	req.Headers = map[string]string{
		"Content-Type": "application/json",
	}
	req.Params = map[string]string{
		"auth_token": hipToken,
	}
	var err error
	var resp *request.Response
	if *command == "create" {
		req.Json = map[string]string{
			"url":   *webhookURL,
			"name":  roomID,
			"event": "room_notification",
		}
		resp, err = req.Post(URL)
	} else if *command == "query" {
		resp, err = req.Get(URL)
	}
	if err != nil {
		defer resp.Body.Close()
		log.Printf("%s failed, err: %v", *command, err)
	} else {
		log.Printf("%s success", *command)
		fmt.Println(resp.Json())
	}
}
