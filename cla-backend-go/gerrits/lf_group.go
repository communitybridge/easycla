package gerrits

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// constants
const (
	DefaultHTTPTimeout = 10 * time.Second
)

// LFGroup contains access information of lf LDAP group
type LFGroup struct {
	LfBaseURL    string
	ClientID     string
	ClientSecret string
	RefreshToken string
}

// LDAPGroup model
type LDAPGroup struct {
	Title string `json:"title"`
}

func (lfg *LFGroup) getAccessToken() (string, error) {
	requestBody, err := json.Marshal(map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": lfg.RefreshToken,
		"scope":         "manage_groups",
	})
	if err != nil {
		return "", err
	}
	OauthURL := fmt.Sprintf("%s/oauth2/token", lfg.LfBaseURL)
	req, err := http.NewRequest("POST", OauthURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(lfg.ClientID, lfg.ClientSecret)
	req.Header.Add("Content-Type", "application/json")

	client := http.Client{
		Timeout: DefaultHTTPTimeout,
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

// GetGroup returns LF LDAP group
func (lfg *LFGroup) GetGroup(groupID string) (*LDAPGroup, error) {
	accessToken, err := lfg.getAccessToken()
	if err != nil {
		return nil, err
	}
	getGroupURL := fmt.Sprintf("%s/rest/auth0/og/%s", lfg.LfBaseURL, groupID)
	req, err := http.NewRequest("GET", getGroupURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	client := http.Client{
		Timeout: DefaultHTTPTimeout,
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
