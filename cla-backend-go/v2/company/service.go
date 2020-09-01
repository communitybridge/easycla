// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package company

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/sirupsen/logrus"

	"github.com/communitybridge/easycla/cla-backend-go/events"
	"github.com/communitybridge/easycla/cla-backend-go/projects_cla_groups"

	"github.com/jinzhu/copier"

	"github.com/LF-Engineering/lfx-kit/auth"
	"github.com/communitybridge/easycla/cla-backend-go/company"

	"github.com/aws/aws-sdk-go/aws"
	v1Company "github.com/communitybridge/easycla/cla-backend-go/company"
	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
	v1ProjectParams "github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/project"
	v1SignatureParams "github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/signatures"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	v2Models "github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/logging"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/signatures"
	"github.com/communitybridge/easycla/cla-backend-go/users"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
	acs_service "github.com/communitybridge/easycla/cla-backend-go/v2/acs-service"

	orgService "github.com/communitybridge/easycla/cla-backend-go/v2/organization-service"
	"github.com/communitybridge/easycla/cla-backend-go/v2/organization-service/client/organizations"
	v2ProjectService "github.com/communitybridge/easycla/cla-backend-go/v2/project-service"
	v2ProjectServiceModels "github.com/communitybridge/easycla/cla-backend-go/v2/project-service/models"
	v2UserService "github.com/communitybridge/easycla/cla-backend-go/v2/user-service"
	v2UserServiceModels "github.com/communitybridge/easycla/cla-backend-go/v2/user-service/models"
)

// errors
var (
	ErrProjectNotFound = errors.New("project not found")
	ErrCLAUserNotFound = errors.New("claUser not found")
	ErrNoLfUsername    = errors.New("user has no LF username")
	// ErrNoValidEmail    = errors.New("user with no valid email")

	//ErrProjectSigned returns error if project already signed
	ErrProjectSigned = errors.New("project already signed")
	//ErrLFXUserNotFound when user-service fails to find user
	ErrLFXUserNotFound = errors.New("lfx user not found")
	//ErrContributorConflict when user is already assigned contributor role
	ErrContributorConflict = errors.New("user already assigned contributor")
	//ErrRoleScopeConflict thrown if user already has role scope
	ErrRoleScopeConflict = errors.New("user is already contributor")
)

// constants
const (
	// used when we want to query all data from dependent service.
	HugePageSize = int64(10000)
	// LoadRepoDetails     = true
	DontLoadRepoDetails = false
	// FoundationType the SF foundation type string - previously was "Foundation", now "Project Group"
	FoundationType = "Project Group"
	// Lead representing type of user
	Lead = "lead"
	//NoAccount
	NoAccount = "Individual - No Account"
)

// Service functions for company
type Service interface {
	GetCompanyProjectCLAManagers(companyID string, projectSFID string) (*models.CompanyClaManagers, error)
	GetCompanyProjectActiveCLAs(companyID string, projectSFID string) (*models.ActiveClaList, error)
	GetCompanyProjectContributors(projectSFID string, companySFID string, searchTerm string) (*models.CorporateContributorList, error)
	GetCompanyProjectCLA(authUser *auth.User, companySFID, projectSFID string) (*models.CompanyProjectClaList, error)
	CreateCompany(companyName string, companyWebsite string, userEmail string, userID string, LFXPortalURL string) (*models.CompanyOutput, error)
	GetCompanyByName(companyName string) (*models.Company, error)
	GetCompanyByID(companyID string) (*models.Company, error)
	GetCompanyBySFID(companySFID string) (*models.Company, error)
	DeleteCompanyByID(companyID string) error
	DeleteCompanyBySFID(companySFID string) error
	GetCompanyCLAGroupManagers(companyID, claGroupID string) (*models.CompanyClaManagers, error)
	AssociateContributor(companySFID, userEmail string) (*models.Contributor, error)
	AssociateContributorByGroup(companySFID, userEmail string, projectCLAGroups []*projects_cla_groups.ProjectClaGroup, f logrus.Fields, ClaGroupID string) ([]*models.Contributor, string, error)
	GetCompanyAdmins(companyID string) (*models.CompanyAdminList, error)
}

// ProjectRepo contains project repo methods
type ProjectRepo interface {
	GetCLAGroupByID(projectID string, loadRepoDetails bool) (*v1Models.Project, error)
	GetCLAGroupsByExternalID(params *v1ProjectParams.GetProjectsByExternalIDParams, loadRepoDetails bool) (*v1Models.Projects, error)
}

// NewService returns instance of company service
func NewService(v1CompanyService v1Company.IService, sigRepo signatures.SignatureRepository, projectRepo ProjectRepo, usersRepo users.UserRepository, companyRepo company.IRepository, pcgRepo projects_cla_groups.Repository, evService events.Service) Service {
	return &service{
		v1CompanyService:     v1CompanyService,
		signatureRepo:        sigRepo,
		projectRepo:          projectRepo,
		userRepo:             usersRepo,
		companyRepo:          companyRepo,
		projectClaGroupsRepo: pcgRepo,
		eventService:         evService,
	}
}

