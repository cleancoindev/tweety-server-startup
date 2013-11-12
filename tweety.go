package main

import (
	"fmt"
	"github.com/thatha/tweety-server-startup/reusable"
	"log"
)

func main() {
	app_name := "tweety-server-startup"
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
