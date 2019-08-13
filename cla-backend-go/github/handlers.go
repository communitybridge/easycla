// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations"
	gh "github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/github"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"github.com/savaki/dynastore"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

const (
	// SessionStoreKey is the key used to lookup the session
	SessionStoreKey = "cla-github"
)

// Configure API call
func Configure(api *operations.ClaAPI, clientID, clientSecret string, sessionStore *dynastore.Store) {
	oauthConfig := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes: []string{
			"read:org",
		},
		Endpoint: github.Endpoint,
	}

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
