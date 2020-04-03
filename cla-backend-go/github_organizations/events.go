package github_organizations

import (
	"fmt"

	"github.com/communitybridge/easycla/cla-backend-go/events"
	"github.com/communitybridge/easycla/cla-backend-go/user"
)

func addGithubOrganizationEvent(eventService events.Service, claUser *user.CLAUser, githubOrgName string) {
	data := fmt.Sprintf("user [%s] added github organization [%s]", claUser.Name, githubOrgName)
	eventService.CreateAuditEvent(events.AddGithubOrganization,
		claUser,
		"",
		"",
		data,
		true)
}

func deleteGithubOrganizationEvent(eventService events.Service, claUser *user.CLAUser, githubOrgName string) {
	data := fmt.Sprintf("user [%s] deleted github organization [%s]", claUser.Name, githubOrgName)
	eventService.CreateAuditEvent(events.DeleteGithubOrganization,
		claUser,
		"",
		"",
		data,
		true)
}