func (s *service) GetCompanyProjectCLAManagers(companyID string, projectSFID string) (*models.CompanyClaManagers, error) {
	f := logrus.Fields{"projectSFID": projectSFID, "companyID": companyID}
	var err error
	claGroups, err := s.getCLAGroupsUnderProjectOrFoundation(projectSFID)
	if err != nil {
		log.WithFields(f).Warnf("problem fetching CLA Groups under project or foundation, error: %+v", err)
		return nil, err
	}
	sigs, err := s.getAllCCLASignatures(companyID)
	if err != nil {
		log.WithFields(f).Warnf("problem fetching CLA signatures, error: %+v", err)
		return nil, err
	}
	claManagers := make([]*models.CompanyClaManager, 0)
	lfUsernames := utils.NewStringSet()
	// Get CLA managers
	for _, sig := range sigs {
		if _, ok := claGroups[sig.ProjectID]; !ok {
			continue
		}
		for _, user := range sig.SignatureACL {
			claManagers = append(claManagers, &models.CompanyClaManager{
				// DB doesn't have approved_on value
				ApprovedOn: sig.SignatureCreated,
				LfUsername: user.LfUsername,
				ProjectID:  sig.ProjectID,
			})
			lfUsernames.Add(user.LfUsername)
		}
	}
	// get userinfo and project info
	var usermap map[string]*v2UserServiceModels.User
	usermap, err = getUsersInfo(lfUsernames.List())
	if err != nil {
		log.WithFields(f).Warnf("problem fetching users information, error: %+v", err)
		return nil, err
	}
	// fill user info
	fillUsersInfo(claManagers, usermap)
	// fill project info
	fillProjectInfo(claManagers, claGroups)
	// sort result by cla manager name
	sort.Slice(claManagers, func(i, j int) bool {
		return claManagers[i].Name < claManagers[j].Name
	})
	return &models.CompanyClaManagers{List: claManagers}, nil
}

func (s *service) GetCompanyAdmins(companySFID string) (*models.CompanyAdminList, error) {
	f := logrus.Fields{"companySFID": companySFID}
	orgClient := orgService.GetClient()

	log.WithFields(f).Info("Getting user admins for company")
	admins, adminErr := orgClient.ListOrgUserAdminScopes(companySFID)
	adminList := make([]*models.AdminSf, 0)
	if adminErr != nil {
		if _, ok := adminErr.(*organizations.ListOrgUsrAdminScopesNotFound); ok {
			log.WithFields(f).Info(" No admins found ")
			return &models.CompanyAdminList{
				List: adminList,
			}, nil
		}
		return nil, adminErr
	}

	// if 404 and no error parse the userroles list
	for _, userRole := range admins.Userroles {
		adminList = append(adminList, &models.AdminSf{
			Email:    userRole.Contact.EmailAddress,
			Username: userRole.Contact.Username,
			ID:       userRole.Contact.ID,
		})
	}

	return &models.CompanyAdminList{
		List: adminList,
	}, nil
}

func (s *service) GetCompanyProjectActiveCLAs(companyID string, projectSFID string) (*models.ActiveClaList, error) {
	f := logrus.Fields{"projectSFID": projectSFID, "companyID": companyID}
	var err error
	claGroups, err := s.getCLAGroupsUnderProjectOrFoundation(projectSFID)
	if err != nil {
		log.WithFields(f).Warnf("problem fetching CLA Groups under project or foundation, error: %+v", err)
		return nil, err
	}
	var out models.ActiveClaList
	sigs, err := s.getAllCCLASignatures(companyID)
	if err != nil {
		log.WithFields(f).Warnf("problem fetching CCLA signatures, error: %+v", err)
		return nil, err
	}
	out.List = make([]*models.ActiveCla, 0, len(sigs))
	if len(sigs) == 0 {
		return &out, nil
	}
	var wg sync.WaitGroup
	wg.Add(len(sigs))
	for _, sig := range sigs {
		if _, ok := claGroups[sig.ProjectID]; !ok {
			// skip the cla_group which are not under current foundation/project
			wg.Done()
			continue
		}
		activeCla := &models.ActiveCla{}
		out.List = append(out.List, activeCla)
		go func(swg *sync.WaitGroup, signature *v1Models.Signature, acla *models.ActiveCla) {
			s.fillActiveCLA(swg, signature, acla, claGroups)
		}(&wg, sig, activeCla)
	}
	wg.Wait()
	return &out, nil
}

func (s *service) GetCompanyProjectContributors(projectSFID string, companySFID string, searchTerm string) (*models.CorporateContributorList, error) {
	f := logrus.Fields{"projectSFID": projectSFID, "companySFID": companySFID, "searchTerm": searchTerm}
	list := make([]*models.CorporateContributor, 0)
	sigs, err := s.getAllCompanyProjectEmployeeSignatures(companySFID, projectSFID)
	if err != nil {
		log.WithFields(f).Warnf("problem fetching all company project employee signatures, error: %+v", err)
		return nil, err
	}
	if len(sigs) == 0 {
		return &models.CorporateContributorList{
			List: list,
		}, nil
	}
	var wg sync.WaitGroup
	result := make(chan *models.CorporateContributor)
	wg.Add(len(sigs))
	go func() {
		wg.Wait()
		close(result)
	}()

	for _, sig := range sigs {
		go fillCorporateContributorModel(&wg, s.userRepo, sig, result, searchTerm)
	}

	for corpContributor := range result {
		list = append(list, corpContributor)
	}

	return &models.CorporateContributorList{
		List: list,
	}, nil
}

