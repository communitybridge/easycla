// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT
//nolint
package events

import (
	"fmt"
)

// EventData returns event data string which is used for event logging and containsPII field
type EventData interface {
	GetEventDetailsString(args *LogEventArgs) (eventData string, containsPII bool)
	GetEventSummaryString(args *LogEventArgs) (eventData string, containsPII bool)
}

// RepositoryAddedEventData . . .
type RepositoryAddedEventData struct {
	RepositoryName string
}

// RepositoryDisabledEventData . . .
type RepositoryDisabledEventData struct {
	RepositoryName string
}

// GerritProjectDeletedEventData . . .
type GerritProjectDeletedEventData struct {
	DeletedCount int
}

// GerritAddedEventData . . .
type GerritAddedEventData struct {
	GerritRepositoryName string
}

// GerritDeletedEventData . . .
type GerritDeletedEventData struct {
	GerritRepositoryName string
}

// GithubProjectDeletedEventData . . .
type GithubProjectDeletedEventData struct {
	DeletedCount int
}

// SignatureProjectInvalidatedEventData . . .
type SignatureProjectInvalidatedEventData struct {
	InvalidatedCount int
}

// UserCreatedEventData . . .
type UserCreatedEventData struct{}

// UserDeletedEventData . . .
type UserDeletedEventData struct {
	DeletedUserID string
}

// UserUpdatedEventData . . .
type UserUpdatedEventData struct{}

// CompanyACLRequestAddedEventData . . .
type CompanyACLRequestAddedEventData struct {
	UserName  string
	UserID    string
	UserEmail string
}

// CompanyACLRequestApprovedEventData . . .
type CompanyACLRequestApprovedEventData struct {
	UserName  string
	UserID    string
	UserEmail string
}

// CompanyACLRequestDeniedEventData . . .
type CompanyACLRequestDeniedEventData struct {
	UserName  string
	UserID    string
	UserEmail string
}

// CompanyACLUserAddedEventData . . .
type CompanyACLUserAddedEventData struct {
	UserLFID string
}

// CLATemplateCreatedEventData . . .
type CLATemplateCreatedEventData struct{}

// GithubOrganizationAddedEventData . . .
type GithubOrganizationAddedEventData struct {
	GithubOrganizationName  string
	AutoEnabled             bool
	AutoEnabledClaGroupID   string
	BranchProtectionEnabled bool
}

// GithubOrganizationDeletedEventData . . .
type GithubOrganizationDeletedEventData struct {
	GithubOrganizationName string
}

// GithubOrganizationUpdatedEventData . . .
type GithubOrganizationUpdatedEventData struct {
	GithubOrganizationName string
	AutoEnabled            bool
	AutoEnabledClaGroupID  string
}

// CCLAApprovalListRequestCreatedEventData . . .
type CCLAApprovalListRequestCreatedEventData struct {
	RequestID string
}

// CCLAApprovalListRequestApprovedEventData . . .
type CCLAApprovalListRequestApprovedEventData struct {
	RequestID string
}

// CCLAApprovalListRequestRejectedEventData . . .
type CCLAApprovalListRequestRejectedEventData struct {
	RequestID string
}

// CLAManagerCreatedEventData . . .
type CLAManagerCreatedEventData struct {
	CompanyName string
	ProjectName string
	UserName    string
	UserEmail   string
	UserLFID    string
}

// CLAManagerDeletedEventData . . .
type CLAManagerDeletedEventData struct {
	CompanyName string
	ProjectName string
	UserName    string
	UserEmail   string
	UserLFID    string
}

// CLAManagerRequestCreatedEventData . . .
type CLAManagerRequestCreatedEventData struct {
	RequestID   string
	CompanyName string
	ProjectName string
	UserName    string
	UserEmail   string
	UserLFID    string
}

// CLAManagerRequestApprovedEventData . . .
type CLAManagerRequestApprovedEventData struct {
	RequestID    string
	CompanyName  string
	ProjectName  string
	UserName     string
	UserEmail    string
	ManagerName  string
	ManagerEmail string
}

