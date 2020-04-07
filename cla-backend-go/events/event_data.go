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

type PendingInviteDeletedEventData struct {
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

func (ed *PendingInviteDeletedEventData) GetEventString(args *LogEventArgs) (string, bool) {
	data := fmt.Sprintf("user [%s] deleted pending invite for user with lf username [%s] for company: [%s]",
		args.UserName, ed.UserLFID, args.companyName)
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
