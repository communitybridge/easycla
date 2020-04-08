package events

import (
	"fmt"
)

type EventData interface {
	GetEventString(args *LogEventArgs) (eventData string, containsPII bool)
}

type GithubRepositoryAddedEventData struct {
	RepositoryName string
}
type GithubRepositoryDeletedEventData struct {
	RepositoryName string
}

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
	data := fmt.Sprintf("user [%s] added github repository [%s] to project [%s]", args.UserName, ed.RepositoryName, args.projectName)
	return data, true
}

func (ed *GithubRepositoryDeletedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] deleted github repository [%s] from project [%s]", args.UserName, ed.RepositoryName, args.projectName)
	return data, true
}

func (ed *UserCreatedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] added. user details = [%+v]", args.UserName, args.UserModel)
	return data, true
}

func (ed *UserUpdatedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] updated. user details = [%+v]", args.UserName, args.UserModel)
	return data, true
}

func (ed *UserDeletedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] deleted user id: [%s]", args.UserName, ed.DeletedUserID)
	return data, true
}

func (ed *CompanyACLRequestDeletedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] deleted pending invite for user with lf username [%s] for company: [%s]",
		args.UserName, ed.UserLFID, args.companyName)
	return data, true
}
func (ed *CompanyACLRequestAddedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] added pending invite for user with user_id [%s] for company: [%s]",
		args.UserName, ed.UserID, args.companyName)
	return data, true
}

func (ed *CompanyACLUserAddedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] added user with lf username [%s] to the ACL for company: [%s]",
		args.UserName, ed.UserLFID, args.companyName)
	return data, true
}

func (ed *CLATemplateCreatedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] created PDF templates for project [%s]", args.UserName, args.projectName)
	return data, true
}

func (ed *GithubOrganizationAddedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] added github organization [%s]",
		args.UserName, ed.GithubOrganizationName)
	return data, true
}

func (ed *GithubOrganizationDeletedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] deleted github organization [%s]",
		args.UserName, ed.GithubOrganizationName)
	return data, true
}

func (ed *CCLAWhitelistRequestDeletedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] deleted a CCLA Whitelist Request for project: [%s], company: [%s] - request id: %s",
		args.UserName, args.projectName, args.companyName, ed.RequestID)
	return data, true
}

func (ed *CCLAWhitelistRequestCreatedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] created a CCLA Whitelist Request for project: [%s], company: [%s] - request id: %s",
		args.UserName, args.projectName, args.companyName, ed.RequestID)
	return data, true
}

func (ed *WhitelistGithubOrganizationAddedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager [%s] added GitHub Organization [%s] to the whitelist for project [%s] company [%s]",
		args.UserName, ed.GithubOrganizationName, args.projectName, args.companyName)
	return data, true
}

func (ed *WhitelistGithubOrganizationDeletedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("CLA Manager [%s] removed GitHub Organization [%s] from the whitelist for project [%s] company [%s]",
		args.UserName, ed.GithubOrganizationName, args.projectName, args.companyName)
	return data, true
}

func (ed *ClaManagerAccessRequestAddedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] has requested to be cla manager for project [%s] company [%s]",
		args.UserName, ed.ProjectName, ed.CompanyName)
	return data, true
}

func (ed *ClaManagerAccessRequestDeletedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] has deleted request with id [%s] to be cla manager",
		args.UserName, ed.RequestID)
	return data, true
}

func (ed *ProjectCreatedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] has created project [%s]",
		args.UserName, args.projectName)
	return data, true
}

func (ed *ProjectUpdatedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] has updated project [%s]",
		args.UserName, args.projectName)
	return data, true
}
func (ed *ProjectDeletedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] has deleted project [%s]",
		args.UserName, args.projectName)
	return data, true
}
