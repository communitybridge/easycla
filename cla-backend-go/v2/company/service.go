package company

import (
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

type Service interface {
	GetCompanyCLAManagers(companyID string) (*models.CompanyClaManagers, error)
}

type ProjectRepo interface {
	GetProjectByID(projectID string) (*v1Models.Project, error)
}

type service struct {
	signatureRepo signatures.SignatureRepository
	projectRepo   ProjectRepo
}

func NewService(sigRepo signatures.SignatureRepository, projectRepo ProjectRepo) Service {
	return &service{
		signatureRepo: sigRepo,
		projectRepo:   projectRepo,
	}
}

func (s service) getAllCCLASignatures(companyID string) ([]v1Models.Signature, error) {
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

func (s service) GetCompanyCLAManagers(companyID string) (*models.CompanyClaManagers, error) {
	sigs, err := s.getAllCCLASignatures(companyID)
	if err != nil {
		return nil, err
	}
	var claManagers []*models.CompanyClaManager
	for _, sig := range sigs {
		for _, user := range sig.SignatureACL {
			claManagers = append(claManagers, &models.CompanyClaManager{
				ApprovedOn: "",
				LfUsername: user.LfUsername,
				ProjectID:  sig.ProjectID,
			})
		}
	}
	err = fillUserInfo(claManagers)
	if err != nil {
		return nil, err
	}
	err = s.fillProjectInfo(claManagers)
	if err != nil {
		return nil, err
	}
	return &models.CompanyClaManagers{List: claManagers}, nil
}

func fillUserInfo(claManagers []*models.CompanyClaManager) error {
	var lfUserNames []string
	for _, cm := range claManagers {
		lfUserNames = append(lfUserNames, cm.LfUsername)
	}
	userServiceClient := v2UserService.GetClient()
	users, err := userServiceClient.GetUsersByUsernames(lfUserNames)
	if err != nil {
		return err
	}
	usermap := make(map[string]*v2UserServiceModels.User)
	for _, user := range users {
		usermap[user.Username] = user
	}
	for _, cm := range claManagers {
		user, ok := usermap[cm.LfUsername]
		if !ok {
			logging.Warnf("Unable to get user with username %s", cm.LfUsername)
			continue
		}
		cm.Name = user.Name
		cm.LogoURL = user.LogoURL
		cm.UserSFID = user.ID
		if user.Email != nil {
			cm.Email = *user.Email
		} else {
			if len(user.Emails) > 0 {
				cm.Email = utils.StringValue(user.Emails[0].EmailAddress)
			}
		}
	}
	return nil
}

func (s service) fillProjectInfo(claManagers []*models.CompanyClaManager) error {
	var wg sync.WaitGroup
	wg.Add(len(claManagers))
	for i, _ := range claManagers {
		go func(claManager *models.CompanyClaManager) {
			defer wg.Done()
			project, err := s.projectRepo.GetProjectByID(claManager.ProjectID)
			if err != nil {
				logging.Warnf("Unable to fetch project details for project %s. error = %s", claManager.ProjectID, err)
				return
			}
			claManager.ProjectName = project.ProjectName
			claManager.ProjectSfid = project.ProjectExternalID
		}(claManagers[i])
	}
	wg.Wait()
	return nil
}