func (s *service) CreateCompany(companyName string, companyWebsite string, userEmail string, userID string, LFXPortalURL string) (*models.CompanyOutput, error) {
	f := logrus.Fields{"companyName": companyName, "companyWebsite": companyWebsite, "userID": userID, "LFXPortalURL": LFXPortalURL}
	var lfUser *v2UserServiceModels.User

	// Create Sales Force company
	orgClient := orgService.GetClient()
	log.Debugf("Creating Organization : %s Website: %s", companyName, companyWebsite)
	org, err := orgClient.CreateOrg(companyName, companyWebsite)
	if err != nil {
		log.WithFields(f).Warnf("unable to create platform organization service, error: %+v", err)
		return nil, err
	}

	// acsClient := acs_service.GetClient()
	userClient := v2UserService.GetClient()

	lfUser, lfErr := userClient.SearchUserByEmail(userEmail)
	if lfErr != nil {
		msg := fmt.Sprintf("User : %s has no LFID", userEmail)
		log.Warn(msg)
	}
	if lfUser != nil {

		log.WithFields(f).Debugf("User :%s has been assigned the company-owner role to organization: %s ", userEmail, org.Name)
		// Associate User (Not associated) with Newly created org
		if lfUser.Account.ID == NoAccount {
			lfUser.Account.ID = org.ID
		}

		//Send Email to User with instructions to complete Company profile
		log.WithFields(f).Debugf("Sending Email to user :%s to complete setup for newly created Org: %s ", userEmail, org.Name)
		sendEmailToUserCompanyProfile(org.Name, userEmail, lfUser.Username, LFXPortalURL)
	}

	// Create Easy CLA Company
	log.WithFields(f).Debugf("Creating EasyCLA company: %s ", companyName)
	// OrgID used as externalID for the easyCLA Company
	// Create a new company model for the create function
	createCompanyModel := &v1Models.Company{
		CompanyACL:        nil,
		CompanyExternalID: org.ID,
		CompanyManagerID:  userID,
		CompanyName:       companyName,
	}

	_, createErr := s.companyRepo.CreateCompany(createCompanyModel)
	//easyCLAErr := s.repo.CreateCompany(companyName, org.ID, userID)
	if createErr != nil {
		log.WithFields(f).Warnf("Failed to create EasyCLA company for company: %s, error: %+v",
			companyName, createErr)
		return nil, createErr
	}

	return &models.CompanyOutput{
		CompanyName:    org.Name,
		CompanyWebsite: companyWebsite,
		LogoURL:        org.LogoURL,
		CompanyID:      org.ID,
	}, nil
}

// GetCompanyByName deletes the company by name
func (s *service) GetCompanyByName(companyName string) (*models.Company, error) {
	companyModel, err := s.companyRepo.GetCompanyByName(companyName)
	if err != nil {
		return nil, err
	}

	if companyModel == nil {
		log.Debugf("search by company name: %s didn't locate the record", companyName)
		return nil, nil
	}

	// Convert from v1 to v2 model - use helper: Copy(toValue interface{}, fromValue interface{})
	var v2CompanyModel v2Models.Company
	copyErr := copier.Copy(&v2CompanyModel, &companyModel)
	if copyErr != nil {
		log.Warnf("problem converting v1 company model to a v2 company model, error: %+v", copyErr)
		return nil, copyErr
	}

	return &v2CompanyModel, nil
}

// GetCompanyByID retrieves the company by internal ID
func (s *service) GetCompanyByID(companyID string) (*models.Company, error) {
	companyModel, err := s.companyRepo.GetCompany(companyID)
	if err != nil {
		return nil, err
	}

	if companyModel == nil {
		log.Debugf("search by company ID: %s didn't locate the record", companyID)
		return nil, nil
	}

	// Convert from v1 to v2 model - use helper: Copy(toValue interface{}, fromValue interface{})
	var v2CompanyModel v2Models.Company
	copyErr := copier.Copy(&v2CompanyModel, &companyModel)
	if copyErr != nil {
		log.Warnf("problem converting v1 company model to a v2 company model, error: %+v", copyErr)
		return nil, copyErr
	}

	return &v2CompanyModel, nil
}

func (s *service) AssociateContributor(companySFID string, userEmail string) (*models.Contributor, error) {
	f := logrus.Fields{"companySFID": companySFID, "userEmail": userEmail}

	orgClient := orgService.GetClient()

	userService := v2UserService.GetClient()
	log.WithFields(f).Info("searching for LFX User")
	lfxUser, userErr := userService.SearchUserByEmail(userEmail)
	if userErr != nil {
		log.WithFields(f).Warnf("unable to get user")
		return nil, userErr
	}

	acsServiceClient := acs_service.GetClient()

	log.WithFields(f).Info("Getting roleID for the contributor role")
	roleID, roleErr := acsServiceClient.GetRoleID("contributor")
	if roleErr != nil {
		log.WithFields(f).Warn("Problem getting roleID for contributor role ")
		return nil, roleErr
	}

	log.WithFields(f).Info("creating contributor role scope")
	scopeErr := orgClient.CreateOrgUserRoleOrgScope(userEmail, companySFID, roleID)
	if scopeErr != nil {
		log.WithFields(f).Warnf("Problem creating role scope")
		return nil, scopeErr
	}

	contributor := &models.Contributor{
		LfUsername:  lfxUser.Username,
		UserSfid:    lfxUser.ID,
		Email:       strfmt.Email(userEmail),
		AssignedOn:  time.Now().String(),
		CompanySfid: companySFID,
		Role:        *aws.String("contributor"),
	}

	return contributor, nil
}

