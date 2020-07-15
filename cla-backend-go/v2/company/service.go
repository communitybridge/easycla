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

	"github.com/communitybridge/easycla/cla-backend-go/projects_cla_groups"

	"github.com/jinzhu/copier"

	"github.com/LF-Engineering/lfx-kit/auth"
	"github.com/communitybridge/easycla/cla-backend-go/company"

	"github.com/aws/aws-sdk-go/aws"
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
)

// constants
const (
	// used when we want to query all data from dependent service.
	HugePageSize = int64(10000)
	// LoadRepoDetails     = true
	DontLoadRepoDetails = false
	FoundationType      = "Foundation"
	// ProjectType         = "Project"
)

// Service functions for company
type Service interface {
	GetCompanyProjectCLAManagers(companyID string, projectSFID string) (*models.CompanyClaManagers, error)
	GetCompanyProjectActiveCLAs(companyID string, projectSFID string) (*models.ActiveClaList, error)
	GetCompanyProjectContributors(projectSFID string, companySFID string, searchTerm string) (*models.CorporateContributorList, error)
	GetCompanyProjectCLA(authUser *auth.User, companySFID, projectSFID string) (*models.CompanyProjectClaList, error)
	CreateCompany(companyName string, companyWebsite string, userID string, LFXPortalURL string) (*models.CompanyOutput, error)
	GetCompanyByName(companyName string) (*models.Company, error)
	GetCompanyByID(companyID string) (*models.Company, error)
	GetCompanyBySFID(companySFID string) (*models.Company, error)
	DeleteCompanyByID(companyID string) error
	DeleteCompanyBySFID(companySFID string) error
	GetCompanyCLAGroupManagers(companyID, claGroupID string) (*models.CompanyClaManagers, error)
}

// ProjectRepo contains project repo methods
type ProjectRepo interface {
	GetProjectByID(projectID string, loadRepoDetails bool) (*v1Models.Project, error)
	GetProjectsByExternalID(params *v1ProjectParams.GetProjectsByExternalIDParams, loadRepoDetails bool) (*v1Models.Projects, error)
}

// NewService returns instance of company service
func NewService(sigRepo signatures.SignatureRepository, projectRepo ProjectRepo, usersRepo users.UserRepository, companyRepo company.IRepository, pcgRepo projects_cla_groups.Repository) Service {
	return &service{
		signatureRepo:        sigRepo,
		projectRepo:          projectRepo,
		userRepo:             usersRepo,
		companyRepo:          companyRepo,
		projectClaGroupsRepo: pcgRepo,
	}
}

