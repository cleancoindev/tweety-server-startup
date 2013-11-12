package reusable

import (
	"bufio"
	"fmt"
	"github.com/garyburd/go-oauth/oauth"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
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

var OauthClient = oauth.Client{
	TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
	ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authenticate",
	TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
	Credentials: oauth.Credentials{
		Token:  "Gqj79MWnbSlx8BZ8wDWYQ",
		Secret: "HyjatdbuRHjc0Y8YaddKEV16RekRnQJZ3M29tqQbs8",
	},
}

func PostTweet(token *oauth.Credentials, msg string) error {
	url_ := "https://api.twitter.com/1.1/statuses/update.json"
	opt := map[string]string{"status": msg}
	param := make(url.Values)
	for k, v := range opt {
		param.Set(k, v)
	}
	OauthClient.SignParam(token, "POST", url_, param)
	res, err := http.PostForm(url_, url.Values(param))
	if err != nil {
		log.Println("failed to post tweet:", err)
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Println("failed to to post:", response_to_string(res), err)
		return err
	}
	return nil
}

func response_to_string(resp *http.Response) string {
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%s", b)
}

func client_auth(requestToken *oauth.Credentials) (*oauth.Credentials, error) {
	cmd := "xdg-open"
	url_ := OauthClient.AuthorizationURL(requestToken, nil)

	args := []string{cmd, url_}
	if runtime.GOOS == "windows" {
		cmd = "rundll32.exe"
		args = []string{cmd, "url.dll,FileProtocolHandler", url_}
	} else if runtime.GOOS == "darwin" {
		cmd = "open"
		args = []string{cmd, url_}
	}
	cmd, err := exec.LookPath(cmd)
	if err != nil {
		log.Fatal("command not found:", err)
	}
	p, err := os.StartProcess(cmd, args, &os.ProcAttr{Dir: "", Files: []*os.File{nil, nil, os.Stderr}})
	if err != nil {
		log.Fatal("failed to start command:", err)
	}
	defer p.Release()

	print("PIN: ")
	stdin := bufio.NewReader(os.Stdin)
	b, err := stdin.ReadBytes('\n')
	if err != nil {
		log.Fatal("canceled")
	}

	if b[len(b)-2] == '\r' {
		b = b[0 : len(b)-2]
	} else {
		b = b[0 : len(b)-1]
	}
	accessToken, _, err := OauthClient.RequestToken(http.DefaultClient, requestToken, string(b))
	if err != nil {
		log.Fatal("failed to request token:", err)
	}
	return accessToken, nil
}

func GetOauthCredentials(config map[string]string) (*oauth.Credentials, bool, error) {
	var result *oauth.Credentials
	var new_credentials bool
	accessToken, foundToken := config["AccessToken"]
	accessSecret, foundSecret := config["AccessSecret"]
	if foundToken && foundSecret {
		result = &oauth.Credentials{accessToken, accessSecret}
		new_credentials = false
	} else {
		requestToken, err := OauthClient.RequestTemporaryCredentials(http.DefaultClient, "", nil)
		if err != nil {
			log.Print("failed to request temporary credentials:", err)
			return nil, false, err
		}
		token, err := client_auth(requestToken)
		if err != nil {
			log.Print("failed to request temporary credentials:", err)
			return nil, false, err
		}

		config["AccessToken"] = token.Token
		config["AccessSecret"] = token.Secret
		result = &oauth.Credentials{token.Token, token.Secret}
		new_credentials = true
	}
	return result, new_credentials, nil
}
