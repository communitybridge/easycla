package dynamo_events

import (
	"errors"
	"fmt"
	"strings"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
	v2AcsService "github.com/communitybridge/easycla/cla-backend-go/v2/acs-service"
	v2OrgService "github.com/communitybridge/easycla/cla-backend-go/v2/organization-service"
	v2UserService "github.com/communitybridge/easycla/cla-backend-go/v2/user-service"
	"github.com/sirupsen/logrus"
)

// constants
const (
	RoleCLAManager         = "cla-manager"
	RoleCLAManagerDesignee = "cla-manager-designee"
)

func (s *service) SetInitialCLAManagerACSPermissions(signatureID string) error {
	sig, err := s.signatureRepo.GetSignature(signatureID)
	if err != nil {
		return err
	}
	if sig.SignatureType != CCLASignatureType {
		return fmt.Errorf("invalid signature set initial cla manager request. %s", signatureID)
	}
	if len(sig.SignatureACL) == 0 {
		return errors.New("initial cla manager details not found")
	}
	// get user details
	userServiceClient := v2UserService.GetClient()
	claManager, err := userServiceClient.GetUserByUsername(sig.SignatureACL[0].LfUsername)
	if err != nil {
		return err
	}
	// get email id from the user details
	var email string
	for _, e := range claManager.Emails {
		if e != nil && e.IsPrimary != nil && *e.IsPrimary {
			email = utils.StringValue(e.EmailAddress)
		}
	}
	company, err := s.companyRepo.GetCompany(sig.SignatureReferenceID)
	if err != nil {
		return err
	}
	// fetch list of projects under cla group
	projectList, err := s.projectsClaGroupRepo.GetProjectsIdsForClaGroup(sig.ProjectID)
	if err != nil {
		return err
	}
	projectSFIDList := utils.NewStringSet()
	for _, p := range projectList {
		projectSFIDList.Add(p.ProjectSFID)
	}
	acsClient := v2AcsService.GetClient()
	claManagerRoleID, roleErr := acsClient.GetRoleID(RoleCLAManager)
	if roleErr != nil {
		return roleErr
	}
	orgService := v2OrgService.GetClient()
	// add user as cla-manager for all projects of cla-group
	for _, projectSFID := range projectSFIDList.List() {
		err = orgService.CreateOrgUserRoleOrgScopeProjectOrg(email, projectSFID, company.CompanyExternalID, claManagerRoleID)
		if err != nil {
			log.WithFields(logrus.Fields{
				"org_id":       company.CompanyExternalID,
				"project_sfid": projectSFID,
				"lf_username":  claManager.Username,
				"email":        email,
			}).Warnf("unable to add cla-manager scope. error = %s", err)
		}
	}
	// delete all cla-manager designee for all project of cla-group
	scopes, err := orgService.ListOrgUserScopes(company.CompanyExternalID, []string{RoleCLAManager, RoleCLAManagerDesignee})
	if err != nil {
		return err
	}
	for _, userRole := range scopes.Userroles {
		for _, roleScope := range userRole.RoleScopes {
			if roleScope.RoleName == RoleCLAManagerDesignee {
				for _, scope := range roleScope.Scopes {
					tmp := strings.Split(scope.ObjectID, "|")
					projectSFID := tmp[0]
					if projectSFIDList.Include(projectSFID) {
						err = orgService.DeleteOrgUserRoleOrgScopeProjectOrg(company.CompanyExternalID, roleScope.RoleID, scope.ScopeID, nil, nil)
						if err != nil {
							log.WithFields(logrus.Fields{"org_id": company.CompanyExternalID, "object_name": scope.ObjectName}).Warnf("unable to delete stale cla-manager-designee scope")
						}
					}
				}
			}
		}
	}
	return nil
}
