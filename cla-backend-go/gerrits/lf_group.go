package gerrits

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	HttpTimeout = time.Duration(10 * time.Second)
)

type LFGroup struct {
	LfBaseUrl    string
	ClientID     string
	ClientSecret string
	RefreshToken string
}

type LDAPGroup struct {
	Title string `json:"title"`
}

func (lfg *LFGroup) getAccessToken() (string, error) {
	requestBody, err := json.Marshal(map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": lfg.RefreshToken,
		"scope":         "manage_groups",
	})
	OauthURL := fmt.Sprintf("%s/oauth2/token", lfg.LfBaseUrl)
	req, err := http.NewRequest("POST", OauthURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(lfg.ClientID, lfg.ClientSecret)
	req.Header.Add("Content-Type", "application/json")

	client := http.Client{
		Timeout: HttpTimeout,
	}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var out struct {
		AccessToken string `json:"access_token"`
	}
	err = json.Unmarshal(body, &out)
	if err != nil {
		return "", err
	}
	return out.AccessToken, nil
}

func (lfg *LFGroup) GetGroup(groupID string) (*LDAPGroup, error) {
	accessToken, err := lfg.getAccessToken()
	if err != nil {
		return nil, err
	}
	getGroupUrl := fmt.Sprintf("%s/rest/auth0/og/%s", lfg.LfBaseUrl, groupID)
	req, err := http.NewRequest("GET", getGroupUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	client := http.Client{
		Timeout: HttpTimeout,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var out LDAPGroup
	err = json.Unmarshal(body, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
