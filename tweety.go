package main

import (
	"fmt"
	"github.com/garyburd/go-oauth/oauth"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Tweet struct {
	Text       string
	Identifier string `json:"id_str"`
	Source     string
	CreatedAt  string `json:"created_at"`
	User       struct {
		Name            string
		ScreenName      string `json:"screen_name"`
		FollowersCount  int    `json:"followers_count"`
		ProfileImageURL string `json:"profile_image_url"`
	}
	Place *struct {
		Id       string
		FullName string `json:"full_name"`
	}
	Entities struct {
		HashTags []struct {
			Indices [2]int
			Text    string
		}
		UserMentions []struct {
			Indices    [2]int
			ScreenName string `json:"screen_name"`
		} `json:"user_mentions"`
		Urls []struct {
			Indices [2]int
			Url     string
		}
	}
}

var oauthClient = oauth.Client{
	TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
	ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authenticate",
	TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
}

func postTweet(token *oauth.Credentials, msg string) error {
	url_ := "https://api.twitter.com/1.1/statuses/update.json"
	opt := map[string]string{"status": msg}
	param := make(url.Values)
	for k, v := range opt {
		param.Set(k, v)
	}
	oauthClient.SignParam(token, "POST", url_, param)
	res, err := http.PostForm(url_, url.Values(param))
	if err != nil {
		log.Println("failed to post tweet:", err)
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Println("failed to get timeline:", err)
		return err
	}
	return nil
}

func filter(s []string, fn func(string) bool) []string {
	var p []string // == nil
	for _, v := range s {
		if fn(v) {
			p = append(p, v)
		}
	}
	return p
}

func is_ipv4(in string) bool {
	return !strings.Contains(in, ":")
}

func get_public_ip() string {
	resp, err := http.Get("http://api.externalip.net/ip")
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	public_ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%s", public_ip)
}

func get_hostname_and_ips() (string, string) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	local_ip_addresses, err := net.LookupHost(hostname)
	if err != nil {
		log.Fatal(err)
	}

	public_ip := get_public_ip()

	ip_addresses := fmt.Sprintf("%s, %s", public_ip, strings.Join(filter(local_ip_addresses, is_ipv4), ", "))

	return hostname, ip_addresses
}

func main() {
	hostname, ip_addresses := get_hostname_and_ips()
	msg := fmt.Sprintf("%s just woke up at %s #TweetyServerStartup", hostname, ip_addresses)

	oauthClient.Credentials.Token = "TOKEN"
	oauthClient.Credentials.Secret = "SECRET"
	token := &oauth.Credentials{"SECRET", "TOKEN"}
	postTweet(token, msg)
}