//CreateContibutor creates contributor for contributor prospect
func (s *service) CreateContibutor(companyID string, projectID string, userEmail string, ClaGroupID string) (*models.Contributor, error) {
	// integrate user,acs,org and project services
	userClient := v2UserService.GetClient()
	acServiceClient := acs_service.GetClient()
	orgClient := orgService.GetClient()

	isSigned, signedErr := s.isSigned(companyID, projectID)
	if signedErr != nil {
		msg := fmt.Sprintf("EasyCLA - 400 Bad Request- %s", signedErr)
		log.Warn(msg)
		return nil, signedErr
	}

	if isSigned {
		msg := fmt.Sprintf("EasyCLA - 400 Bad Request - Project :%s is already signed ", projectID)
		log.Warn(msg)
		return nil, ErrProjectSigned
	}

	user, userErr := userClient.SearchUserByEmail(userEmail)
	if userErr != nil {
		log.Debugf("Failed to get user by email: %s , error: %+v", userEmail, userErr)
		return nil, ErrLFXUserNotFound
	}

	// Check if user is already contributor of project|organization scope
	hasRoleScope, hasRoleScopeErr := orgClient.IsUserHaveRoleScope("contributor", user.ID, companyID, projectID)
	if hasRoleScopeErr != nil {
		// Skip 404 for ListOrgUsrServiceScopes endpoint
		if _, ok := hasRoleScopeErr.(*organizations.ListOrgUsrServiceScopesNotFound); !ok {
			log.Debugf("Failed to check roleScope: contributor  for user: %s", user.Username)
			return nil, hasRoleScopeErr
		}
	}
	if hasRoleScope {
		log.Debugf("Conflict ")
		return nil, ErrContributorConflict
	}

	roleID, designeeErr := acServiceClient.GetRoleID("contributor")
	if designeeErr != nil {
		msg := "Problem getting role ID for contributor"
		log.Warn(msg)
		return nil, designeeErr
	}

	scopeErr := orgClient.CreateOrgUserRoleOrgScopeProjectOrg(userEmail, projectID, companyID, roleID)
	if scopeErr != nil {
		msg := fmt.Sprintf("Problem creating projectOrg scope for email: %s , projectID: %s, companyID: %s", userEmail, projectID, companyID)
		log.Warn(msg)
		if _, ok := scopeErr.(*organizations.CreateOrgUsrRoleScopesConflict); ok {
			return nil, ErrRoleScopeConflict
		}
		return nil, scopeErr
	}

	v1Company, companyErr := s.v1CompanyService.GetCompanyByExternalID(companyID)
	if companyErr != nil {
		log.Error("company not found", companyErr)
	}

	projectModel, projErr := s.projectRepo.GetCLAGroupByID(ClaGroupID, DontLoadRepoDetails)
	if projErr != nil {
		msg := fmt.Sprintf("unable to query CLA Group ID: %s, error: %+v", ClaGroupID, projErr)
		log.Warnf(msg)
	}

	// Log Event
	s.eventService.LogEvent(
		&events.LogEventArgs{
			EventType:         events.AssignUserRoleScopeType,
			LfUsername:        user.Username,
			UserID:            user.ID,
			ExternalProjectID: projectID,
			CompanyModel:      v1Company,
			ProjectModel:      projectModel,
			UserModel:         &v1Models.User{LfUsername: user.Username, UserID: user.ID},
			EventData: &events.AssignRoleScopeData{
				Role:  "contributor",
				Scope: fmt.Sprintf("%s|%s", projectID, companyID),
			},
		})

	if user.Type == Lead {
		log.Debugf("Converting user: %s from lead to contact ", userEmail)
		contactErr := userClient.ConvertToContact(user.ID)
		if contactErr != nil {
			log.Debugf("failed to convert user: %s to contact ", userEmail)
			return nil, contactErr
		}
		// Log user conversion event
		s.eventService.LogEvent(&events.LogEventArgs{
			EventType:         events.ConvertUserToContactType,
			LfUsername:        user.Username,
			ExternalProjectID: projectID,
			EventData:         &events.UserConvertToContactData{},
		})
	}

	contributor := &models.Contributor{
		LfUsername:  user.Username,
		UserSfid:    user.ID,
		AssignedOn:  time.Now().String(),
		Email:       strfmt.Email(userEmail),
		CompanySfid: companyID,
		Role:        *aws.String("contributor"),
	}
	return contributor, nil
}

// Helper function to check if project/claGroup is signed
func (s *service) isSigned(companyID string, projectID string) (bool, error) {
	isSigned := false
	// Check for company contributor

	v1Company, companyErr := s.v1CompanyService.GetCompanyByExternalID(companyID)
	if companyErr != nil {
		msg := fmt.Sprintf("EasyCLA - 400 Bad Request - %s", companyErr)
		log.Warn(msg)
		return isSigned, companyErr
	}

	claManagers, err := s.GetCompanyProjectCLAManagers(v1Company.CompanyID, projectID)
	if err != nil {
		msg := fmt.Sprintf("EasyCLA - 400 Bad Request : %v", err)
		log.Warn(msg)
		return isSigned, err
	}

	if len(claManagers.List) > 0 {
		msg := fmt.Sprintf("EasyCLA - 400 Bad Request - CLA Group signed for company: %s and project : %s", companyID, projectID)
		log.Warn(msg)
		isSigned = true
	}

	return isSigned, nil
}

//AssociateContributorByGroup creates contributor by group for contributor prospect
func (s *service) AssociateContributorByGroup(companySFID, userEmail string, projectCLAGroups []*projects_cla_groups.ProjectClaGroup, f logrus.Fields, ClaGroupID string) ([]*models.Contributor, string, error) {
	var contributors []*models.Contributor
	foundationSFID := projectCLAGroups[0].FoundationSFID
	if foundationSFID != "" {
		contributor, err := s.CreateContibutor(companySFID, foundationSFID, userEmail, ClaGroupID)
		if err != nil {
			if err == ErrContributorConflict {
				msg := fmt.Sprintf("Conflict assigning contributor role for Foundation SFID: %s ", foundationSFID)
				return nil, msg, err
			}
			msg := fmt.Sprintf("Creating contributor failed for Foundation SFID: %s ", foundationSFID)
			return nil, msg, err
		}
		contributors = append(contributors, contributor)
	}

	for _, pcg := range projectCLAGroups {
		log.WithFields(f).Debugf("creating contributor for Project SFID: %s", pcg.ProjectSFID)
		if foundationSFID != pcg.ProjectSFID {
			contributor, err := s.CreateContibutor(companySFID, pcg.ProjectSFID, userEmail, ClaGroupID)
			if err != nil {
				if err == ErrContributorConflict {
					msg := fmt.Sprintf("Conflict assigning contributor role for Project SFID: %s ", pcg.ProjectSFID)
					return nil, msg, err
				}
				msg := fmt.Sprintf("Creating contributor failed for Project SFID: %s ", pcg.ProjectSFID)
				return nil, msg, err
			}
			contributors = append(contributors, contributor)
		}

	}
	return contributors, "", nil
}

