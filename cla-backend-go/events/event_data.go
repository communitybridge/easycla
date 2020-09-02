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
	GithubOrganizationName string
	AutoEnabled            bool
}

// GithubOrganizationDeletedEventData . . .
type GithubOrganizationDeletedEventData struct {
	GithubOrganizationName string
}

// GithubOrganizationUpdatedEventData . . .
type GithubOrganizationUpdatedEventData struct {
	GithubOrganizationName string
	AutoEnabled            bool
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
type CLAGroupUpdatedEventData struct{}

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
	data := fmt.Sprintf("user [%s] added github repository [%s] to project [%s]", args.userName, ed.RepositoryName, args.projectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *RepositoryDisabledEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] deleted github repository [%s] from project [%s]", args.userName, ed.RepositoryName, args.projectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *UserCreatedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] added. user details = [%+v]", args.userName, args.UserModel)
	return data, true
}

// GetEventDetailsString . . .
func (ed *UserUpdatedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	return fmt.Sprintf("user [%s] updated. user details = [%+v]", args.userName, *args.UserModel), true
}

// GetEventDetailsString . . .
func (ed *UserDeletedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] deleted user id: [%s]", args.userName, ed.DeletedUserID)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CompanyACLRequestAddedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] added pending invite with id [%s], email [%s] for company: [%s]",
		ed.UserName, ed.UserID, ed.UserEmail, args.companyName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CompanyACLRequestApprovedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] company invite was approved access with id [%s], email [%s] for company: [%s]",
		ed.UserName, ed.UserID, ed.UserEmail, args.companyName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CompanyACLRequestDeniedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] company invite was denied access with id [%s], email [%s] for company: [%s]",
		ed.UserName, ed.UserID, ed.UserEmail, args.companyName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CompanyACLUserAddedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] added user with lf username [%s] to the ACL for company: [%s]",
		args.userName, ed.UserLFID, args.companyName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLATemplateCreatedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] created PDF templates for project [%s]", args.userName, args.projectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *GithubOrganizationAddedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] added github organization [%s] with auto-enabled: %t",
		args.userName, ed.GithubOrganizationName, ed.AutoEnabled)
	return data, true
}

