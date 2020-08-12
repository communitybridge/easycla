package dynamo_events

import (
	"errors"
	"fmt"
	"strings"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/projects_cla_groups"
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
	company, err := s.companyRepo.GetCompany(sig.SignatureReferenceID.String())
	if err != nil {
		return err
	}

	// fetch list of projects under cla group
	projectList, err := s.projectsClaGroupRepo.GetProjectsIdsForClaGroup(sig.ProjectID)
	if err != nil {
		return err
	}

	// Assign cla manager role based on level
	err = s.assignCLAManager(email, claManager.Username, company.CompanyExternalID, projectList)
	if err != nil {
		return err
	}

	return nil
}

func (s service) assignCLAManager(email, username, companySFID string, projectList []*projects_cla_groups.ProjectClaGroup) error {
	// check if project is signed at foundation level
	signedAtFoundation := s.projectService.SignedAtFoundationLevel(projectList)
	acsClient := v2AcsService.GetClient()
	claManagerRoleID, roleErr := acsClient.GetRoleID(RoleCLAManager)
	if roleErr != nil {
		return roleErr
	}
	orgService := v2OrgService.GetClient()

	scopes, err := orgService.ListOrgUserScopes(companySFID, []string{RoleCLAManager, RoleCLAManagerDesignee})
	if err != nil {
		return err
	}

	if signedAtFoundation {
		foundationID := projectList[0].FoundationSFID
		// add cla manager role at foundation level
		err = orgService.CreateOrgUserRoleOrgScopeProjectOrg(email, foundationID, companySFID, claManagerRoleID)
		if err != nil {
			log.WithFields(logrus.Fields{
				"org_id":          companySFID,
				"foundation_sfid": foundationID,
				"lf_username":     username,
				"email":           email,
			}).Warnf("unable to add cla-manager scope. error = %s", err)
		}

		// delete cla manager designee role for foundation
		for _, userRole := range scopes.Userroles {
			for _, roleScope := range userRole.RoleScopes {
				if roleScope.RoleName == RoleCLAManagerDesignee {
					for _, scope := range roleScope.Scopes {
						tmp := strings.Split(scope.ObjectID, "|")
						projectSFID := tmp[0]
						if foundationID == projectSFID {
							err = orgService.DeleteOrgUserRoleOrgScopeProjectOrg(companySFID, roleScope.RoleID, scope.ScopeID, nil, nil)
							if err != nil {
								log.WithFields(logrus.Fields{"org_id": companySFID, "object_name": scope.ObjectName}).Warnf("unable to delete stale cla-manager-designee scope")
							}
						}
					}
				}
			}
		}

	} else {
		projectSFIDList := utils.NewStringSet()
		for _, p := range projectList {
			projectSFIDList.Add(p.ProjectSFID)
		}

		// add user as cla-manager for all projects of cla-group
		for _, projectSFID := range projectSFIDList.List() {
			err = orgService.CreateOrgUserRoleOrgScopeProjectOrg(email, projectSFID, companySFID, claManagerRoleID)
			if err != nil {
				log.WithFields(logrus.Fields{
					"org_id":       companySFID,
					"project_sfid": projectSFID,
					"lf_username":  username,
					"email":        email,
				}).Warnf("unable to add cla-manager scope. error = %s", err)
			}
		}
		// delete all cla-manager designee for all project of cla-group
		for _, userRole := range scopes.Userroles {
			for _, roleScope := range userRole.RoleScopes {
				if roleScope.RoleName == RoleCLAManagerDesignee {
					for _, scope := range roleScope.Scopes {
						tmp := strings.Split(scope.ObjectID, "|")
						projectSFID := tmp[0]
						if projectSFIDList.Include(projectSFID) {
							err = orgService.DeleteOrgUserRoleOrgScopeProjectOrg(companySFID, roleScope.RoleID, scope.ScopeID, nil, nil)
							if err != nil {
								log.WithFields(logrus.Fields{"org_id": companySFID, "object_name": scope.ObjectName}).Warnf("unable to delete stale cla-manager-designee scope")
							}
						}
					}
				}
			}
		}
	}

	return nil
}
