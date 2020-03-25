package token

import (
	"time"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/imroc/req"
)

const (
	tokenPath = "oauth/token" // nolint
)

var (
	clientID      string
	clientSecret  string
	audience      string
	oauthTokenURL string
	token         string
	expiry        time.Time
)

type tokenGen struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Audience     string `json:"audience"`
}

type tokenReturn struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// Init is the token initialization logic
func Init(paramClientID, paramClientSecret, paramAuth0URL, paramAudience string) {
	clientID = paramClientID

	clientSecret = paramClientSecret

	audience = paramAudience

	oauthTokenURL = paramAuth0URL

	if expiry.Year() == 1 {
		expiry = time.Now()
	}
}

func retrieveToken() {
	log.Debugf("Refreshing auth0 token...")

	tg := tokenGen{
		GrantType:    "client_credentials",
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Audience:     audience,
	}

	resp, err := req.Post(oauthTokenURL, req.BodyJSON(&tg))
	if err != nil {
		log.Warnf("refresh token request failed")
		return
	}

	if resp.Response().StatusCode < 200 || resp.Response().StatusCode > 299 {
		log.Fatalf("invalid response from auth0 service %s - received error code: %d, response: %s",
			oauthTokenURL, resp.Response().StatusCode, resp.String())
	}

	var tr tokenReturn
	err = resp.ToJSON(&tr)
	if err != nil {
		log.Warnf("refresh token::json unmarshal failed of response: %s, error: %+v", resp.String(), err)
		return
	}

	//log.Infof("%+v", tr)
	//log.Infof("Token response: %s", resp.String())
	token = tr.TokenType + " " + tr.AccessToken
	if tr.AccessToken == "" || tr.TokenType == "" {
		log.Warnf("Error fetching authentication token - response value is empty.")
	}

	expiry = time.Now().Add(time.Second * time.Duration(tr.ExpiresIn))
}

// GetToken returns the Auth0 Token - in necessary, refreshes the token when expired
func GetToken() string {
	if expiry.Unix()-time.Now().Unix() < 120 {
		retrieveToken()
	}
	return token
}