// GetEventDetailsString . . .
func (ed *GithubOrganizationDeletedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] deleted github organization [%s]",
		args.userName, ed.GithubOrganizationName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *GithubOrganizationUpdatedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] updated github organization [%s] with auto-enabled: %t",
		args.userName, ed.GithubOrganizationName, ed.AutoEnabled)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CCLAApprovalListRequestApprovedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] approved a CCLA Approval Request for project: [%s], company: [%s] - request id: %s",
		args.userName, args.projectName, args.companyName, ed.RequestID)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CCLAApprovalListRequestRejectedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] rejected a CCLA Approval Request for project: [%s], company: [%s] - request id: %s",
		args.userName, args.projectName, args.companyName, ed.RequestID)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAManagerRequestCreatedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s / %s / %s] added CLA Manager Request [%s] for Company: %s, Project: %s",
		ed.UserLFID, ed.UserName, ed.UserEmail, ed.RequestID, ed.CompanyName, ed.ProjectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAManagerCreatedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s / %s / %s] was added as CLA Manager for Company: %s, Project: %s",
		ed.UserLFID, ed.UserName, ed.UserEmail, ed.CompanyName, ed.ProjectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAManagerDeletedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s / %s / %s] was removed as CLA Manager for Company: %s, Project: %s",
		ed.UserLFID, ed.UserName, ed.UserEmail, ed.CompanyName, ed.ProjectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAManagerRequestApprovedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager Request [%s] for user [%s / %s] was approved by [%s / %s] for Company: %s, Project: %s",
		ed.RequestID, ed.UserName, ed.UserEmail, ed.ManagerName, ed.ManagerEmail, ed.CompanyName, ed.ProjectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAManagerRequestDeniedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager Request [%s] for user [%s / %s] was denied by [%s / %s] for Company: %s, Project: %s",
		ed.RequestID, ed.UserName, ed.UserEmail, ed.ManagerName, ed.ManagerEmail, ed.CompanyName, ed.ProjectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAManagerRequestDeletedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager Request [%s] for user [%s / %s] was deleted by [%s / %s] for Company: %s, Project: %s",
		ed.RequestID, ed.UserName, ed.UserEmail, ed.ManagerName, ed.ManagerEmail, ed.CompanyName, ed.ProjectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAApprovalListAddEmailData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager [%s / %s / %s] added Email %s to the approval list for Company: %s, Project: %s",
		ed.UserName, ed.UserEmail, ed.UserLFID, ed.ApprovalListEmail, args.companyName, args.projectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAApprovalListRemoveEmailData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager [%s / %s / %s] removed Email %s from the approval list for Company: %s, Project: %s",
		ed.UserName, ed.UserEmail, ed.UserLFID, ed.ApprovalListEmail, args.companyName, args.projectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAApprovalListAddDomainData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager [%s / %s / %s] added Domain %s to the approval list for Company: %s, Project: %s",
		ed.UserName, ed.UserEmail, ed.UserLFID, ed.ApprovalListDomain, args.companyName, args.projectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAApprovalListRemoveDomainData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager [%s / %s / %s] removed Domain %s from the approval list for Company: %s, Project: %s",
		ed.UserName, ed.UserEmail, ed.UserLFID, ed.ApprovalListDomain, args.companyName, args.projectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAApprovalListAddGitHubUsernameData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager [%s / %s / %s] added GitHub Username %s to the approval list for Company: %s, Project: %s",
		ed.UserName, ed.UserEmail, ed.UserLFID, ed.ApprovalListGitHubUsername, args.companyName, args.projectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAApprovalListRemoveGitHubUsernameData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager [%s / %s / %s] removed GitHub Username %s from the approval list for Company: %s, Project: %s",
		ed.UserName, ed.UserEmail, ed.UserLFID, ed.ApprovalListGitHubUsername, args.companyName, args.projectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAApprovalListAddGitHubOrgData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager [%s / %s / %s] added GitHub Org %s to the approval list for Company: %s, Project: %s",
		ed.UserName, ed.UserEmail, ed.UserLFID, ed.ApprovalListGitHubOrg, args.companyName, args.projectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAApprovalListRemoveGitHubOrgData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager [%s / %s / %s] removed GitHub Org %s from the approval list for Company: %s, Project: %s",
		ed.UserName, ed.UserEmail, ed.UserLFID, ed.ApprovalListGitHubOrg, args.companyName, args.projectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CCLAApprovalListRequestCreatedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] created a CCLA Approval Request for project: [%s], company: [%s] - request id: %s",
		args.userName, args.projectName, args.companyName, ed.RequestID)
	return data, true
}

// GetEventDetailsString . . .
func (ed *ApprovalListGithubOrganizationAddedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager [%s] added GitHub Organization [%s] to the whitelist for project [%s] company [%s]",
		args.userName, ed.GithubOrganizationName, args.projectName, args.companyName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *ApprovalListGithubOrganizationDeletedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager [%s] removed GitHub Organization [%s] from the whitelist for project [%s] company [%s]",
		args.userName, ed.GithubOrganizationName, args.projectName, args.companyName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *ClaManagerAccessRequestAddedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] has requested to be cla manager for project [%s] company [%s]",
		args.userName, ed.ProjectName, ed.CompanyName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *ClaManagerAccessRequestDeletedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] has deleted request with id [%s] to be cla manager",
		args.userName, ed.RequestID)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAGroupCreatedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] has created a CLA Group [%s - %s]",
		args.userName, args.projectName, args.ProjectID)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAGroupUpdatedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] has updated CLA Group [%s - %s]",
		args.userName, args.projectName, args.ProjectID)
	return data, true
}

// GetEventDetailsString . . .
func (ed *CLAGroupDeletedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] has deleted CLA Group [%s - %s]",
		args.userName, args.projectName, args.ProjectID)
	return data, true
}

