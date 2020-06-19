package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type authDetails struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	OAuthURL     *url.URL
}

// AuthToken is a struct containing the OAuth token details for an authorized user
type AuthToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int32  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type authorizationCodeRequestData struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
}

type discordUserInfoResponse struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Avatar        string `json:"avatar"`
	Discriminator string `json:"discriminator"`
}

func (s *Server) setupAuth() {
	a := &authDetails{
		ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
		ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
		RedirectURI:  os.Getenv("DISCORD_REDIRECT_URI"),
	}

	link, err := url.Parse("https://discord.com")
	if err != nil {
		log.Fatal(err)
	}
	link.Path += "/api/oauth2/authorize"
	params := url.Values{}
	params.Add("client_id", a.ClientID)
	params.Add("redirect_uri", a.RedirectURI)
	params.Add("response_type", "code")
	params.Add("scope", "identify")
	params.Add("prompt", "none")
	link.RawQuery = params.Encode()
	a.OAuthURL = link
	s.authDetails = a
}

func (s *Server) getAuthorizationCode(code string) (*AuthToken, error) {
	link := "https://discord.com/api/oauth2/token"

	// For some reason Discord gives me errors if I pass in the url encoded data generated from
	// url.Values.Encode() so I am just going to concatenate strings like this until I figure out
	// what is causing the problem.
	data := "code=" + code + "&client_id=722724036706041976&client_secret=" + s.authDetails.ClientSecret + "&redirect_uri=http%3A//localhost%3A3000/api/auth/discord/callback&scope=identify&grant_type=authorization_code"
	req, err := http.NewRequest("POST", link, strings.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "PostmanRuntime/7.24.1")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("There was an error authorizing your Discord account")
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	authToken := &AuthToken{}
	fmt.Println(string(body))
	err = json.Unmarshal(body, authToken)
	if err != nil {
		return nil, err
	}

	return authToken, nil
}

func (s *Server) getDiscordUserInfo(accessToken string, refreshToken string) (*discordUserInfoResponse, error) {
	// TO-DO: Add logic to use refresh token if access token is expired

	url := "https://discord.com/api/v6/users/@me"
	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	data := &discordUserInfoResponse{}
	if err := decoder.Decode(data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *Server) isAuthenticated() {

}