func (s *service) GetCompanyProjectCLAManagers(companyID string, projectSFID string) (*models.CompanyClaManagers, error) {
	var err error
	claGroups, err := s.getCLAGroupsUnderProjectOrFoundation(projectSFID)
	if err != nil {
		return nil, err
	}
	sigs, err := s.getAllCCLASignatures(companyID)
	if err != nil {
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

func (s *service) GetCompanyProjectActiveCLAs(companyID string, projectSFID string) (*models.ActiveClaList, error) {
	var err error
	claGroups, err := s.getCLAGroupsUnderProjectOrFoundation(projectSFID)
	if err != nil {
		return nil, err
	}
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
	list := make([]*models.CorporateContributor, 0)
	sigs, err := s.getAllCompanyProjectEmployeeSignatures(companySFID, projectSFID)
	if err != nil {
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

func (s *service) CreateCompany(companyName string, companyWebsite string, userID string, LFXPortalURL string) (*models.CompanyOutput, error) {

	var lfUser *v2UserServiceModels.User

	// Get EasyCLA User
	claUser, claUserErr := s.userRepo.GetUser(userID)
	if claUserErr != nil {
		return nil, ErrCLAUserNotFound
	}

	// Create Sales Force company
	orgClient := orgService.GetClient()
	log.Debugf("Creating Organization : %s Website: %s", companyName, companyWebsite)
	org, err := orgClient.CreateOrg(companyName, companyWebsite)
	if err != nil {
		return nil, err
	}

	acsClient := acs_service.GetClient()
	userClient := v2UserService.GetClient()

	// Ensure to send user invite with public email (GH)
	var filteredEmails []string
	for _, email := range claUser.Emails {
		// Ignore noreply emails...
		if !(strings.Contains(email, "noreply.github.com")) {
			filteredEmails = append(filteredEmails, email)
		}
	}

	var userEmail *string

	if len(filteredEmails) > 0 {
		// Check if userEmail is an LF user
		for i := range filteredEmails {
			found, lfErr := userClient.SearchUserByEmail(filteredEmails[i])
			if lfErr != nil {
				userEmail = &filteredEmails[i]
				log.Debugf("User email :%s is not associated with LF", filteredEmails[i])
				continue
			}
			lfUser = found
			log.Debugf("User email : %s has an LFx account", filteredEmails[i])
			break
		}
	}

	if lfUser != nil || claUser.LfEmail != "" {
		// Set lf email
		var email string
		if claUser.LfEmail != "" {
			log.Debugf("Setting email to: %s", claUser.LfEmail)
			email = claUser.LfEmail
		} else if lfUser != nil && len(lfUser.Emails) > 0 {
			log.Debugf("Setting email to: %s", *lfUser.Emails[0].EmailAddress)
			email = *lfUser.Emails[0].EmailAddress
		}

		// set lf username
		var lfUsername string
		if claUser.LfUsername != "" {
			log.Debugf("Setting username to: %s", claUser.LfUsername)
			lfUsername = claUser.LfUsername
		} else if lfUser != nil && lfUser.Username != "" {
			log.Debugf("Setting username to: %s", lfUser.Username)
			lfUsername = lfUser.Username
		}

		// Get Role ID
		roleID, designeeErr := acsClient.GetRoleID("company-owner")
		if designeeErr != nil {
			msg := "Problem getting role ID for company-owner"
			log.Warn(msg)
			return nil, designeeErr
		}

		if lfUsername == "" {
			msg := fmt.Sprintf("EasyCLA - Bad Request - User with lfEmail: %s has no username required for setting scope ", email)
			log.Warn(msg)
			return nil, ErrNoLfUsername
		}
		err = orgClient.CreateOrgUserRoleOrgScope(email, org.ID, roleID)
		if err != nil {
			log.Warnf("Organization Service - Failed to assign company-owner role to user: %s, error: %+v ", email, err)
			return nil, err
		}
		log.Debugf("User :%s has been assigned the company-owner role to organization: %s ", email, org.Name)
		//Send Email to User with instructions to complete Company profile
		log.Debugf("Sending Email to user :%s to complete setup for newly created Org : %s ", email, org.Name)
		sendEmailToUserCompanyProfile(org.Name, email, lfUsername, LFXPortalURL)

	} else if userEmail != nil {
		// Send User invite
		log.Debugf("User invite initiated for user with email : %s", *userEmail)
		acsErr := acsClient.SendUserInvite(userEmail, "contributor", "organization", org.ID, "userinvite", nil, nil)
		if acsErr != nil {
			return nil, acsErr
		}
	}

	// Create Easy CLA Company
	log.Debugf("Creating EasyCLA company : %s ", companyName)
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
		log.Warnf("Failed to create EasyCLA company for company: %s ", companyName)
		return nil, createErr
	}

	return &models.CompanyOutput{
		CompanyName:    org.Name,
		CompanyWebsite: companyWebsite,
		LogoURL:        org.LogoURL,
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

// GetCompanyBySFID retrieves the company by external SFID
func (s *service) GetCompanyBySFID(companySFID string) (*models.Company, error) {
	companyModel, err := s.companyRepo.GetCompanyByExternalID(companySFID)
	if err != nil {
		return nil, err
	}

	if companyModel == nil {
		log.Debugf("search by company SFID: %s didn't locate the record", companySFID)
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

// DeleteCompanyByID deletes the company by ID
func (s *service) DeleteCompanyByID(companyID string) error {
	return s.companyRepo.DeleteCompanyByID(companyID)
}

// DeleteCompanyBySFID deletes the company by SFID
func (s *service) DeleteCompanyBySFID(companyID string) error {
	return s.companyRepo.DeleteCompanyBySFID(companyID)
}

func (s *service) GetCompanyProjectCLA(authUser *auth.User, companySFID, projectSFID string) (*models.CompanyProjectClaList, error) {
	var canSign bool
	resources := authUser.ResourceIDsByTypeAndRole(auth.ProjectOrganization, "cla-manager-designee")
	projectOrg := fmt.Sprintf("%s|%s", projectSFID, companySFID)
	for _, r := range resources {
		if r == projectOrg {
			canSign = true
			break
		}
	}
	companyModel, err := s.companyRepo.GetCompanyByExternalID(companySFID)
	if err != nil {
		return nil, err
	}

	claGroups, err := s.getCLAGroupsUnderProjectOrFoundation(projectSFID)
	if err != nil {
		return nil, err
	}

	activeCLAList, err := s.GetCompanyProjectActiveCLAs(companyModel.CompanyID, projectSFID)
	if err != nil {
		return nil, err
	}
	resp := &models.CompanyProjectClaList{
		SignedClaList:       activeCLAList.List,
		UnsignedProjectList: make([]*models.UnsignedProject, 0),
	}
	for _, activeCLA := range activeCLAList.List {
		// remove cla groups for which we have signed cla
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
		resp.UnsignedProjectList = append(resp.UnsignedProjectList, unsignedProject)
	}
	return resp, nil
}

// GetCompanyCLAGroupManagers when provided the internal company ID and CLA Groups ID, this routine returns the list of
// corresponding CLA managers
func (s *service) GetCompanyCLAGroupManagers(companyID, claGroupID string) (*models.CompanyClaManagers, error) {
	signed, approved := true, true
	pageSize := int64(10)
	sigModel, err := s.signatureRepo.GetProjectCompanySignature(companyID, claGroupID, &signed, &approved, nil, &pageSize)
	if err != nil {
		log.Warnf("unable to query CCLA signature using Company ID: %s and CLA Group ID: %s, signed: true, approved: true, error: %+v",
			companyID, claGroupID, err)
		return nil, err
	}

	if sigModel == nil {
		log.Warnf("unable to query CCLA signature using Company ID: %s and CLA Group ID: %s, signed: true, approved: true - no signature found",
			companyID, claGroupID)
		return nil, nil
	}

	projectModel, projErr := s.projectRepo.GetProjectByID(claGroupID, DontLoadRepoDetails)
	if projErr != nil {
		log.Warnf("unable to query CLA Group ID: %s, error: %+v", claGroupID, err)
		return nil, err
	}

	if projectModel == nil {
		log.Warnf("unable to query CLA Group ID: %s - no CLA Group found", claGroupID)
		return nil, nil
	}

	companyModel, companyErr := s.companyRepo.GetCompany(companyID)
	if companyErr != nil {
		log.Warnf("unable to query Company ID: %s, error: %+v", companyID, companyErr)
		return nil, err
	}

	if companyModel == nil {
		log.Warnf("unable to query Company ID: %s - no company by ID found", companyID)
		return nil, nil
	}

	claManagers := make([]*models.CompanyClaManager, 0)
	for _, user := range sigModel.SignatureACL {
		claManagers = append(claManagers, &models.CompanyClaManager{
			// DB doesn't have approved_on value - just use sig created date/time
			ApprovedOn:       sigModel.SignatureCreated,
			LfUsername:       user.LfUsername,
			Email:            user.LfEmail,
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
			cginfo, err := s.projectRepo.GetProjectByID(claGroupID, DontLoadRepoDetails)
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
			if email != nil && email.IsPrimary != nil && *email.IsPrimary {
				cm.Email = utils.StringValue(email.EmailAddress)
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
	activeCla.SignedOn = sig.SignatureCreated
	activeCla.ClaGroupName = cg.ClaGroupName
	activeCla.SignatureID = sig.SignatureID

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
		cclaURL, err = utils.GetDownloadLink(utils.SignedCLAFilename(sig.ProjectID, sig.SignatureType, sig.SignatureReferenceID, sig.SignatureID))
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
			project, err := s.projectRepo.GetProjectsByExternalID(&v1ProjectParams.GetProjectsByExternalIDParams{
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
	user, err := usersRepo.GetUser(sig.SignatureReferenceID)
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
		claGroup, projectErr = s.projectRepo.GetProjectByID(pm.ClaGroupID, DontLoadRepoDetails)
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

func sendEmailToUserCompanyProfile(orgName string, userEmail string, username string, LFXPortalURL string) {
	subject := "EasyCLA: Invitation to create LFID and complete process of becoming CLA Manager"
	recipients := []string{userEmail}
	body := fmt.Sprintf(`
<p>Hello %s,</p>
<p>This is a notification email from EasyCLA regarding the newly created Salesforce Organization %s .</p>
<p>The organization profile can be completed via <a href="%s/company/manage/" target="_blank">clicking this link</a>
<p>If you need help or have questions about EasyCLA, you can
<a href="https://docs.linuxfoundation.org/docs/communitybridge/communitybridge-easycla" target="_blank">read the documentation</a> or
<a href="https://jira.linuxfoundation.org/servicedesk/customer/portal/4/create/143" target="_blank">reach out to us for
support</a>.</p>

<p>Thanks,</p>
<p>EasyCLA support team</p>`,
		username, orgName, LFXPortalURL)
	err := utils.SendEmail(subject, body, recipients)
	if err != nil {
		log.Warnf("problem sending email with subject: %s to recipients: %+v, error: %+v", subject, recipients, err)
	} else {
		log.Debugf("sent email with subject: %s to recipients: %+v", subject, recipients)
	}
}