// GetEventDetailsString . . .
func (ed *GerritProjectDeletedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("Deleted %d Gerrit Repositories due to CLA Group/Project: [%s] deletion",
		ed.DeletedCount, args.projectName)
	containsPII := false
	return data, containsPII
}

// GetEventDetailsString . . .
func (ed *GerritAddedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] has added gerrit [%s]", args.userName, ed.GerritRepositoryName)
	containsPII := true
	return data, containsPII
}

// GetEventDetailsString . . .
func (ed *GerritDeletedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] has deleted gerrit [%s]", args.userName, ed.GerritRepositoryName)
	containsPII := true
	return data, containsPII
}

// GetEventDetailsString . . .
func (ed *GithubProjectDeletedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("Deleted %d Github Repositories  due to CLA Group/Project: [%s] deletion",
		ed.DeletedCount, args.projectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *SignatureProjectInvalidatedEventData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("Invalidated %d signatures (approved set to false) due to CLA Group/Project: [%s] deletion",
		ed.InvalidatedCount, args.projectName)
	return data, true
}

// GetEventDetailsString . . .
func (ed *ContributorNotifyCompanyAdminData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] notified company admin by email: %s %s for company [%s / %s]",
		args.userName, ed.AdminName, ed.AdminEmail, args.companyName, args.CompanyID)
	return data, true
}

// GetEventDetailsString . . .
func (ed *ContributorNotifyCLADesignee) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] notified CLA Designee by email: %s %s for project [%s / %s] company [%s / %s]",
		args.userName, ed.DesigneeName, ed.DesigneeEmail,
		args.projectName, args.ExternalProjectID,
		args.companyName, args.CompanyID)
	return data, true
}

// GetEventDetailsString . . .
func (ed *ContributorAssignCLADesignee) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] assigned user: [%s / %s] as CLA Manager Designee for project [%s / %s] company [%s / %s]",
		args.userName, ed.DesigneeName, ed.DesigneeEmail,
		args.projectName, args.ExternalProjectID,
		args.companyName, args.CompanyID)
	return data, true
}

// GetEventDetailsString . . .
func (ed *UserConvertToContactData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] converted to Contact state for project [%s]",
		args.LfUsername, args.ExternalProjectID)
	return data, true
}

// GetEventDetailsString . . .
func (ed *AssignRoleScopeData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] assigned scope [%s] with role [%s] for project [%s]",
		args.LfUsername,
		ed.Scope, ed.Role, args.ExternalProjectID)
	return data, true
}

// GetEventDetailsString . . .
func (ed *ClaManagerRoleCreatedData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] added user %s/%s from role: %s with scope: %s", args.userName, ed.UserName, ed.UserEmail, ed.Role, ed.Scope)
	return data, false
}

// GetEventDetailsString . . .
func (ed *ClaManagerRoleDeletedData) GetEventDetailsString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] removed user %s/%s from role: %s with scope: %s", args.userName, ed.UserName, ed.UserEmail, ed.Role, ed.Scope)
	return data, false
}

// Event Summary started

// GetEventSummaryString . . .
func (ed *RepositoryAddedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s added github repository %s to project %s", args.userName, ed.RepositoryName, args.projectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *RepositoryDisabledEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s deleted github repository %s from project %s", args.userName, ed.RepositoryName, args.projectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *UserCreatedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s added. user details %+v", args.userName, args.UserModel)
	return data, true
}

// GetEventSummaryString . . .
func (ed *UserUpdatedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	return fmt.Sprintf("user %s updated. user details %+v", args.userName, *args.UserModel), true
}

// GetEventSummaryString . . .
func (ed *UserDeletedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s deleted user id: %s", args.userName, ed.DeletedUserID)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CompanyACLRequestAddedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s added pending invite with id %s, email %s for company: %s",
		ed.UserName, ed.UserID, ed.UserEmail, args.companyName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CompanyACLRequestApprovedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s company invite was approved access with id %s, email %s for company: %s",
		ed.UserName, ed.UserID, ed.UserEmail, args.companyName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CompanyACLRequestDeniedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s company invite was denied access with id %s, email %s for company: %s",
		ed.UserName, ed.UserID, ed.UserEmail, args.companyName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CompanyACLUserAddedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s added user with lf username %s to the ACL for company: %s",
		args.userName, ed.UserLFID, args.companyName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLATemplateCreatedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s created PDF templates for project %s", args.userName, args.projectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *GithubOrganizationAddedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s added github organization %s with auto-enabled: %t",
		args.userName, ed.GithubOrganizationName, ed.AutoEnabled)
	return data, true
}

