package github_organizations

import (
	"fmt"

	"github.com/LF-Engineering/lfx-kit/auth"

	"github.com/communitybridge/easycla/cla-backend-go/events"
)

func addGithubOrganizationEvent(eventService events.Service, authUser *auth.User, githubOrgName string) {
	data := fmt.Sprintf("user [%s] added github organization [%s]", authUser.UserName, githubOrgName)
	eventService.CreateAuditEventWithUserID(events.AddGithubOrganization,
		authUser.UserName,
		"",
		"",
		data,
		true)
}

func deleteGithubOrganizationEvent(eventService events.Service, authUser *auth.User, githubOrgName string) {
	data := fmt.Sprintf("user [%s] deleted github organization [%s]", authUser.UserName, githubOrgName)
	eventService.CreateAuditEventWithUserID(events.DeleteGithubOrganization,
		authUser.UserName,
		"",
		"",
		data,
		true)
}
