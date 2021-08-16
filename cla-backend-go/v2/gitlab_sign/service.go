package gitlab_sign

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/communitybridge/easycla/cla-backend-go/auth"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"

	"github.com/communitybridge/easycla/cla-backend-go/events"
	"github.com/communitybridge/easycla/cla-backend-go/v2/gitlab_organizations"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/user"
)

type service struct {
	gitlabRepository gitlab_organizations.RepositoryInterface
}

type Service interface {
	GitlabSignRequest(ctx context.Context, req *http.Request, installationID, repositoryID, mergeRequestID, contributorConsoleV2Base string, eventService events.Service) error
}

func NewService(gitlabRepository gitlab_organizations.RepositoryInterface) Service {
	return &service{
		gitlabRepository: gitlabRepository,
	}
}

func (s service) GitlabSignRequest(ctx context.Context,req *http.Request, installationID, repositoryID, mergeRequestID string, eventService events.Service) error {
	
}

func(s service) redirectToConsole(ctx context.Context,req *http.Request,authorizer auth.Authorizer, installationID, repositoryID, mergeRequestID, originURL, contributorBaseURL string, eventService events.Service) error{
	f := logrus.Fields{
		"functionName": "v2.gitlab_sign.service.redirectToConsole",
		"installationID": installationID,
		"repositoryID": repositoryID,
		"mergeRequestID": mergeRequestID,
		"originURL": originURL,
	}

	claUser, err := s.getOrCreateUser(ctx, req, eventService)
	if err != nil {
		msg := fmt.Sprintf("unable to get or create user : %+v ", err)
		log.WithFields(f).Warn(msg)
		return err
	}

	s.gitlabRepository.GetGitlabOrganization()

	consoleURL := fmt.Sprintf("https://%s/#/cla/project/%s/user/%s")
	resp, err := http.Get(consoleURL)

	if err != nil {
		msg := fmt.Sprintf("unable to redirect to : %s , error: %+v ", consoleURL, err)
		log.WithFields(f).Warn(msg)
		return err
	}

}

func(s service) getOrCreateUser(ctx context.Context, r *http.Request, eventsService events.Service) (*user.CLAUser, error) {

	f := logrus.Fields {
		"functionName": "v2.gitlab_sign.service.getOrCreateUser"
	}

	claUser, err := authorizer.SecurityAuth(t[1], []string{})
	if err != nil {
		log.WithFields(f).WithError(err).Warn("parsing failed")
		return nil, err 
	}
	f["claUserName"] = claUser.Name
	f["claUserID"] = claUser.UserID
	f["claUserLFUsername"] = claUser.LFUsername
	f["claUserLFEmail"] = claUser.LFEmail
	f["claUserEmails"] = strings.Join(claUser.Emails, ",")

	// search if user exist in database by username
	userModel, err := usersService.GetUserByLFUserName(claUser.LFUsername)
	if err != nil {
		if err, ok := err.(*utils.UserNotFound); ok {
			log.WithFields(f).Debug("unable to locate user by lf-email")
		} else {
			log.WithFields(f).WithError(err).Warn("searching user by lf-username failed")
			return
		}
	}
	// If found - just return
	if userModel != nil {
		return
	}

	// search if user exist in database by username
	userModel, err = usersService.GetUserByEmail(claUser.LFEmail)
	if err != nil {
		if err, ok := err.(*utils.UserNotFound); ok {
			log.WithFields(f).Debug("unable to locate user by lf-email")
		} else {
			log.WithFields(f).WithError(err).Warn("searching user by lf-email failed")
			return
		}
	}
	// If found - just return
	if userModel != nil {
		return
	}

	// Attempt to create the user
	newUser := &models.User{
		LfEmail:    strfmt.Email(claUser.LFEmail),
		LfUsername: claUser.LFUsername,
		Username:   claUser.Name,
	}
	log.WithFields(f).Debug("creating new user")
	userModel, err = usersService.CreateUser(newUser, nil)
	if err != nil {
		log.WithFields(f).WithField("user", newUser).WithError(err).Warn("creating new user failed")
		return
	}

	// Log the event
	eventsService.LogEvent(&events.LogEventArgs{
		EventType: events.UserCreated,
		UserID:    userModel.UserID,
		UserModel: userModel,
		EventData: &events.UserCreatedEventData{},
	})

}
