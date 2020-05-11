//nolint
package events

import (
	"fmt"
)

// EventData returns event data string which is used for event logging and containsPII field
type EventData interface {
	GetEventString(args *LogEventArgs) (eventData string, containsPII bool)
}

type GithubRepositoryAddedEventData struct {
	RepositoryName string
}
type GithubRepositoryDeletedEventData struct {
	RepositoryName string
}

type GerritProjectDeletedEventData struct{}

type GerritAddedEventData struct {
	GerritRepositoryName string
}

type GerritDeletedEventData struct {
	GerritRepositoryName string
}

type GithubProjectDeletedEventData struct{}

type SignatureProjectInvalidatedEventData struct{}

type UserCreatedEventData struct{}
type UserDeletedEventData struct {
	DeletedUserID string
}
type UserUpdatedEventData struct{}

type CompanyACLRequestAddedEventData struct {
	UserName  string
	UserID    string
	UserEmail string
}

type CompanyACLRequestApprovedEventData struct {
	UserName  string
	UserID    string
	UserEmail string
}

type CompanyACLRequestDeniedEventData struct {
	UserName  string
	UserID    string
	UserEmail string
}

type CompanyACLUserAddedEventData struct {
	UserLFID string
}

type CLATemplateCreatedEventData struct{}

type GithubOrganizationAddedEventData struct {
	GithubOrganizationName string
}

type GithubOrganizationDeletedEventData struct {
	GithubOrganizationName string
}

type CCLAWhitelistRequestCreatedEventData struct {
	RequestID string
}

type CCLAWhitelistRequestApprovedEventData struct {
	RequestID string
}

type CCLAWhitelistRequestRejectedEventData struct {
	RequestID string
}

type CLAManagerCreatedEventData struct {
	CompanyName string
	ProjectName string
	UserName    string
	UserEmail   string
	UserLFID    string
}

type CLAManagerDeletedEventData struct {
	CompanyName string
	ProjectName string
	UserName    string
	UserEmail   string
	UserLFID    string
}

type CLAManagerRequestCreatedEventData struct {
	RequestID   string
	CompanyName string
	ProjectName string
	UserName    string
	UserEmail   string
	UserLFID    string
}

type CLAManagerRequestApprovedEventData struct {
	RequestID    string
	CompanyName  string
	ProjectName  string
	UserName     string
	UserEmail    string
	ManagerName  string
	ManagerEmail string
}

type CLAManagerRequestDeniedEventData struct {
	RequestID    string
	CompanyName  string
	ProjectName  string
	UserName     string
	UserEmail    string
	ManagerName  string
	ManagerEmail string
}

type CLAManagerRequestDeletedEventData struct {
	RequestID    string
	CompanyName  string
	ProjectName  string
	UserName     string
	UserEmail    string
	ManagerName  string
	ManagerEmail string
}

type WhitelistGithubOrganizationAddedEventData struct {
	GithubOrganizationName string
}
type WhitelistGithubOrganizationDeletedEventData struct {
	GithubOrganizationName string
}
type ClaManagerAccessRequestAddedEventData struct {
	ProjectName string
	CompanyName string
}
type ClaManagerAccessRequestDeletedEventData struct {
	RequestID string
}

type ProjectCreatedEventData struct{}
type ProjectUpdatedEventData struct{}
type ProjectDeletedEventData struct{}

func (ed *GithubRepositoryAddedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] added github repository [%s] to project [%s]", args.userName, ed.RepositoryName, args.projectName)
	return data, true
}

func (ed *GithubRepositoryDeletedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] deleted github repository [%s] from project [%s]", args.userName, ed.RepositoryName, args.projectName)
	return data, true
}

func (ed *UserCreatedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] added. user details = [%+v]", args.userName, args.UserModel)
	return data, true
}

func (ed *UserUpdatedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] updated. user details = [%+v]", args.userName, args.UserModel)
	return data, true
}

func (ed *UserDeletedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] deleted user id: [%s]", args.userName, ed.DeletedUserID)
	return data, true
}

func (ed *CompanyACLRequestAddedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] added pending invite with id [%s], email [%s] for company: [%s]",
		ed.UserName, ed.UserID, ed.UserEmail, args.companyName)
	return data, true
}

func (ed *CompanyACLRequestApprovedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] company invite was approved access with id [%s], email [%s] for company: [%s]",
		ed.UserName, ed.UserID, ed.UserEmail, args.companyName)
	return data, true
}

func (ed *CompanyACLRequestDeniedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] company invite was denied access with id [%s], email [%s] for company: [%s]",
		ed.UserName, ed.UserID, ed.UserEmail, args.companyName)
	return data, true
}

func (ed *CompanyACLUserAddedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] added user with lf username [%s] to the ACL for company: [%s]",
		args.userName, ed.UserLFID, args.companyName)
	return data, true
}

