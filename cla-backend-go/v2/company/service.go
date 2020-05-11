package company

import (
	"sort"
	"strings"
	"sync"

	"github.com/communitybridge/easycla/cla-backend-go/logging"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/aws/aws-sdk-go/aws"
	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
	v1ProjectParams "github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/project"
	v1SignatureParams "github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/signatures"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/signatures"
	v2ProjectService "github.com/communitybridge/easycla/cla-backend-go/v2/project-service"
	v2ProjectServiceModels "github.com/communitybridge/easycla/cla-backend-go/v2/project-service/models"
	v2UserService "github.com/communitybridge/easycla/cla-backend-go/v2/user-service"
	v2UserServiceModels "github.com/communitybridge/easycla/cla-backend-go/v2/user-service/models"
)

// Service functions for company
type Service interface {
	GetCompanyCLAManagers(companyID string) (*models.CompanyClaManagers, error)
	GetCompanyActiveCLAs(companyID string) (*models.ActiveClaList, error)
}

// ProjectRepo contains project repo methods
type ProjectRepo interface {
	GetProjectByID(projectID string) (*v1Models.Project, error)
	GetProjectsByExternalID(params *v1ProjectParams.GetProjectsByExternalIDParams) (*v1Models.Projects, error)
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

func signedCLAFilename(projectID string, claType string, identifier string, signatureID string) string {
	return strings.Join([]string{"contract-group", projectID, claType, identifier, signatureID}, "/") + ".pdf"
}

func (s *service) getAllCCLASignatures(companyID string) ([]*v1Models.Signature, error) {
	var sigs []*v1Models.Signature
	var lastScannedKey *string
	for {
		signatures, err := s.signatureRepo.GetCompanySignatures(v1SignatureParams.GetCompanySignaturesParams{
			CompanyID:     companyID,
			SignatureType: aws.String("ccla"),
			NextKey:       lastScannedKey,
		}, 1000, signatures.DontLoadACLDetails)
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
				ApprovedOn: sig.SignatureCreated,
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
				log.Warnf("Unable to fetch project details for project %s. error = %s", projectID, err)
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

func (s *service) GetCompanyActiveCLAs(companyID string) (*models.ActiveClaList, error) {
	var out models.ActiveClaList
	sigs, err := s.getAllCCLASignatures(companyID)
	if err != nil {
		return nil, err
	}
	out.List = make([]*models.ActiveCla, 0, len(sigs))
	if len(sigs) == 0 {
		return &out, nil
	}
	var wg sync.WaitGroup
	wg.Add(len(sigs))
	for _, sig := range sigs {
		activeCla := &models.ActiveCla{}
		out.List = append(out.List, activeCla)
		go func(swg *sync.WaitGroup, signature *v1Models.Signature, acla *models.ActiveCla) {
			s.fillActiveCLA(swg, signature, acla)
		}(&wg, sig, activeCla)
	}
	wg.Wait()
	return &out, nil
}

func (s *service) fillActiveCLA(wg *sync.WaitGroup, sig *v1Models.Signature, activeCla *models.ActiveCla) {
	defer wg.Done()
	p, err := s.projectRepo.GetProjectByID(sig.ProjectID)
	if err != nil {
		log.Error("fillActiveCLA : unable to get project", err)
		return
	}
	psc := v2ProjectService.GetClient()
	projectDetails, err := psc.GetProject(p.ProjectExternalID)
	if err != nil {
		log.Error("fillActiveCLA : unable to get project details", err)
		return
	}

	activeCla.ProjectID = sig.ProjectID
	activeCla.ProjectName = projectDetails.Name
	activeCla.ProjectSfid = p.ProjectExternalID
	activeCla.ProjectType = projectDetails.ProjectType
	activeCla.ProjectLogo = projectDetails.ProjectLogo
	activeCla.SignedOn = sig.SignatureCreated
	activeCla.ClaGroupName = p.ProjectName
	activeCla.SubProjects = make([]*models.SubProject, 0, len(projectDetails.Projects))

	var subProjects []*v2ProjectServiceModels.ProjectOutput
	var signatoryName string
	var cwg sync.WaitGroup
	cwg.Add(3)

	var cclaURL string
	go func() {
		defer cwg.Done()
		cclaURL, err = utils.GetDownloadLink(signedCLAFilename(sig.ProjectID, sig.SignatureType, sig.SignatureReferenceID, sig.SignatureID))
		if err != nil {
			log.Error("fillActiveCLA : unable to get ccla s3 link", err)
			return
		}
	}()

	go func() {
		defer cwg.Done()
		usc := v2UserService.GetClient()
		if len(sig.SignatureACL) == 0 {
			log.Warnf("signature : %s have empty signature_acl", sig.SignatureID)
			return
		}
		lfUsername := sig.SignatureACL[0].LfUsername
		user, err := usc.GetUserByUsername(lfUsername)
		if err != nil {
			log.Warnf("unable to get user with lf username : %s", lfUsername)
			return
		}
		signatoryName = user.Name
	}()

	go func() {
		defer cwg.Done()
		subProjects = s.filterClaProjects(projectDetails.Projects)
	}()
	cwg.Wait()

	activeCla.SignatoryName = signatoryName
	activeCla.CclaURL = cclaURL
	for _, subProject := range subProjects {
		sp := &models.SubProject{
			ProjectName: subProject.Name,
			ProjectSfid: subProject.ID,
			ProjectLogo: subProject.ProjectLogo,
			ProjectType: subProject.ProjectType,
		}
		activeCla.SubProjects = append(activeCla.SubProjects, sp)
	}
}

// return projects output for which cla_group is present in cla
func (s *service) filterClaProjects(projects []*v2ProjectServiceModels.ProjectOutput) []*v2ProjectServiceModels.ProjectOutput {
	results := make([]*v2ProjectServiceModels.ProjectOutput, 0)
	prChan := make(chan *v2ProjectServiceModels.ProjectOutput)
	for _, v := range projects {
		go func(projectOutput *v2ProjectServiceModels.ProjectOutput) {
			project, err := s.projectRepo.GetProjectsByExternalID(&v1ProjectParams.GetProjectsByExternalIDParams{
				ExternalID: projectOutput.ID,
				PageSize:   aws.Int64(1),
			})
			if err != nil {
				log.Warnf("Unable to fetch project details for project with external id %s. error = %s", projectOutput.ID, err)
				prChan <- nil
				return
			}
			if project.ResultCount == 0 {
				prChan <- nil
				return
			}
			prChan <- projectOutput
		}(v)
	}
	for range projects {
		project := <-prChan
		if project != nil {
			results = append(results, project)
		}
	}
	return results
}
