// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

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

// SetInitialCLAManagerACSPermissions
func (s *service) SetInitialCLAManagerACSPermissions(signatureID string) error {
	f := logrus.Fields{
		"functionName": "SetInitialCLAManagerACSPermissions",
		"signatureID":  signatureID,
	}

	sig, err := s.signatureRepo.GetSignature(signatureID)
	if err != nil {
		log.WithFields(f).Warnf("problem locating signature by ID, error: %+v", err)
		return err
	}

	f["signatureType"] = sig.SignatureType
	if sig.SignatureType != CCLASignatureType {
		log.WithFields(f).Warn("invalid signature type for setting initial cla manager request")
		return fmt.Errorf("invalid signature type for setting initial cla manager request. %s", signatureID)
	}
	if len(sig.SignatureACL) == 0 {
		log.WithFields(f).Warn("initial cla manager details not found")
		return errors.New("initial cla manager details not found")
	}
	// get user details
	userServiceClient := v2UserService.GetClient()
	claManager, err := userServiceClient.GetUserByUsername(sig.SignatureACL[0].LfUsername)
	if err != nil {
		log.WithFields(f).Warnf("unable to lookup user by username: %s", sig.SignatureACL[0].LfUsername)
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
		log.WithFields(f).Warnf("unable to lookup company by ID: %s, error: %+v",
			sig.SignatureReferenceID.String(), err)
		return err
	}

	// fetch list of projects under cla group
	projectList, err := s.projectsClaGroupRepo.GetProjectsIdsForClaGroup(sig.ProjectID)
	if err != nil {
		log.WithFields(f).Warnf("unable to fetch list of projects associated with CLA Group: %s, error: %+v",
			sig.ProjectID, err)
		return err
	}

	// Assign cla manager role based on level
	err = s.assignCLAManager(email, claManager.Username, company.CompanyExternalID, projectList)
	if err != nil {
		log.WithFields(f).Warnf("unable to assign CLA Manager %s for company: %s, error: %+v",
			claManager.Username, company.CompanyExternalID, err)
		return err
	}

	return nil
}

func (s service) assignCLAManager(email, username, companySFID string, projectList []*projects_cla_groups.ProjectClaGroup) error {
	f := logrus.Fields{
		"functionName": "assignCLAManager",
		"email":        email,
		"username":     username,
		"companySFID":  companySFID,
	}

	// check if project is signed at foundation level
	foundationID := projectList[0].FoundationSFID
	signedAtFoundation, signedErr := s.projectService.SignedAtFoundationLevel(foundationID)
	if signedErr != nil {
		return signedErr
	}
	acsClient := v2AcsService.GetClient()
	claManagerRoleID, roleErr := acsClient.GetRoleID(utils.CLAManagerRole)
	if roleErr != nil {
		log.WithFields(f).Warnf("problem getting role for %s, error: %+v", utils.CLAManagerRole, roleErr)
		return roleErr
	}
	orgService := v2OrgService.GetClient()

	scopes, err := orgService.ListOrgUserScopes(companySFID, []string{utils.CLAManagerRole, utils.CLADesigneeRole})
	if err != nil {
		log.WithFields(f).Warnf("problem listing organization user scopes for company: %s, error: %+v",
			companySFID, err)
		return err
	}

	if signedAtFoundation {
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
				if roleScope.RoleName == utils.CLADesigneeRole {
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
				if roleScope.RoleName == utils.CLADesigneeRole {
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