// GetEventSummaryString . . .
func (ed *GithubOrganizationDeletedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s deleted github organization %s",
		args.userName, ed.GithubOrganizationName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *GithubOrganizationUpdatedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s deleted github organization %s with auto-enabled: %t",
		args.userName, ed.GithubOrganizationName, ed.AutoEnabled)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CCLAApprovalListRequestApprovedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s approved a CCLA Approval Request for project: %s, company: %s - request id: %s",
		args.userName, args.projectName, args.companyName, ed.RequestID)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CCLAApprovalListRequestRejectedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s rejected a CCLA Approval Request for project: %s, company: %s - request id: %s",
		args.userName, args.projectName, args.companyName, ed.RequestID)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAManagerRequestCreatedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s added CLA Manager Request %s for Company: %s, Project: %s",
		ed.UserName, ed.RequestID, ed.CompanyName, ed.ProjectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAManagerCreatedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s was added as CLA Manager for Company: %s, Project: %s",
		ed.UserName, ed.CompanyName, ed.ProjectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAManagerDeletedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s was removed as CLA Manager for Company: %s, Project: %s",
		ed.UserLFID, ed.CompanyName, ed.ProjectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAManagerRequestApprovedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager Request %s for user %s was approved by %s for Company: %s, Project: %s",
		ed.RequestID, ed.UserName, ed.ManagerName, ed.CompanyName, ed.ProjectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAManagerRequestDeniedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager Request %s for user %s was denied by %s for Company: %s, Project: %s",
		ed.RequestID, ed.UserName, ed.ManagerName, ed.CompanyName, ed.ProjectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAManagerRequestDeletedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager Request %s for user %s was deleted by %s for Company: %s, Project: %s",
		ed.RequestID, ed.UserName, ed.ManagerName, ed.CompanyName, ed.ProjectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAApprovalListAddEmailData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager %s added Email %s to the approval list for Company: %s, Project: %s",
		ed.UserName, ed.ApprovalListEmail, args.companyName, args.projectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAApprovalListRemoveEmailData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager %s removed Email %s from the approval list for Company: %s, Project: %s",
		ed.UserName, ed.ApprovalListEmail, args.companyName, args.projectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAApprovalListAddDomainData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager %s added Domain %s to the approval list for Company: %s, Project: %s",
		ed.UserName, ed.ApprovalListDomain, args.companyName, args.projectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAApprovalListRemoveDomainData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager %s removed Domain %s from the approval list for Company: %s, Project: %s",
		ed.UserName, ed.ApprovalListDomain, args.companyName, args.projectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAApprovalListAddGitHubUsernameData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager %s added GitHub Username %s to the approval list for Company: %s, Project: %s",
		ed.UserName, ed.ApprovalListGitHubUsername, args.companyName, args.projectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAApprovalListRemoveGitHubUsernameData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager %s removed GitHub Username %s from the approval list for Company: %s, Project: %s",
		ed.UserName, ed.ApprovalListGitHubUsername, args.companyName, args.projectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAApprovalListAddGitHubOrgData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager %s added GitHub Org %s to the approval list for Company: %s, Project: %s",
		ed.UserName, ed.ApprovalListGitHubOrg, args.companyName, args.projectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAApprovalListRemoveGitHubOrgData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager %s removed GitHub Org %s from the approval list for Company: %s, Project: %s",
		ed.UserName, ed.ApprovalListGitHubOrg, args.companyName, args.projectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CCLAApprovalListRequestCreatedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s created a CCLA Approval Request for project: %s, company: %s - request id: %s",
		args.userName, args.projectName, args.companyName, ed.RequestID)
	return data, true
}

