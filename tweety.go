package main

import (
	"github.com/thatha/go-service"
	"fmt"
	"github.com/thatha/tweety-server-startup/reusable"
	"log"
	"strings"
)

func main() {
	app_name := "tweety-server-startup"
	desc := "Tweet on Startup"

	var s, err = service.NewService(app_name, desc, desc)

	err = s.Install()
	if err != nil && !strings.HasPrefix(err.Error(), "Init already exists") {
		fmt.Printf("Failed to install: %s\n", err)
		return
	}

	config := reusable.GetConfig(app_name)
	token, should_save_config, err := reusable.GetOauthCredentials(config)
	if err != nil {
		log.Fatal("failed to get access token:", err)
	}
	if should_save_config {
		reusable.SetConfig(app_name, config)
	}

	hostname, ip_addresses := reusable.GetHostnameAndIps()
	msg := fmt.Sprintf("%s just woke up at %s #TweetyServerStartup", hostname, ip_addresses)
	reusable.PostTweet(token, msg)
}