// GetCompanyBySFID retrieves the company by external SFID
func (s *service) GetCompanyBySFID(companySFID string) (*models.Company, error) {
	f := logrus.Fields{"companySFID": companySFID}
	companyModel, err := s.companyRepo.GetCompanyByExternalID(companySFID)
	if err != nil {
		// If we were unable to find the company/org in our local database, try to auto-create based
		// on the existing SF record
		if err == company.ErrCompanyDoesNotExist {
			log.WithFields(f).Debug("company not found in EasyCLA database - attempting to auto-create from platform organization service record")
			newCompanyModel, createCompanyErr := s.autoCreateCompany(companySFID)
			if createCompanyErr != nil {
				log.WithFields(f).Warnf("problem creating company from platform organization SF record, error: %+v",
					createCompanyErr)
				return nil, createCompanyErr
			}
			if newCompanyModel == nil {
				log.WithFields(f).Warnf("problem creating company from SF records - created model is nil")
				return nil, company.ErrCompanyDoesNotExist
			}
			// Success, fall through and continue processing
			companyModel = newCompanyModel
		}
		return nil, err
	}

	if companyModel == nil {
		log.WithFields(f).Debugf("search by company SFID didn't locate the record")
		return nil, nil
	}

	// Convert from v1 to v2 model - use helper: Copy(toValue interface{}, fromValue interface{})
	var v2CompanyModel v2Models.Company
	copyErr := copier.Copy(&v2CompanyModel, &companyModel)
	if copyErr != nil {
		log.WithFields(f).Warnf("problem converting v1 company model to a v2 company model, error: %+v", copyErr)
		return nil, copyErr
	}

	return &v2CompanyModel, nil
}

// DeleteCompanyByID deletes the company by ID
func (s *service) DeleteCompanyByID(companyID string) error {
	return s.companyRepo.DeleteCompanyByID(companyID)
}

// DeleteCompanyBySFID deletes the company by SFID
func (s *service) DeleteCompanyBySFID(companyID string) error {
	return s.companyRepo.DeleteCompanyBySFID(companyID)
}

func (s *service) GetCompanyProjectCLA(authUser *auth.User, companySFID, projectSFID string) (*models.CompanyProjectClaList, error) {
	f := logrus.Fields{"authUserName": authUser.UserName, "authUserEmail": authUser.Email, "companySFID": companySFID, "projectSFID": projectSFID}
	var canSign bool
	resources := authUser.ResourceIDsByTypeAndRole(auth.ProjectOrganization, "cla-manager-designee")
	projectOrg := fmt.Sprintf("%s|%s", projectSFID, companySFID)
	for _, r := range resources {
		if r == projectOrg {
			canSign = true
			break
		}
	}

	// Attempt to locate the company model in our database
	log.WithFields(f).Debug("locating company by SF ID")
	var companyModel *v1Models.Company
	companyModel, companyErr := s.companyRepo.GetCompanyByExternalID(companySFID)
	if companyErr != nil {
		// If we were unable to find the company/org in our local database, try to auto-create based
		// on the existing SF record
		if companyErr == company.ErrCompanyDoesNotExist {

			log.WithFields(f).Debug("company not found in EasyCLA database - attempting to auto-create from platform organization service record")
			var createCompanyErr error
			companyModel, createCompanyErr = s.autoCreateCompany(companySFID)
			if createCompanyErr != nil {
				log.WithFields(f).Warnf("problem creating company from platform organization SF record, error: %+v",
					createCompanyErr)
				return nil, createCompanyErr
			}
			if companyModel == nil {
				log.WithFields(f).Warnf("problem creating company from SF records - created model is nil")
				return nil, company.ErrCompanyDoesNotExist
			}
			// Success, fall through and continue processing
		} else {
			return nil, companyErr
		}
	}

	claGroups, err := s.getCLAGroupsUnderProjectOrFoundation(projectSFID)
	if err != nil {
		log.WithFields(f).Warnf("problem fetching CLA Groups under project or foundation, error: %+v", err)
		return nil, err
	}

	activeCLAList, err := s.GetCompanyProjectActiveCLAs(companyModel.CompanyID, projectSFID)
	if err != nil {
		log.WithFields(f).Warnf("problem fetching company project active CLAs, error: %+v", err)
		return nil, err
	}

	resp := &models.CompanyProjectClaList{
		SignedClaList:       activeCLAList.List,
		UnsignedProjectList: make([]*models.UnsignedProject, 0),
	}

	for _, activeCLA := range activeCLAList.List {
		// remove cla groups for which we have signed cla
		log.WithFields(f).Debugf("removing CLA Groups with active CLA, CLA Group: %+v, error: %+v", activeCLA, err)
		delete(claGroups, activeCLA.ProjectID)
	}

	// fill details for not signed cla
	for claGroupID, claGroup := range claGroups {
		unsignedProject := &models.UnsignedProject{
			CanSign:      canSign,
			ClaGroupID:   claGroupID,
			ClaGroupName: claGroup.ClaGroupName,
			ProjectName:  claGroup.ProjectName,
			ProjectSfid:  claGroup.ProjectSFID,
			SubProjects:  claGroup.SubProjects,
			IclaEnabled:  claGroup.IclaEnabled,
			CclaEnabled:  claGroup.CclaEnabled,
		}
		log.WithFields(f).Debugf("adding unsigned CLA Group: %+v, error: %+v", unsignedProject, err)
		resp.UnsignedProjectList = append(resp.UnsignedProjectList, unsignedProject)
	}

	return resp, nil
}