func (ed *CLATemplateCreatedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] created PDF templates for project [%s]", args.userName, args.projectName)
	return data, true
}

func (ed *GithubOrganizationAddedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] added github organization [%s]",
		args.userName, ed.GithubOrganizationName)
	return data, true
}

func (ed *GithubOrganizationDeletedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] deleted github organization [%s]",
		args.userName, ed.GithubOrganizationName)
	return data, true
}

func (ed *CCLAWhitelistRequestApprovedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] approved a CCLA Approval Request for project: [%s], company: [%s] - request id: %s",
		args.userName, args.projectName, args.companyName, ed.RequestID)
	return data, true
}

func (ed *CCLAWhitelistRequestRejectedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] rejected a CCLA Approval Request for project: [%s], company: [%s] - request id: %s",
		args.userName, args.projectName, args.companyName, ed.RequestID)
	return data, true
}

func (ed *CLAManagerRequestCreatedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s / %s / %s] added CLA Manager Request [%s] for Company: %s, Project: %s",
		ed.UserLFID, ed.UserName, ed.UserEmail, ed.RequestID, ed.CompanyName, ed.ProjectName)
	return data, true
}

func (ed *CLAManagerCreatedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s / %s / %s] was added as CLA Manager for Company: %s, Project: %s",
		ed.UserLFID, ed.UserName, ed.UserEmail, ed.CompanyName, ed.ProjectName)
	return data, true
}

func (ed *CLAManagerDeletedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s / %s / %s] was removed as CLA Manager for Company: %s, Project: %s",
		ed.UserLFID, ed.UserName, ed.UserEmail, ed.CompanyName, ed.ProjectName)
	return data, true
}

func (ed *CLAManagerRequestApprovedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager Request [%s] for user [%s / %s] was approved by [%s / %s] for Company: %s, Project: %s",
		ed.RequestID, ed.UserName, ed.UserEmail, ed.ManagerName, ed.ManagerEmail, ed.CompanyName, ed.ProjectName)
	return data, true
}

func (ed *CLAManagerRequestDeniedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager Request [%s] for user [%s / %s] was denied by [%s / %s] for Company: %s, Project: %s",
		ed.RequestID, ed.UserName, ed.UserEmail, ed.ManagerName, ed.ManagerEmail, ed.CompanyName, ed.ProjectName)
	return data, true
}

func (ed *CLAManagerRequestDeletedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager Request [%s] for user [%s / %s] was deleted by [%s / %s] for Company: %s, Project: %s",
		ed.RequestID, ed.UserName, ed.UserEmail, ed.ManagerName, ed.ManagerEmail, ed.CompanyName, ed.ProjectName)
	return data, true
}

func (ed *CCLAWhitelistRequestCreatedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] created a CCLA Approval Request for project: [%s], company: [%s] - request id: %s",
		args.userName, args.projectName, args.companyName, ed.RequestID)
	return data, true
}

func (ed *WhitelistGithubOrganizationAddedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager [%s] added GitHub Organization [%s] to the whitelist for project [%s] company [%s]",
		args.userName, ed.GithubOrganizationName, args.projectName, args.companyName)
	return data, true
}

func (ed *WhitelistGithubOrganizationDeletedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager [%s] removed GitHub Organization [%s] from the whitelist for project [%s] company [%s]",
		args.userName, ed.GithubOrganizationName, args.projectName, args.companyName)
	return data, true
}

func (ed *ClaManagerAccessRequestAddedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] has requested to be cla manager for project [%s] company [%s]",
		args.userName, ed.ProjectName, ed.CompanyName)
	return data, true
}

func (ed *ClaManagerAccessRequestDeletedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] has deleted request with id [%s] to be cla manager",
		args.userName, ed.RequestID)
	return data, true
}

func (ed *ProjectCreatedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] has created project [%s]",
		args.userName, args.projectName)
	return data, true
}

func (ed *ProjectUpdatedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] has updated project [%s]",
		args.userName, args.projectName)
	return data, true
}
func (ed *ProjectDeletedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] has deleted project [%s]",
		args.userName, args.projectName)
	return data, true
}

func (ed *GerritProjectDeletedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("Gerrit Repository Deleted  due to CLA Group/Project: [%s] deletion",
		args.projectName)
	containsPII := false
	return data, containsPII
}

func (ed *GerritAddedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] has added gerrit [%s]", args.userName, ed.GerritRepositoryName)
	containsPII := true
	return data, containsPII
}

func (ed *GerritDeletedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] has deleted gerrit [%s]", args.userName, ed.GerritRepositoryName)
	containsPII := true
	return data, containsPII
}

func (ed *GithubProjectDeletedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("Github Repository Deleted  due to CLA Group/Project: [%s] deletion",
		args.projectName)
	return data, true
}

func (ed *SignatureProjectInvalidatedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("Signature invalidated (approved set to false) due to CLA Group/Project: [%s] deletion",
		args.projectName)
	return data, true
}
