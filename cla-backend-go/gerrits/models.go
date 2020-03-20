package gerrits

import "github.com/communitybridge/easycla/cla-backend-go/gen/models"

type Gerrit struct {
	DateCreated   string `json:"date_created,omitempty"`
	DateModified  string `json:"date_modified,omitempty"`
	GerritID      string `json:"gerrit_id,omitempty"`
	GerritName    string `json:"gerrit_name,omitempty"`
	GerritURL     string `json:"gerrit_url,omitempty"`
	GroupIDCcla   string `json:"group_id_ccla,omitempty"`
	GroupIDIcla   string `json:"group_id_icla,omitempty"`
	GroupNameCcla string `json:"group_name_ccla,omitempty"`
	GroupNameIcla string `json:"group_name_icla,omitempty"`
	ProjectID     string `json:"project_id,omitempty"`
	Version       string `json:"version,omitempty"`
}

func (g *Gerrit) toModel() *models.Gerrit {
	return &models.Gerrit{
		DateCreated:   g.DateCreated,
		DateModified:  g.DateModified,
		GerritID:      g.GerritID,
		GerritName:    g.GerritName,
		GerritURL:     g.GerritURL,
		GroupIDCcla:   g.GroupIDCcla,
		GroupIDIcla:   g.GroupIDIcla,
		GroupNameCcla: g.GroupNameCcla,
		GroupNameIcla: g.GroupNameIcla,
		ProjectID:     g.ProjectID,
		Version:       g.Version,
	}
}
