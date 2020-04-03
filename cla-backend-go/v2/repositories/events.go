package repositories

import (
	"fmt"

	"github.com/LF-Engineering/lfx-kit/auth"

	"github.com/communitybridge/easycla/cla-backend-go/events"
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
)

func addGithubRepositoryEvent(eventService events.Service, authUser *auth.User, input *models.GithubRepositoryInput) {
	data := fmt.Sprintf("user [%s] added github repository [%s]", authUser.UserName, utils.StringValue(input.RepositoryName))
	eventService.CreateAuditEventWithUserID(events.AddGithubRepository,
		authUser.UserName,
		utils.StringValue(input.RepositoryProjectID),
		"",
		data,
		true)
}

func deleteGithubRepositoryEvent(eventService events.Service, authUser *auth.User, repositoryName string, projectID string) {
	data := fmt.Sprintf("user [%s] deleted github repository [%s]", authUser.UserName, repositoryName)
	eventService.CreateAuditEventWithUserID(events.DeleteGithubRepository,
		authUser.UserName,
		projectID,
		"",
		data,
		true)
}