// CLAManagerRequestDeniedEventData . . .
type CLAManagerRequestDeniedEventData struct {
	RequestID    string
	CompanyName  string
	ProjectName  string
	UserName     string
	UserEmail    string
	ManagerName  string
	ManagerEmail string
}

// CLAManagerRequestDeletedEventData . . .
type CLAManagerRequestDeletedEventData struct {
	RequestID    string
	CompanyName  string
	ProjectName  string
	UserName     string
	UserEmail    string
	ManagerName  string
	ManagerEmail string
}

// CLAApprovalListAddEmailData . . .
type CLAApprovalListAddEmailData struct {
	UserName          string
	UserEmail         string
	UserLFID          string
	ApprovalListEmail string
}

// CLAApprovalListRemoveEmailData . . .
type CLAApprovalListRemoveEmailData struct {
	UserName          string
	UserEmail         string
	UserLFID          string
	ApprovalListEmail string
}

// CLAApprovalListAddDomainData . . .
type CLAApprovalListAddDomainData struct {
	UserName           string
	UserEmail          string
	UserLFID           string
	ApprovalListDomain string
}

// CLAApprovalListRemoveDomainData . . .
type CLAApprovalListRemoveDomainData struct {
	UserName           string
	UserEmail          string
	UserLFID           string
	ApprovalListDomain string
}

// CLAApprovalListAddGitHubUsernameData . . .
type CLAApprovalListAddGitHubUsernameData struct {
	UserName                   string
	UserEmail                  string
	UserLFID                   string
	ApprovalListGitHubUsername string
}

// CLAApprovalListRemoveGitHubUsernameData . . .
type CLAApprovalListRemoveGitHubUsernameData struct {
	UserName                   string
	UserEmail                  string
	UserLFID                   string
	ApprovalListGitHubUsername string
}

// CLAApprovalListAddGitHubOrgData . . .
type CLAApprovalListAddGitHubOrgData struct {
	UserName              string
	UserEmail             string
	UserLFID              string
	ApprovalListGitHubOrg string
}

// CLAApprovalListRemoveGitHubOrgData . . .
type CLAApprovalListRemoveGitHubOrgData struct {
	UserName              string
	UserEmail             string
	UserLFID              string
	ApprovalListGitHubOrg string
}

// ApprovalListGithubOrganizationAddedEventData . . .
type ApprovalListGithubOrganizationAddedEventData struct {
	GithubOrganizationName string
}

// ApprovalListGithubOrganizationDeletedEventData . . .
type ApprovalListGithubOrganizationDeletedEventData struct {
	GithubOrganizationName string
}

// ClaManagerAccessRequestAddedEventData . . .
type ClaManagerAccessRequestAddedEventData struct {
	ProjectName string
	CompanyName string
}

// ClaManagerAccessRequestDeletedEventData . . .
type ClaManagerAccessRequestDeletedEventData struct {
	RequestID string
}

// CLAGroupCreatedEventData . . .
type CLAGroupCreatedEventData struct{}

// CLAGroupUpdatedEventData . . .
type CLAGroupUpdatedEventData struct {
	ClaGroupName        string
	ClaGroupDescription string
}

// CLAGroupDeletedEventData . . .
type CLAGroupDeletedEventData struct{}

// ContributorNotifyCompanyAdminData . . .
type ContributorNotifyCompanyAdminData struct {
	AdminName  string
	AdminEmail string
}

// ContributorNotifyCLADesignee . . .
type ContributorNotifyCLADesignee struct {
	DesigneeName  string
	DesigneeEmail string
}

// ContributorAssignCLADesignee . . .
type ContributorAssignCLADesignee struct {
	DesigneeName  string
	DesigneeEmail string
}

// UserConvertToContactData . . .
type UserConvertToContactData struct{}

// AssignRoleScopeData . . .
type AssignRoleScopeData struct {
	Role  string
	Scope string
}

// ClaManagerRoleCreatedData . . .
type ClaManagerRoleCreatedData struct {
	Role      string
	Scope     string
	UserName  string
	UserEmail string
}