// GetCompanyCLAGroupManagers when provided the internal company ID and CLA Groups ID, this routine returns the list of
// corresponding CLA managers
func (s *service) GetCompanyCLAGroupManagers(companyID, claGroupID string) (*models.CompanyClaManagers, error) {
	f := logrus.Fields{"companyID": companyID, "claGroupID": claGroupID}
	signed, approved := true, true
	pageSize := int64(10)
	sigModel, err := s.signatureRepo.GetProjectCompanySignature(companyID, claGroupID, &signed, &approved, nil, &pageSize)
	if err != nil {
		log.WithFields(f).Warnf("unable to query CCLA signature using Company ID: %s and CLA Group ID: %s, signed: true, approved: true, error: %+v",
			companyID, claGroupID, err)
		return nil, err
	}

	if sigModel == nil {
		log.WithFields(f).Warnf("unable to query CCLA signature using Company ID: %s and CLA Group ID: %s, signed: true, approved: true - no signature found",
			companyID, claGroupID)
		return nil, nil
	}

	projectModel, projErr := s.projectRepo.GetCLAGroupByID(claGroupID, DontLoadRepoDetails)
	if projErr != nil {
		log.WithFields(f).Warnf("unable to query CLA Group ID: %s, error: %+v", claGroupID, err)
		return nil, err
	}

	if projectModel == nil {
		log.WithFields(f).Warnf("unable to query CLA Group ID: %s - no CLA Group found", claGroupID)
		return nil, nil
	}

	companyModel, companyErr := s.companyRepo.GetCompany(companyID)
	if companyErr != nil {
		log.WithFields(f).Warnf("unable to query Company ID: %s, error: %+v", companyID, companyErr)
		return nil, err
	}

	if companyModel == nil {
		log.WithFields(f).Warnf("unable to query Company ID: %s - no company by ID found", companyID)
		return nil, nil
	}

	claManagers := make([]*models.CompanyClaManager, 0)
	for _, user := range sigModel.SignatureACL {
		claManagers = append(claManagers, &models.CompanyClaManager{
			// DB doesn't have approved_on value - just use sig created date/time
			ApprovedOn:       sigModel.SignatureCreated,
			LfUsername:       user.LfUsername,
			Email:            strfmt.Email(user.LfEmail),
			Name:             user.Username,
			UserSfid:         user.UserExternalID,
			ProjectID:        sigModel.ProjectID,
			ProjectName:      projectModel.ProjectName,
			ClaGroupName:     projectModel.ProjectName,
			OrganizationName: companyModel.CompanyName,
			OrganizationSfid: companyModel.CompanyExternalID,
		})
	}

	return &models.CompanyClaManagers{List: claManagers}, nil
}

func v2ProjectToMap(projectDetails *v2ProjectServiceModels.ProjectOutputDetailed) (map[string]*v2ProjectServiceModels.ProjectOutput, error) {
	epmap := make(map[string]*v2ProjectServiceModels.ProjectOutput) // key project_sfid
	var pr v2ProjectServiceModels.ProjectOutput
	err := copier.Copy(&pr, projectDetails)
	if err != nil {
		return nil, err
	}
	epmap[projectDetails.ID] = &pr
	for _, p := range projectDetails.Projects {
		epmap[p.ID] = p
	}
	return epmap, nil
}

func (s *service) getCLAGroupsUnderProjectOrFoundation(id string) (map[string]*claGroupModel, error) {
	result := make(map[string]*claGroupModel)
	psc := v2ProjectService.GetClient()
	projectDetails, err := psc.GetProject(id)
	if err != nil {
		return nil, err
	}
	var allProjectMapping []*projects_cla_groups.ProjectClaGroup
	if projectDetails.ProjectType == FoundationType {
		// get all projects for all cla group under foundation
		allProjectMapping, err = s.projectClaGroupsRepo.GetProjectsIdsForFoundation(id)
		if err != nil {
			return nil, err
		}
	} else {
		// get cla group id from project
		projectMapping, perr := s.projectClaGroupsRepo.GetClaGroupIDForProject(id)
		if perr != nil {
			return nil, err
		}
		// get all projects for that cla group
		allProjectMapping, err = s.projectClaGroupsRepo.GetProjectsIdsForClaGroup(projectMapping.ClaGroupID)
		if err != nil {
			return nil, err
		}
		if len(allProjectMapping) > 1 {
			// reload data in projectDetails for all projects of foundation
			projectDetails, err = psc.GetProject(projectDetails.Foundation.ID)
			if err != nil {
				return nil, err
			}
		}
	}
	// v2ProjectMap will contains projectSFID -> salesforce details of that project
	v2ProjectMap, err := v2ProjectToMap(projectDetails)
	if err != nil {
		return nil, err
	}
	// for all cla-groups create claGroupModel
	for _, pm := range allProjectMapping {
		cg, ok := result[pm.ClaGroupID]
		if !ok {
			cg = &claGroupModel{
				FoundationSFID: pm.FoundationSFID,
				SubProjects:    make([]string, 0),
			}
			result[pm.ClaGroupID] = cg
		}
		cg.SubProjectIDs = append(cg.SubProjectIDs, pm.ProjectSFID)
	}
	// if no cla-group found, return empty result
	if len(result) == 0 {
		return result, nil
	}
	var wg sync.WaitGroup
	wg.Add(len(result))
	for id, cg := range result {
		go func(claGroupID string, claGroup *claGroupModel) {
			defer wg.Done()
			// get cla-group info
			cginfo, err := s.projectRepo.GetCLAGroupByID(claGroupID, DontLoadRepoDetails)
			if err != nil || cginfo == nil {
				log.Warnf("Unable to get details of cla_group: %s", claGroupID)
				return
			}
			claGroup.ClaGroupName = cginfo.ProjectName
			claGroup.IclaEnabled = cginfo.ProjectICLAEnabled
			claGroup.CclaEnabled = cginfo.ProjectCCLAEnabled

			var pid string
			if len(claGroup.SubProjectIDs) == 1 {
				// use project info if cla-group have only one project
				pid = claGroup.SubProjectIDs[0]
			} else {
				// use foundation info if cla-group have multiple project
				pid = claGroup.FoundationSFID
				for _, spid := range claGroup.SubProjectIDs {
					subProject, ok := v2ProjectMap[spid]
					if !ok {
						log.Warnf("Unable to fill details for cla_group: %s with project details of %s", claGroupID, spid)
						return
					}
					claGroup.SubProjects = append(claGroup.SubProjects, subProject.Name)
				}
			}
			project, ok := v2ProjectMap[pid]
			if !ok {
				log.Warnf("Unable to fill details for cla_group: %s with project details of %s", claGroupID, claGroup.ProjectSFID)
				return
			}
			claGroup.ProjectLogo = project.ProjectLogo
			claGroup.ProjectName = project.Name
			claGroup.ProjectType = project.ProjectType
			claGroup.ProjectSFID = pid
		}(id, cg)
	}
	wg.Wait()
	return result, nil
}

