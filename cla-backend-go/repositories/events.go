package repositories

import (
	"fmt"

	"github.com/communitybridge/easycla/cla-backend-go/events"
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/user"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
)

func addGithubRepositoryEvent(eventService events.Service, claUser *user.CLAUser, input *models.GithubRepositoryInput) {
	data := fmt.Sprintf("user [%s] added github repository [%s]", claUser.Name, utils.StringValue(input.RepositoryName))
	eventService.CreateAuditEvent(events.AddGithubRepository,
		claUser,
		utils.StringValue(input.RepositoryProjectID),
		"",
		data,
		true)
}

func deleteGithubRepositoryEvent(eventService events.Service, claUser *user.CLAUser, repositoryName string, projectID string) {
	data := fmt.Sprintf("user [%s] deleted github repository [%s]", claUser.Name, repositoryName)
	eventService.CreateAuditEvent(events.DeleteGithubRepository,
		claUser,
		projectID,
		"",
		data,
		true)
}