// ClaManagerRoleDeletedData . . .
type ClaManagerRoleDeletedData struct {
	Role      string
	Scope     string
	UserName  string
	UserEmail string
}

// GetEventDetailsString . . .
func (ed *RepositoryAddedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("Github Repository: %s added to Project: %s by: %s.", ed.RepositoryName, args.projectName, args.userName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *RepositoryDisabledEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("Github Repository: %s deleted from Project: %s by: %s.", ed.RepositoryName, args.projectName, args.userName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *UserCreatedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s added. User Details: %+v.", args.userName, args.UserModel)
	return data, true
}

// GetEventDetailsString . . .
func (ed *UserUpdatedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	return fmt.Sprintf("User: %s updated. User Details:  %+v.", args.userName, *args.UserModel), true
}

// GetEventDetailsString . . .
func (ed *UserDeletedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s deleted. User ID: %s.", args.userName, ed.DeletedUserID)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CompanyACLRequestAddedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s added pending invite with ID: %s and Email: %s for Company: %s.",
		ed.UserName, ed.UserID, ed.UserEmail, args.companyName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CompanyACLRequestApprovedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("Access Aproved for User: %s, ID: %s, Email: %s Company Group: %s.",
		ed.UserName, args.companyName, ed.UserID, ed.UserEmail)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CompanyACLRequestDeniedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("Access Denied for User: %s, ID: %s, Email: %s Company Group: %s.",
		ed.UserName, args.companyName, ed.UserID, ed.UserEmail)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CompanyACLUserAddedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User with LF Username: %s added to the ACL for Company: %s by: %s.",
		ed.UserLFID, args.companyName, args.userName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLATemplateCreatedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("PDF Templates created for Project: %s by: %s.", args.userName, args.projectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *GithubOrganizationAddedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("Github Organization: %s was added with auto-enabled: %t, with branch protection enabled: %t",
		ed.GithubOrganizationName, ed.AutoEnabled, ed.BranchProtectionEnabled)
	if ed.AutoEnabledClaGroupID != "" {
		data = data + fmt.Sprintf(" with auto-enabled-cla-group: %s", ed.AutoEnabledClaGroupID)
	}
	data = data + fmt.Sprintf(" by: %s.", args.userName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *GithubOrganizationDeletedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("Github Organization: %s was deleted by: %s.",
		ed.GithubOrganizationName, args.userName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *GithubOrganizationUpdatedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("Github Organization:%s was updated with auto-enabled: %t",
		ed.GithubOrganizationName, ed.AutoEnabled)
	if ed.AutoEnabledClaGroupID != "" {
		data = data + fmt.Sprintf(" with auto-enabled-cla-group: %s", ed.AutoEnabledClaGroupID)
	}
	data = data + fmt.Sprintf("by: %s.", args.userName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CCLAApprovalListRequestApprovedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s approved a CCLA Approval Request for Project: %s and Company: %s with Request ID: %s.",
		args.userName, args.projectName, args.companyName, ed.RequestID)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CCLAApprovalListRequestRejectedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s rejected a CCLA Approval Request for Project: %s, Company: %s with Request ID: %s.",
		args.userName, args.projectName, args.companyName, ed.RequestID)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAManagerRequestCreatedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s, LFID: %s, Email: %s added CLA Manager Request: %s for Company: %s, Project: %s.",
		ed.UserName, ed.UserLFID, ed.UserEmail, ed.RequestID, ed.CompanyName, ed.ProjectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAManagerCreatedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s, LFID: %s, Email: %s was added as CLA Manager for Company: %s, Project: %s.",
		ed.UserName, ed.UserLFID, ed.UserEmail, ed.CompanyName, ed.ProjectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAManagerDeletedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s LFID: %s, Email: %s was removed as CLA Manager for Company: %s, Project: %s.",
		ed.UserName, ed.UserLFID, ed.UserEmail, ed.CompanyName, ed.ProjectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAManagerRequestApprovedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager Request: %s was approved for User %s, Email: %s by Manager: %s, Email: %s for Company: %s, Project: %s.",
		ed.RequestID, ed.UserName, ed.UserEmail, ed.ManagerName, ed.ManagerEmail, ed.CompanyName, ed.ProjectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAManagerRequestDeniedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager Request: %s was denied for User %s, Email: %s by Manager: %s, Email: %s for Company: %s, Project: %s.",
		ed.RequestID, ed.UserName, ed.UserEmail, ed.ManagerName, ed.ManagerEmail, ed.CompanyName, ed.ProjectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAManagerRequestDeletedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager Request: %s was deleted for User %s, Email: %s by Manager: %s, Email: %s for Company: %s, Project: %s.",
		ed.RequestID, ed.UserName, ed.UserEmail, ed.ManagerName, ed.ManagerEmail, ed.CompanyName, ed.ProjectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAApprovalListAddEmailData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager: %s, Email: %s, LFID: %s added Email: %s to the approval list for Company: %s, Project: %s.",
		ed.UserName, ed.UserEmail, ed.UserLFID, ed.ApprovalListEmail, args.companyName, args.projectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAApprovalListRemoveEmailData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager: %s, Email: %s, LFID: %s removed Email: %s from the approval list for Company: %s, Project: %s.",
		ed.UserName, ed.UserEmail, ed.UserLFID, ed.ApprovalListEmail, args.companyName, args.projectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAApprovalListAddDomainData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager: %s, Email: %s, LFID: %s added Domain: %s to the approval list for Company: %s, Project: %s.",
		ed.UserName, ed.UserEmail, ed.UserLFID, ed.ApprovalListDomain, args.companyName, args.projectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAApprovalListRemoveDomainData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager: %s, Email: %s, LFID: %s removed Domain %s from the approval list for Company: %s, Project: %s.",
		ed.UserName, ed.UserEmail, ed.UserLFID, ed.ApprovalListDomain, args.companyName, args.projectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAApprovalListAddGitHubUsernameData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager: %s, Email: %s, LFID: %s added GitHub Username: %s to the approval list for Company: %s, Project: %s.",
		ed.UserName, ed.UserEmail, ed.UserLFID, ed.ApprovalListGitHubUsername, args.companyName, args.projectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAApprovalListRemoveGitHubUsernameData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager: %s, Email: %s, LFID: %s removed GitHub Username: %s from the approval list for Company: %s, Project: %s.",
		ed.UserName, ed.UserEmail, ed.UserLFID, ed.ApprovalListGitHubUsername, args.companyName, args.projectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAApprovalListAddGitHubOrgData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager: %s, Email: %s, LFID: %s added GitHub Organization: %s to the approval list for Company: %s, Project: %s.",
		ed.UserName, ed.UserEmail, ed.UserLFID, ed.ApprovalListGitHubOrg, args.companyName, args.projectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAApprovalListRemoveGitHubOrgData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager: %s, Email: %s, LFID: %s removed GitHub Organization: %s from the approval list for Company: %s, Project: %s.",
		ed.UserName, ed.UserEmail, ed.UserLFID, ed.ApprovalListGitHubOrg, args.companyName, args.projectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CCLAApprovalListRequestCreatedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s created a CCLA Approval Request for Project: %s, Company: %s with Request ID: %s.",
		args.userName, args.projectName, args.companyName, ed.RequestID)
	return data, true
}

// GetEventDetailsString . . .
func (ed *ApprovalListGithubOrganizationAddedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager: %s added GitHub Organization: %s to the whitelist for Company %s, Project: %s.",
		args.userName, ed.GithubOrganizationName, args.companyName, args.projectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *ApprovalListGithubOrganizationDeletedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager: %s removed GitHub Organization: %s from the whitelist for Company: %s, Project: %s.",
		args.userName, ed.GithubOrganizationName, args.companyName, args.projectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *ClaManagerAccessRequestAddedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s has requested to be CLA Manager for Company %s, Project: %s.",
		args.userName, ed.CompanyName, ed.ProjectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *ClaManagerAccessRequestDeletedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s has deleted CLA Manager Request with ID: %s.",
		args.userName, ed.RequestID)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAGroupCreatedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Group ID: %s, Name: %s was created by: %s.",
		args.ProjectID, args.projectName, args.userName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAGroupUpdatedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Group ID: %s was updated by: %s with Name: %s, Description: %s.",
		args.ProjectID, args.userName, ed.ClaGroupName, ed.ClaGroupDescription)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAGroupDeletedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Group ID: %s was deleted by: %s.",
		args.ProjectID, args.userName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *GerritProjectDeletedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("%d Gerrit Repositories were deleted due to CLA Group/Project: %s deletion.",
		ed.DeletedCount, args.projectName)
	return data, false
}

// GetEventDetailsString . . .
func (ed *GerritAddedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("Gerrit Repository: %s was added by: %s.", ed.GerritRepositoryName, args.userName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *GerritDeletedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("Gerrit Repository: %s was deleted by: %s.", ed.GerritRepositoryName, args.userName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *GithubProjectDeletedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("%d Github Repositories were deleted due to CLA Group/Project: [%s] deletion.",
		ed.DeletedCount, args.projectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *SignatureProjectInvalidatedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("%d Signatures were invalidated (approved set to false) due to CLA Group/Project: %s deletion.",
		ed.InvalidatedCount, args.projectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *ContributorNotifyCompanyAdminData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s notified Company Admin: %s by Email: %s for Company ID: %s, Name: %s.",
		args.userName, ed.AdminName, ed.AdminEmail, args.companyName, args.CompanyID)
	return data, true
}

// GetEventDetailsString . . .
func (ed *ContributorNotifyCLADesignee) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s notified CLA Designee: %s by Email: %s for Project Name : %s, ID: %s and Company Name: %s, ID: %s.",
		args.userName, ed.DesigneeName, ed.DesigneeEmail,
		args.projectName, args.ExternalProjectID,
		args.companyName, args.CompanyID)
	return data, true
}

// GetEventDetailsString . . .
func (ed *ContributorAssignCLADesignee) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User Name: %s, Email: %s was assigned as CLA Manager Designee for project Name: %s, ID:  %s and Company Name: %s, ID: %s by: %s.",
		ed.DesigneeName, ed.DesigneeEmail,
		args.projectName, args.ExternalProjectID,
		args.companyName, args.CompanyID, args.userName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *UserConvertToContactData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s was converted to Contact state for Project: %s.",
		args.LfUsername, args.ExternalProjectID)
	return data, true
}

