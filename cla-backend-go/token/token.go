package token

import (
	"errors"
	"fmt"
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
	go retrieveToken() //nolint
}

func retrieveToken() error {
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
		return err
	}

	if resp.Response().StatusCode < 200 || resp.Response().StatusCode > 299 {
		err = fmt.Errorf("invalid response from auth0 service %s - received error code: %d, response: %s",
			oauthTokenURL, resp.Response().StatusCode, resp.String())
		log.WithError(err).Warn("invalid response from auth0 service")
		return err
	}

	var tr tokenReturn
	err = resp.ToJSON(&tr)
	if err != nil {
		log.Warnf("refresh token::json unmarshal failed of response: %s, error: %+v", resp.String(), err)
		return err
	}

	//token = tr.TokenType + " " + tr.AccessToken
	token = tr.AccessToken
	if tr.AccessToken == "" || tr.TokenType == "" {
		err = errors.New("error fetching authentication token - response value is empty")
		log.WithError(err).Warn("empty response from auth server")
		return err
	}

	expiry = time.Now().Add(time.Second * time.Duration(tr.ExpiresIn))
	return nil
}

// GetToken returns the Auth0 Token - in necessary, refreshes the token when expired
func GetToken() (string, error) {
	if expiry.Unix()-time.Now().Unix() < 120 {
		err := retrieveToken()
		if err != nil {
			return "", err
		}
	}
	return token, nil
}
