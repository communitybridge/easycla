// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package gerrits

// WebLink contains the name and url
type WebLink struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// GerritRepoInfo a simplified gerrit repo information data model
type GerritRepoInfo struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	State       string    `json:"state"`
	WebLinks    []WebLink `json:"web_links"`
}

// GerritRepoListResponse a response model for gerrit repository list
type GerritRepoListResponse struct {
	Repos map[string]GerritRepoInfo `json:"*"`
}
