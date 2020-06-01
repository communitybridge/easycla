// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package service_discovery

import (
	"github.com/communitybridge/easycla/cla-backend-go/company"
	"github.com/communitybridge/easycla/cla-backend-go/events"
	"github.com/communitybridge/easycla/cla-backend-go/project"
	"github.com/communitybridge/easycla/cla-backend-go/signatures"
	"github.com/communitybridge/easycla/cla-backend-go/users"
	v2CompanyService "github.com/communitybridge/easycla/cla-backend-go/v2/company"
	v2ProjectService "github.com/communitybridge/easycla/cla-backend-go/v2/project"
)

// ServiceDiscovery interface
type ServiceDiscovery interface {
	SetEventService(eventService events.Service)
	GetEventService() events.Service
	SetUserService(userService users.Service)
	GetUserService() users.Service
	SetV1CompanyService(companyService company.IService)
	GetV1CompanyService() company.IService
	SetV2CompanyService(companyService v2CompanyService.Service)
	GetV2CompanyService() v2CompanyService.Service
	GetV1ProjectService() project.Service
	GetV2ProjectService() v2ProjectService.Service
	GetV1SignatureService() signatures.SignatureService
}

type service struct {
	eventService       events.Service
	userService        users.Service
	v1CompanyService   company.IService
	v2CompanyService   v2CompanyService.Service
	v1ProjectService   project.Service
	v2ProjectService   v2ProjectService.Service
	v1SignatureService signatures.SignatureService
}

// NewServiceDiscovery creates a new whitelist service
func NewServiceDiscovery() ServiceDiscovery {
	return service{}
}

// GetEventService sets the event service reference
func (s service) SetEventService(eventService events.Service) {
	s.eventService = eventService
}

// GetEventService returns a reference to the the event service
func (s service) GetEventService() events.Service {
	return s.eventService
}

// SetUserService sets the user service reference
func (s service) SetUserService(userService users.Service) {
	s.userService = userService
}

// GetUserService returns a reference to the the user service
func (s service) GetUserService() users.Service {
	return s.userService
}

// SetV1CompanyService sets the v1 company service reference
func (s service) SetV1CompanyService(companyService company.IService) {
	s.v1CompanyService = companyService
}

// GetV1CompanyService returns a reference to the the v1 company service
func (s service) GetV1CompanyService() company.IService {
	return s.v1CompanyService
}

// SetV2CompanyService sets the v1 company service reference
func (s service) SetV2CompanyService(companyService v2CompanyService.Service) {
	s.v2CompanyService = companyService
}

// GetV2CompanyService returns a reference to the the v2 company service
func (s service) GetV2CompanyService() v2CompanyService.Service {
	return s.v2CompanyService
}

// GetV1ProjectService returns a reference to the the v1 project service
func (s service) GetV1ProjectService() project.Service { return s.v1ProjectService }

// GetV2ProjectService returns a reference to the the v2 project service
func (s service) GetV2ProjectService() v2ProjectService.Service { return s.v2ProjectService }

// GetV1SignatureService returns a reference to the the v1 signature service
func (s service) GetV1SignatureService() signatures.SignatureService { return s.v1SignatureService }
