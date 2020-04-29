package company

import (
	"sort"
	"sync"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/aws/aws-sdk-go/aws"
	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
	v1SignatureParams "github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/signatures"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/signatures"
	v2UserService "github.com/communitybridge/easycla/cla-backend-go/v2/user-service"
	v2UserServiceModels "github.com/communitybridge/easycla/cla-backend-go/v2/user-service/models"
)

// Service functions for company
type Service interface {
	GetCompanyCLAManagers(companyID string) (*models.CompanyClaManagers, error)
}

// ProjectRepo contains project repo methods
type ProjectRepo interface {
	GetProjectByID(projectID string) (*v1Models.Project, error)
}

type service struct {
	signatureRepo signatures.SignatureRepository
	projectRepo   ProjectRepo
}

// NewService returns instance of company service
func NewService(sigRepo signatures.SignatureRepository, projectRepo ProjectRepo) Service {
	return &service{
		signatureRepo: sigRepo,
		projectRepo:   projectRepo,
	}
}

func (s *service) getAllCCLASignatures(companyID string) ([]v1Models.Signature, error) {
	var sigs []v1Models.Signature
	var lastScannedKey *string
	for {
		signatures, err := s.signatureRepo.GetCompanySignatures(v1SignatureParams.GetCompanySignaturesParams{
			CompanyID:     companyID,
			SignatureType: aws.String("ccla"),
			NextKey:       lastScannedKey,
		}, 1000)
		if err != nil {
			return nil, err
		}
		sigs = append(sigs, signatures.Signatures...)
		if signatures.LastKeyScanned == "" {
			break
		}
		lastScannedKey = aws.String(signatures.LastKeyScanned)
	}
	return sigs, nil
}

func (s *service) GetCompanyCLAManagers(companyID string) (*models.CompanyClaManagers, error) {
	sigs, err := s.getAllCCLASignatures(companyID)
	if err != nil {
		return nil, err
	}
	var claManagers []*models.CompanyClaManager
	lfUsernames := utils.NewStringSet()
	projectIDs := utils.NewStringSet()
	// Get CLA managers
	for _, sig := range sigs {
		for _, user := range sig.SignatureACL {
			claManagers = append(claManagers, &models.CompanyClaManager{
				// DB doesn't have approved_on value
				ApprovedOn: "",
				LfUsername: user.LfUsername,
				ProjectID:  sig.ProjectID,
			})
			lfUsernames.Add(user.LfUsername)
			projectIDs.Add(sig.ProjectID)
		}
	}
	// get userinfo and project info
	var usermap map[string]*v2UserServiceModels.User
	var projects map[string]*v1Models.Project
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		usermap, err = getUsersInfo(lfUsernames.List())
	}()
	go func() {
		defer wg.Done()
		projects = s.getProjects(projectIDs.List())
	}()
	wg.Wait()
	if err != nil {
		return nil, err
	}
	// fill user info
	fillUsersInfo(claManagers, usermap)
	// fill project info
	fillProjectInfo(claManagers, projects)
	// sort result by cla manager name
	sort.Slice(claManagers, func(i, j int) bool {
		return claManagers[i].Name < claManagers[j].Name
	})
	return &models.CompanyClaManagers{List: claManagers}, nil
}

func getUsersInfo(lfUsernames []string) (map[string]*v2UserServiceModels.User, error) {
	userServiceClient := v2UserService.GetClient()
	users, err := userServiceClient.GetUsersByUsernames(lfUsernames)
	if err != nil {
		return nil, err
	}
	usermap := make(map[string]*v2UserServiceModels.User)
	for _, user := range users {
		usermap[user.Username] = user
	}
	return usermap, nil
}

func fillUsersInfo(claManagers []*models.CompanyClaManager, usermap map[string]*v2UserServiceModels.User) {
	for _, cm := range claManagers {
		user, ok := usermap[cm.LfUsername]
		if !ok {
			logging.Warnf("Unable to get user with username %s", cm.LfUsername)
			continue
		}
		cm.Name = user.Name
		cm.LogoURL = user.LogoURL
		cm.UserSfid = user.ID
		if user.Email != nil {
			cm.Email = *user.Email
		} else {
			if len(user.Emails) > 0 {
				cm.Email = utils.StringValue(user.Emails[0].EmailAddress)
			}
		}
	}
}

func (s *service) getProjects(projectIDs []string) map[string]*v1Models.Project {
	projects := make(map[string]*v1Models.Project)
	prChan := make(chan *v1Models.Project)
	for _, id := range projectIDs {
		go func(projectID string) {
			project, err := s.projectRepo.GetProjectByID(projectID)
			if err != nil {
				logging.Warnf("Unable to fetch project details for project %s. error = %s", projectID, err)
			}
			prChan <- project
		}(id)
	}
	for range projectIDs {
		project := <-prChan
		if project != nil {
			projects[project.ProjectID] = project
		}
	}
	return projects
}

func fillProjectInfo(claManagers []*models.CompanyClaManager, projects map[string]*v1Models.Project) {
	for _, claManager := range claManagers {
		project, ok := projects[claManager.ProjectID]
		if !ok {
			continue
		}
		claManager.ProjectName = project.ProjectName
		claManager.ProjectSfid = project.ProjectExternalID
	}
}