// GetEventSummaryString . . .
func (ed *ApprovalListGithubOrganizationAddedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager %s added GitHub Organization %s to the whitelist for project %s company %s",
		args.userName, ed.GithubOrganizationName, args.projectName, args.companyName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *ApprovalListGithubOrganizationDeletedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager %s removed GitHub Organization %s from the whitelist for project %s company %s",
		args.userName, ed.GithubOrganizationName, args.projectName, args.companyName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *ClaManagerAccessRequestAddedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s has requested to be cla manager for project %s company %s",
		args.userName, ed.ProjectName, ed.CompanyName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *ClaManagerAccessRequestDeletedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s has deleted request with id %s to be cla manager",
		args.userName, ed.RequestID)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAGroupCreatedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s has created a CLA Group %s - %s",
		args.userName, args.projectName, args.ProjectID)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAGroupUpdatedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s has updated CLA Group %s - %s",
		args.userName, args.projectName, args.ProjectID)
	return data, true
}

// GetEventSummaryString . . .
func (ed *CLAGroupDeletedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s has deleted CLA Group %s - %s",
		args.userName, args.projectName, args.ProjectID)
	return data, true
}

// GetEventSummaryString . . .
func (ed *GerritProjectDeletedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("Deleted %d Gerrit Repositories due to CLA Group/Project: %s deletion",
		ed.DeletedCount, args.projectName)
	containsPII := false
	return data, containsPII
}

// GetEventSummaryString . . .
func (ed *GerritAddedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s has added gerrit %s", args.userName, ed.GerritRepositoryName)
	containsPII := true
	return data, containsPII
}

// GetEventSummaryString . . .
func (ed *GerritDeletedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s has deleted gerrit %s", args.userName, ed.GerritRepositoryName)
	containsPII := true
	return data, containsPII
}

// GetEventSummaryString . . .
func (ed *GithubProjectDeletedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("Deleted %d Github Repositories  due to CLA Group/Project: %s deletion",
		ed.DeletedCount, args.projectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *SignatureProjectInvalidatedEventData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("Invalidated %d signatures (approved set to false) due to CLA Group/Project: %s deletion",
		ed.InvalidatedCount, args.projectName)
	return data, true
}

// GetEventSummaryString . . .
func (ed *ContributorNotifyCompanyAdminData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s notified company admin by email: %s %s for company %s / %s",
		args.userName, ed.AdminName, ed.AdminEmail, args.companyName, args.CompanyID)
	return data, true
}

// GetEventSummaryString . . .
func (ed *ContributorNotifyCLADesignee) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s notified CLA Designee by email: %s %s for project %s / %s company %s / %s",
		args.userName, ed.DesigneeName, ed.DesigneeEmail,
		args.projectName, args.ExternalProjectID,
		args.companyName, args.CompanyID)
	return data, true
}

// GetEventSummaryString . . .
func (ed *ContributorAssignCLADesignee) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s assigned user: %s as CLA Manager Designee for project %s / %s company %s / %s",
		args.userName, ed.DesigneeName,
		args.projectName, args.ExternalProjectID,
		args.companyName, args.CompanyID)
	return data, true
}

// GetEventSummaryString . . .
func (ed *UserConvertToContactData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s converted to Contact state for project %s",
		args.LfUsername, args.ExternalProjectID)
	return data, true
}

// GetEventSummaryString . . .
func (ed *AssignRoleScopeData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s assigned scope %s with role %s for project %s",
		args.LfUsername,
		ed.Scope, ed.Role, args.ExternalProjectID)
	return data, true
}

// GetEventSummaryString . . .
func (ed *ClaManagerRoleCreatedData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s added user %s from role: %s with scope: %s", args.userName, ed.UserName, ed.Role, ed.Scope)
	return data, false
}

// GetEventSummaryString . . .
func (ed *ClaManagerRoleDeletedData) GetEventSummaryString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user %s removed user %s from role: %s with scope: %s", args.userName, ed.UserName, ed.Role, ed.Scope)
	return data, false
}