func (s *service) getAllCCLASignatures(companyID string) ([]*v1Models.Signature, error) {
	var sigs []*v1Models.Signature
	var lastScannedKey *string
	for {
		sigModels, err := s.signatureRepo.GetCompanySignatures(v1SignatureParams.GetCompanySignaturesParams{
			CompanyID:     companyID,
			SignatureType: aws.String("ccla"),
			NextKey:       lastScannedKey,
		}, HugePageSize, signatures.DontLoadACLDetails)
		if err != nil {
			return nil, err
		}
		sigs = append(sigs, sigModels.Signatures...)
		if sigModels.LastKeyScanned == "" {
			break
		}
		lastScannedKey = aws.String(sigModels.LastKeyScanned)
	}
	return sigs, nil
}

func getUsersInfo(lfUsernames []string) (map[string]*v2UserServiceModels.User, error) {
	userMap := make(map[string]*v2UserServiceModels.User)
	if len(lfUsernames) == 0 {
		return userMap, nil
	}
	userServiceClient := v2UserService.GetClient()
	userModels, err := userServiceClient.GetUsersByUsernames(lfUsernames)
	if err != nil {
		return nil, err
	}
	for _, user := range userModels {
		userMap[user.Username] = user
	}
	return userMap, nil
}

func fillUsersInfo(claManagers []*models.CompanyClaManager, usermap map[string]*v2UserServiceModels.User) {
	for _, cm := range claManagers {
		user, ok := usermap[cm.LfUsername]
		if !ok {
			logging.Warnf("Unable to get user with username %s", cm.LfUsername)
			continue
		}
		cm.Name = user.Name
		// cm.LogoURL = user.LogoURL
		cm.UserSfid = user.ID
		for _, email := range user.Emails {
			if email != nil && email.IsPrimary != nil && *email.IsPrimary && email.EmailAddress != nil {
				cm.Email = strfmt.Email(*email.EmailAddress)
				break
			}
		}
	}
}

func fillProjectInfo(claManagers []*models.CompanyClaManager, claGroups map[string]*claGroupModel) {
	for _, claManager := range claManagers {
		cg, ok := claGroups[claManager.ProjectID]
		if !ok {
			continue
		}
		claManager.ClaGroupName = cg.ClaGroupName
		claManager.ProjectSfid = cg.ProjectSFID
		claManager.ProjectName = cg.ProjectName
	}
}

