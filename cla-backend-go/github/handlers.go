// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"

	"github.com/communitybridge/easycla/cla-backend-go/user"

	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations"
	gh "github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/github"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	ghLib "github.com/google/go-github/github"
	"github.com/savaki/dynastore"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

const (
	// SessionStoreKey is the key used to lookup the session
	SessionStoreKey = "cla-github"
)

// Configure API call
func Configure(api *operations.ClaAPI, clientID, clientSecret, accessToken string, sessionStore *dynastore.Store) {
	oauthConfig := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes: []string{
			"read:org",
		},
		Endpoint: github.Endpoint,
	}

	api.GithubGetOrgHandler = gh.GetOrgHandlerFunc(func(params gh.GetOrgParams, user *user.CLAUser) middleware.Responder {
		if params.OrgName == "" {
			return gh.NewGetOrgBadRequest().WithPayload(&models.ErrorResponse{
				Code:    "400",
				Message: "GithubGetOrgHandler - Missing Organization Name",
			})
		}

		if accessToken == "" {
			return gh.NewGetOrgBadRequest().WithPayload(&models.ErrorResponse{
				Code:    "400",
				Message: "GithubGetOrgHandler - Unable to create oauth2 client for GitHub API requests",
			})
		}

		ctx := context.Background()
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
		tc := oauth2.NewClient(ctx, ts)
		if tc == nil {
			return gh.NewGetOrgBadRequest().WithPayload(&models.ErrorResponse{
				Code:    "400",
				Message: "GithubGetOrgHandler - Unable to create oauth2 client",
			})
		}

		client := ghLib.NewClient(tc)
		if client == nil {
			return gh.NewGetOrgBadRequest().WithPayload(&models.ErrorResponse{
				Code:    "400",
				Message: "GithubGetOrgHandler - Unable to create GitHub client",
			})
		}

		org, resp, err := client.Organizations.Get(ctx, params.OrgName)
		if err != nil {
			log.Warnf("GithubGetOrgHandler - GitHub response error looking up org: %s, error: %+v", params.OrgName, err)
		}

		if resp.Response.StatusCode < 200 || resp.Response.StatusCode > 299 {
			log.Warnf("GithubGetOrgHandler - Non success response code from GitHub: %d while querying for GitHub Org: %s",
				resp.Response.StatusCode, params.OrgName)
			return gh.NewGetOrgNotFound()
		}

		log.Debugf("GithubGetOrgHandler - Success looking up GitHub Organization %s - ID is %d",
			params.OrgName, *org.ID)
		return gh.NewGetOrgOK()
	})

	api.GithubLoginHandler = gh.LoginHandlerFunc(func(params gh.LoginParams) middleware.Responder {
		return middleware.ResponderFunc(
			func(w http.ResponseWriter, pr runtime.Producer) {

				// Get a session. Get() always returns a session, even if empty.
				session, err := sessionStore.Get(params.HTTPRequest, SessionStoreKey)
				if err != nil {
					log.Warnf("Error fetching session store value from key: %s, error: %v", SessionStoreKey, err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				log.Debugf("GH Login Handler loaded the http session (%s): %v", session.Name(), session)

				// Store the callback url so we can redirect back to it once logged in.
				session.Values["callback"] = params.Callback
				//session.Values[""] = params.

				// Generate a csrf token to send
				state, err := uuid.NewV4()
				if err != nil {
					log.Warnf("Error creating new UUIDv4, error: %v", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				session.Values["state"] = state.String()

				err = session.Save(params.HTTPRequest, w)
				if err != nil {
					log.Warnf("Error saving session, error: %v", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				log.Debugf("GH Login handler saved the http session: %v", session)
				log.Debugf("redirecting flow to %s", oauthConfig.AuthCodeURL(state.String()))
				http.Redirect(w, params.HTTPRequest, oauthConfig.AuthCodeURL(state.String()), http.StatusFound)
			})
	})

	api.GithubRedirectHandler = gh.RedirectHandlerFunc(func(params gh.RedirectParams) middleware.Responder {
		return middleware.ResponderFunc(
			func(w http.ResponseWriter, pr runtime.Producer) {
				// Verify csrf token
				session, err := sessionStore.Get(params.HTTPRequest, SessionStoreKey)
				if err != nil {
					log.Warnf("error with session store lookup, error: %v", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				persistedState, ok := session.Values["state"].(string)
				if !ok {
					log.Warnf("Error getting session state, error: %v", err)
					http.Error(w, "no session state", http.StatusInternalServerError)
					return
				}

				if params.State != persistedState {
					msg := fmt.Sprintf("mismatch state, received: %s from callback, but loaded our state as: %s",
						params.State, persistedState)
					log.Warnf(msg)
					http.Error(w, msg, http.StatusInternalServerError)
					return
				}

				// trade temporary code for access token
				token, err := oauthConfig.Exchange(context.TODO(), params.Code)
				if err != nil {
					log.Warnf("unable to exchange oath code, error: %v", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				// persist access token
				session.Values["github_access_token"] = token.AccessToken

				err = session.Save(params.HTTPRequest, w)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				callback, ok := session.Values["callback"].(string)
				if !ok {
					log.Warn("unable to find callback to redirect to")
					http.Error(w, "unable to find callback to redirect to", http.StatusInternalServerError)
					return
				}

				http.Redirect(w, params.HTTPRequest, callback, http.StatusFound)
			})
	})
}
