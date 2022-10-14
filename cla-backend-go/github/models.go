// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package github

// Status Github PR status
type Status struct {
	URL         string        `json:"url,omitempty"`
	AvatarURL   string        `json:"avatar_url,omitempty"`
	ID          int           `json:"id,omitempty"`
	NodeID      string        `json:"node_id,omitempty"`
	State       string        `json:"state,omitempty"`
	Description string        `json:"description,omitempty"`
	TargetURL   string        `json:"target_url,omitempty"`
	Context     string        `json:"context,omitempty"`
	CreatedAt   string        `json:"created_at,omitempty"`
	UpdatedAt   string        `json:"updated_at,omitempty"`
	Creator     StatusCreator `json:"creator,omitempty"`
}

// StatusCreator representative of a Github Status creator
type StatusCreator struct {
	Login             string `json:"login,omitempty"`
	ID                int    `json:"id,omitempty"`
	NodeID            string `json:"node_id,omitempty"`
	AvatarURL         string `json:"avatar_url,omitempty"`
	GravatarID        string `json:"gravatar_id,omitempty"`
	URL               string `json:"url,omitempty"`
	HTMLURL           string `json:"html_url,omitempty"`
	FollowersURL      string `json:"followers_url,omitempty"`
	FollowingURL      string `json:"following_url,omitempty"`
	GistsURL          string `json:"gists_url,omitempty"`
	StarredURL        string `json:"starred_url,omitempty"`
	SubscriptionsURL  string `json:"subscriptions_url,omitempty"`
	OrganizationsURL  string `json:"organizations_url,omitempty"`
	ReposURL          string `json:"repos_url,omitempty"`
	EventsURL         string `json:"events_url,omitempty"`
	ReceivedEventsURL string `json:"received_events_url,omitempty"`
	Type              string `json:"type,omitempty"`
	SiteAdmin         bool   `json:"site_admin,omitempty"`
}

type StatusRequest struct {
	State       *string `json:"state,omitempty"`
	TargetURL   *string `json:"target_url,omitempty"`
	Description *string `json:"description,omitempty"`
	Context     *string `json:"context,omitempty"`
}

type CLABadge struct {
	CLALogoURL     string
	CLALandingPage string
}
