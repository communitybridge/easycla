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

type GithubProjectDeletedEventData struct{}

type SignatureProjectInvalidatedEventData struct{}

type UserCreatedEventData struct{}
type UserDeletedEventData struct {
	DeletedUserID string
}
type UserUpdatedEventData struct{}

type CompanyACLRequestAddedEventData struct {
	UserID string
}

type CompanyACLRequestDeletedEventData struct {
	UserLFID string
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

type CCLAWhitelistRequestDeletedEventData struct {
	RequestID string
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

func (ed *CompanyACLRequestDeletedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] deleted pending invite for user with lf username [%s] for company: [%s]",
		args.userName, ed.UserLFID, args.companyName)
	return data, true
}
func (ed *CompanyACLRequestAddedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] added pending invite for user with user_id [%s] for company: [%s]",
		args.userName, ed.UserID, args.companyName)
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

func (ed *CCLAWhitelistRequestDeletedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] deleted a CCLA Whitelist Request for project: [%s], company: [%s] - request id: %s",
		args.userName, args.projectName, args.companyName, ed.RequestID)
	return data, true
}

func (ed *CCLAWhitelistRequestCreatedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] created a CCLA Whitelist Request for project: [%s], company: [%s] - request id: %s",
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
	return data, true
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