func (s *service) fillActiveCLA(wg *sync.WaitGroup, sig *v1Models.Signature, activeCla *models.ActiveCla, claGroups map[string]*claGroupModel) {
	defer wg.Done()
	cg, ok := claGroups[sig.ProjectID]
	if !ok {
		log.Warn("unable to get project details")
		return
	}

	// fill details from dynamodb
	activeCla.ProjectID = sig.ProjectID
	if sig.SignedOn == "" {
		activeCla.SignedOn = sig.SignatureCreated
	} else {
		activeCla.SignedOn = sig.SignedOn
	}
	activeCla.ClaGroupName = cg.ClaGroupName
	activeCla.SignatureID = sig.SignatureID.String()

	// fill details from project service
	activeCla.ProjectName = cg.ProjectName
	activeCla.ProjectSfid = cg.ProjectSFID
	activeCla.ProjectType = cg.ProjectType
	activeCla.ProjectLogo = cg.ProjectLogo
	activeCla.SubProjects = cg.SubProjects
	var signatoryName string
	var cwg sync.WaitGroup
	cwg.Add(2)

	var cclaURL string
	go func() {
		var err error
		defer cwg.Done()
		cclaURL, err = utils.GetDownloadLink(utils.SignedCLAFilename(sig.ProjectID, sig.SignatureType, sig.SignatureReferenceID.String(), sig.SignatureID.String()))
		if err != nil {
			log.Error("fillActiveCLA : unable to get ccla s3 link", err)
			return
		}
	}()

	go func() {
		defer cwg.Done()
		if sig.SignatoryName != "" {
			signatoryName = sig.SignatoryName
			return
		}
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

	cwg.Wait()

	activeCla.SignatoryName = signatoryName
	activeCla.CclaURL = cclaURL
}

// return projects output for which cla_group is present in cla
func (s *service) filterClaProjects(projects []*v2ProjectServiceModels.ProjectOutput) []*v2ProjectServiceModels.ProjectOutput { //nolint
	results := make([]*v2ProjectServiceModels.ProjectOutput, 0)
	prChan := make(chan *v2ProjectServiceModels.ProjectOutput)
	for _, v := range projects {
		go func(projectOutput *v2ProjectServiceModels.ProjectOutput) {
			project, err := s.projectRepo.GetCLAGroupsByExternalID(&v1ProjectParams.GetProjectsByExternalIDParams{
				ProjectSFID: projectOutput.ID,
				PageSize:    aws.Int64(1),
			}, DontLoadRepoDetails)
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

func fillCorporateContributorModel(wg *sync.WaitGroup, usersRepo users.UserRepository, sig *v1Models.Signature, result chan *models.CorporateContributor, searchTerm string) {
	defer wg.Done()
	user, err := usersRepo.GetUser(sig.SignatureReferenceID.String())
	if err != nil {
		log.Error("fillCorporateContributorModel: unable to get user info", err)
		return
	}
	if searchTerm != "" {
		ls := strings.ToLower(searchTerm)
		if !(strings.Contains(strings.ToLower(user.Username), ls) || strings.Contains(strings.ToLower(user.LfUsername), ls)) {
			return
		}
	}
	var contributor models.CorporateContributor
	var sigSignedTime = sig.SignatureCreated
	contributor.GithubID = user.GithubID
	contributor.LinuxFoundationID = user.LfUsername
	contributor.Name = user.Username
	t, err := utils.ParseDateTime(sig.SignatureCreated)
	if err != nil {
		log.Error("fillCorporateContributorModel: unable to parse time", err)
	} else {
		sigSignedTime = utils.TimeToString(t)
	}
	contributor.Timestamp = sigSignedTime
	contributor.SignatureVersion = fmt.Sprintf("v%s.%s", sig.SignatureMajorVersion, sig.SignatureMinorVersion)

	// send contributor struct on result channel
	result <- &contributor
}

func (s *service) getAllCompanyProjectEmployeeSignatures(companySFID string, projectSFID string) ([]*v1Models.Signature, error) {
	comp, claGroup, err := s.getCompanyAndClaGroup(companySFID, projectSFID)
	if err != nil {
		return nil, err
	}
	companyID := comp.CompanyID
	params := v1SignatureParams.GetProjectCompanyEmployeeSignaturesParams{
		HTTPRequest: nil,
		CompanyID:   companyID,
		ProjectID:   claGroup.ProjectID,
	}
	sigs, err := s.signatureRepo.GetProjectCompanyEmployeeSignatures(params, HugePageSize)
	if err != nil {
		return nil, err
	}
	return sigs.Signatures, nil
}

// get company and project parallely
func (s *service) getCompanyAndClaGroup(companySFID, projectSFID string) (*v1Models.Company, *v1Models.Project, error) {
	var comp *v1Models.Company
	var claGroup *v1Models.Project
	var companyErr, projectErr error
	// query projects and company
	var cp sync.WaitGroup
	cp.Add(2)
	go func() {
		defer cp.Done()
		comp, companyErr = s.companyRepo.GetCompanyByExternalID(companySFID)
	}()
	go func() {
		defer cp.Done()
		t := time.Now()
		var pm *projects_cla_groups.ProjectClaGroup
		pm, projectErr = s.projectClaGroupsRepo.GetClaGroupIDForProject(projectSFID)
		if projectErr != nil {
			log.Debugf("cla group mapping not found for projectSFID %s", projectSFID)
			return
		}
		claGroup, projectErr = s.projectRepo.GetCLAGroupByID(pm.ClaGroupID, DontLoadRepoDetails)
		if claGroup == nil {
			projectErr = ErrProjectNotFound
		}
		log.WithField("time_taken", time.Since(t).String()).Debugf("getting project by external id : %s completed", projectSFID)
	}()
	cp.Wait()
	if companyErr != nil {
		return nil, nil, companyErr
	}
	if projectErr != nil {
		return nil, nil, projectErr
	}
	return comp, claGroup, nil
}

// autoCreateCompany helper function to create a new company record based on the SF ID and underlying record in SF
func (s service) autoCreateCompany(companySFID string) (*v1Models.Company, error) {
	f := logrus.Fields{"companySFID": companySFID}
	// Get a reference to the platform organization service client
	orgClient := orgService.GetClient()
	log.WithFields(f).Debug("locating Organization in SF")

	// Lookup organization by ID in the Org Service
	sfOrgModel, sfOrgErr := orgClient.GetOrganization(companySFID)
	if sfOrgErr != nil {
		log.WithFields(f).Warnf("unable to locate platform organization record by SF ID, error: %+v", sfOrgErr)
		return nil, sfOrgErr
	}

	// If we were unable to lookup the company record in SF - we tried our best - return not exist error
	if sfOrgModel == nil {
		log.WithFields(f).Warn("unable to locate platform organization record by SF ID - record not found")
		return nil, company.ErrCompanyDoesNotExist
	}

	log.WithFields(f).Debug("found platform organization record in SF")
	// Auto-create based on the SF record information
	companyModel, companyCreateErr := s.companyRepo.CreateCompany(&v1Models.Company{
		CompanyExternalID: companySFID,
		CompanyName:       sfOrgModel.Name,
		Note:              "created on-demand by v4 service based on SF Organization Service record",
	})

	if companyCreateErr != nil || companyModel == nil {
		log.WithFields(f).Warnf("unable to create EasyCLA company from platform SF organization record, error: %+v",
			companyCreateErr)
		return nil, companyCreateErr
	}

	log.WithFields(f).Debugf("successfully created EasyCLA company record: %+v", companyModel)
	return companyModel, nil
}

func sendEmailToUserCompanyProfile(orgName string, userEmail string, username string, LFXPortalURL string) {
	subject := "EasyCLA: Company Profile "
	recipients := []string{userEmail}
	body := fmt.Sprintf(`
<p>Hello %s,</p>
<p>This is a notification email from EasyCLA regarding the newly created Salesforce Organization %s.</p>
<p>The organization profile can be completed via <a href="%s/company/manage/" target="_blank">clicking this link</a>
%s
%s`,
		username, orgName, LFXPortalURL,
		utils.GetEmailHelpContent(true), utils.GetEmailSignOffContent())
	err := utils.SendEmail(subject, body, recipients)
	if err != nil {
		log.Warnf("problem sending email with subject: %s to recipients: %+v, error: %+v", subject, recipients, err)
	} else {
		log.Debugf("sent email with subject: %s to recipients: %+v", subject, recipients)
	}
}