// GetEventDetailsString . . .
func (ed *AssignRoleScopeData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s was assigned Scope: %s with Role: %s for Project: %s.",
		args.LfUsername,
		ed.Scope, ed.Role, args.ExternalProjectID)
	return data, true
}

// GetEventDetailsString . . .
func (ed *ClaManagerRoleCreatedData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s, Email: %s was added to Role: %s with Scope: %s by: %s.", ed.UserName, ed.UserEmail, ed.Role, ed.Scope, args.userName)
	return data, false
}

// GetEventDetailsString . . .
func (ed *ClaManagerRoleDeletedData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s, Email: %s was removed from Role: %s with Scope: %s by: %s.", ed.UserName, ed.UserEmail, ed.Role, ed.Scope, args.userName)
	return data, false
}

// Event Summary started

// GetEventSummaryString . . .
func (ed *RepositoryAddedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("Github Repository: %s was added to Project: %s by: %s.", ed.RepositoryName, args.projectName, args.userName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *RepositoryDisabledEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("Github Repository: %s was deleted from Project: %s by: %s.", ed.RepositoryName, args.projectName, args.userName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *UserCreatedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s was added, User Details: %+v.", args.userName, args.UserModel)
	return data, true
}

// GetEventSummaryString . . .
func (ed *UserUpdatedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	return fmt.Sprintf("User: %s was updated, User Details: %+v.", args.userName, *args.UserModel), true
}

// GetEventSummaryString . . .
func (ed *UserDeletedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User ID : %s was deleted by: %s.", ed.DeletedUserID, args.userName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CompanyACLRequestAddedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s with ID: %s, Email: %s requested Company Invite for Company: %s.",
		ed.UserName, ed.UserID, ed.UserEmail, args.companyName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CompanyACLRequestApprovedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("Company Invite was approved access for User: %s with ID: %s, Email: %s for company: %s.",
		ed.UserName, ed.UserID, ed.UserEmail, args.companyName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CompanyACLRequestDeniedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("Company Invite was denied access for User: %s with ID: %s, Email %s for Company: %s.",
		ed.UserName, ed.UserID, ed.UserEmail, args.companyName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CompanyACLUserAddedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User with LF Username %s was added to the ACL for Company: %s by: %s.",
		ed.UserLFID, args.companyName, args.userName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLATemplateCreatedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("PDF templates were created for Project %s by: %s.", args.projectName, args.userName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *GithubOrganizationAddedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("Github Organization: %s was added with auto-enabled: %t, branch protection enabled: %t",
		ed.GithubOrganizationName, ed.AutoEnabled, ed.BranchProtectionEnabled)
	if ed.AutoEnabledClaGroupID != "" {
		data = data + fmt.Sprintf(" with auto-enabled-cla-group: %s", ed.AutoEnabledClaGroupID)
	}
	data = data + fmt.Sprintf(" by: %s.", args.userName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *GithubOrganizationDeletedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("Github Organization: %s was deleted by: %s.",
		ed.GithubOrganizationName, args.userName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *GithubOrganizationUpdatedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("Github Organization: %s was updated with auto-enabled: %t",
		ed.GithubOrganizationName, ed.AutoEnabled)
	if ed.AutoEnabledClaGroupID != "" {
		data = data + fmt.Sprintf(" with auto-enabled-cla-group: %s", ed.AutoEnabledClaGroupID)
	}
	data = data + fmt.Sprintf(" by: %s.", args.userName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CCLAApprovalListRequestApprovedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s approved a CCLA Approval Request for Project: %s, Company: %s.",
		args.userName, args.projectName, args.companyName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CCLAApprovalListRequestRejectedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s rejected a CCLA Approval Request for Project: %s, Company: %s.",
		args.userName, args.projectName, args.companyName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAManagerRequestCreatedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s added CLA Manager Request: %s for Company: %s, Project: %s.",
		ed.UserName, ed.RequestID, ed.CompanyName, ed.ProjectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAManagerCreatedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s was added as CLA Manager for Company: %s, Project: %s.",
		ed.UserName, ed.CompanyName, ed.ProjectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAManagerDeletedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s was removed as CLA Manager for Company: %s, Project: %s.",
		ed.UserLFID, ed.CompanyName, ed.ProjectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAManagerRequestApprovedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager Request: %s for User: %s was approved by: %s for Company: %s, Project: %s.",
		ed.RequestID, ed.UserName, ed.ManagerName, ed.CompanyName, ed.ProjectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAManagerRequestDeniedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager Request: %s for User: %s was denied by: %s for Company: %s, Project: %s.",
		ed.RequestID, ed.UserName, ed.ManagerName, ed.CompanyName, ed.ProjectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAManagerRequestDeletedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager Request: %s for User: %s was deleted by: %s for Company: %s, Project: %s.",
		ed.RequestID, ed.UserName, ed.ManagerName, ed.CompanyName, ed.ProjectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAApprovalListAddEmailData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager: %s added Email: %s to the approval list for Company: %s, Project: %s.",
		ed.UserName, ed.ApprovalListEmail, args.companyName, args.projectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAApprovalListRemoveEmailData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager: %s removed Email: %s from the approval list for Company: %s, Project: %s.",
		ed.UserName, ed.ApprovalListEmail, args.companyName, args.projectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAApprovalListAddDomainData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager: %s added Domain: %s to the approval list for Company: %s, Project: %s.",
		ed.UserName, ed.ApprovalListDomain, args.companyName, args.projectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAApprovalListRemoveDomainData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager: %s removed Domain: %s from the approval list for Company: %s, Project: %s.",
		ed.UserName, ed.ApprovalListDomain, args.companyName, args.projectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAApprovalListAddGitHubUsernameData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager: %s added GitHub Username: %s to the approval list for Company: %s, Project: %s.",
		ed.UserName, ed.ApprovalListGitHubUsername, args.companyName, args.projectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAApprovalListRemoveGitHubUsernameData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager: %s removed GitHub Username: %s from the approval list for Company: %s, Project: %s.",
		ed.UserName, ed.ApprovalListGitHubUsername, args.companyName, args.projectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAApprovalListAddGitHubOrgData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager: %s added GitHub Organization: %s to the approval list for Company: %s, Project: %s.",
		ed.UserName, ed.ApprovalListGitHubOrg, args.companyName, args.projectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAApprovalListRemoveGitHubOrgData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager: %s removed GitHub Organization: %s from the approval list for Company: %s, Project: %s.",
		ed.UserName, ed.ApprovalListGitHubOrg, args.companyName, args.projectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CCLAApprovalListRequestCreatedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s created a CCLA Approval Request for Project: %s, Company: %s.",
		args.userName, args.projectName, args.companyName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *ApprovalListGithubOrganizationAddedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager: %s added GitHub Organization: %s to the whitelist for Project: %s, Company: %s.",
		args.userName, ed.GithubOrganizationName, args.projectName, args.companyName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *ApprovalListGithubOrganizationDeletedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager: %s removed GitHub Organization: %s from the whitelist for Project: %s, Company: %s.",
		args.userName, ed.GithubOrganizationName, args.projectName, args.companyName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *ClaManagerAccessRequestAddedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s has requested to be CLA Manager for Project: %s, Company: %s.",
		args.userName, ed.ProjectName, ed.CompanyName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *ClaManagerAccessRequestDeletedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s has deleted a request to be CLA Manager.",
		args.userName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAGroupCreatedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Group: %s was created by: %s.",
		args.projectName, args.userName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAGroupUpdatedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Group: %s was updated by: %s.", args.projectName, args.userName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAGroupDeletedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Group: %s was deleted by: %s.",
		args.projectName, args.userName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *GerritProjectDeletedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("%d Gerrit Repositories deleted  due to CLA Group/Project: %s deletion.",
		ed.DeletedCount, args.projectName)
	return data, false
}

// GetEventSummaryString . . .
func (ed *GerritAddedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("Gerrit Repository: %s added by: %s.", ed.GerritRepositoryName, args.userName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *GerritDeletedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("Gerrit Repository: %s was deleted by: %s.", ed.GerritRepositoryName, args.userName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *GithubProjectDeletedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("%d Github Repositories deleted due to CLA Group/Project: %s deletion.",
		ed.DeletedCount, args.projectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *SignatureProjectInvalidatedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("%d Signatures invalidated (approved set to false) due to CLA Group/Project: %s deletion.",
		ed.InvalidatedCount, args.projectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *ContributorNotifyCompanyAdminData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s notified Company Admin: %s by Email: %s for Company: %s.",
		args.userName, ed.AdminName, ed.AdminEmail, args.companyName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *ContributorNotifyCLADesignee) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s notified CLA Designee: %s by Email: %s for Project: %s, Company: %s.",
		args.userName, ed.DesigneeName, ed.DesigneeEmail,
		args.projectName, args.companyName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *ContributorAssignCLADesignee) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s was assigned as CLA Manager Designee for Project: %s, Company: %s by: %s.",
		ed.DesigneeName,
		args.projectName, args.companyName, args.userName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *UserConvertToContactData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s was converted to Contact state for Project: %s.",
		args.LfUsername, args.projectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *AssignRoleScopeData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User %s was added to Role: %s for Project: %s.", args.LfUsername, ed.Role, args.projectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *ClaManagerRoleCreatedData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s was added with Role: %s by: %s.", ed.UserName, ed.Role, args.userName)
	return data, false
}

// GetEventSummaryString . . .
func (ed *ClaManagerRoleDeletedData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("User: %s was removed from Role: %s by: %s.", ed.UserName, ed.Role, args.userName)
	return data, false
}
