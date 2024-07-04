// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

// Code generated by MockGen. DO NOT EDIT.
// Source: company/repository.go

// Package mock_company is a generated GoMock package.
package mock_company

import (
	context "context"
	reflect "reflect"

	company "github.com/communitybridge/easycla/cla-backend-go/company"
	models "github.com/communitybridge/easycla/cla-backend-go/gen/v1/models"
	user "github.com/communitybridge/easycla/cla-backend-go/user"
	gomock "github.com/golang/mock/gomock"
)

// MockIRepository is a mock of IRepository interface.
type MockIRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIRepositoryMockRecorder
}

// MockIRepositoryMockRecorder is the mock recorder for MockIRepository.
type MockIRepositoryMockRecorder struct {
	mock *MockIRepository
}

// NewMockIRepository creates a new mock instance.
func NewMockIRepository(ctrl *gomock.Controller) *MockIRepository {
	mock := &MockIRepository{ctrl: ctrl}
	mock.recorder = &MockIRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIRepository) EXPECT() *MockIRepositoryMockRecorder {
	return m.recorder
}

// AddPendingCompanyInviteRequest mocks base method.
func (m *MockIRepository) AddPendingCompanyInviteRequest(ctx context.Context, companyID string, userModel user.User) (*company.Invite, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPendingCompanyInviteRequest", ctx, companyID, userModel)
	ret0, _ := ret[0].(*company.Invite)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddPendingCompanyInviteRequest indicates an expected call of AddPendingCompanyInviteRequest.
func (mr *MockIRepositoryMockRecorder) AddPendingCompanyInviteRequest(ctx, companyID, userModel interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPendingCompanyInviteRequest", reflect.TypeOf((*MockIRepository)(nil).AddPendingCompanyInviteRequest), ctx, companyID, userModel)
}

// ApproveCompanyAccessRequest mocks base method.
func (m *MockIRepository) ApproveCompanyAccessRequest(ctx context.Context, companyInviteID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ApproveCompanyAccessRequest", ctx, companyInviteID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ApproveCompanyAccessRequest indicates an expected call of ApproveCompanyAccessRequest.
func (mr *MockIRepositoryMockRecorder) ApproveCompanyAccessRequest(ctx, companyInviteID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApproveCompanyAccessRequest", reflect.TypeOf((*MockIRepository)(nil).ApproveCompanyAccessRequest), ctx, companyInviteID)
}

// CreateCompany mocks base method.
func (m *MockIRepository) CreateCompany(ctx context.Context, in *models.Company) (*models.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCompany", ctx, in)
	ret0, _ := ret[0].(*models.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCompany indicates an expected call of CreateCompany.
func (mr *MockIRepositoryMockRecorder) CreateCompany(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCompany", reflect.TypeOf((*MockIRepository)(nil).CreateCompany), ctx, in)
}

// DeleteCompanyByID mocks base method.
func (m *MockIRepository) DeleteCompanyByID(ctx context.Context, companyID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCompanyByID", ctx, companyID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCompanyByID indicates an expected call of DeleteCompanyByID.
func (mr *MockIRepositoryMockRecorder) DeleteCompanyByID(ctx, companyID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCompanyByID", reflect.TypeOf((*MockIRepository)(nil).DeleteCompanyByID), ctx, companyID)
}

// DeleteCompanyBySFID mocks base method.
func (m *MockIRepository) DeleteCompanyBySFID(ctx context.Context, companySFID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCompanyBySFID", ctx, companySFID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCompanyBySFID indicates an expected call of DeleteCompanyBySFID.
func (mr *MockIRepositoryMockRecorder) DeleteCompanyBySFID(ctx, companySFID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCompanyBySFID", reflect.TypeOf((*MockIRepository)(nil).DeleteCompanyBySFID), ctx, companySFID)
}

// GetCompanies mocks base method.
func (m *MockIRepository) GetCompanies(ctx context.Context) (*models.Companies, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompanies", ctx)
	ret0, _ := ret[0].(*models.Companies)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompanies indicates an expected call of GetCompanies.
func (mr *MockIRepositoryMockRecorder) GetCompanies(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompanies", reflect.TypeOf((*MockIRepository)(nil).GetCompanies), ctx)
}

// GetCompaniesByExternalID mocks base method.
func (m *MockIRepository) GetCompaniesByExternalID(ctx context.Context, companySFID string, includeChildCompanies bool) ([]*models.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompaniesByExternalID", ctx, companySFID, includeChildCompanies)
	ret0, _ := ret[0].([]*models.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompaniesByExternalID indicates an expected call of GetCompaniesByExternalID.
func (mr *MockIRepositoryMockRecorder) GetCompaniesByExternalID(ctx, companySFID, includeChildCompanies interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompaniesByExternalID", reflect.TypeOf((*MockIRepository)(nil).GetCompaniesByExternalID), ctx, companySFID, includeChildCompanies)
}

// GetCompaniesByUserManager mocks base method.
func (m *MockIRepository) GetCompaniesByUserManager(ctx context.Context, userID string, userModel user.User) (*models.Companies, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompaniesByUserManager", ctx, userID, userModel)
	ret0, _ := ret[0].(*models.Companies)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompaniesByUserManager indicates an expected call of GetCompaniesByUserManager.
func (mr *MockIRepositoryMockRecorder) GetCompaniesByUserManager(ctx, userID, userModel interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompaniesByUserManager", reflect.TypeOf((*MockIRepository)(nil).GetCompaniesByUserManager), ctx, userID, userModel)
}

// GetCompaniesByUserManagerWithInvites mocks base method.
func (m *MockIRepository) GetCompaniesByUserManagerWithInvites(ctx context.Context, userID string, userModel user.User) (*models.CompaniesWithInvites, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompaniesByUserManagerWithInvites", ctx, userID, userModel)
	ret0, _ := ret[0].(*models.CompaniesWithInvites)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompaniesByUserManagerWithInvites indicates an expected call of GetCompaniesByUserManagerWithInvites.
func (mr *MockIRepositoryMockRecorder) GetCompaniesByUserManagerWithInvites(ctx, userID, userModel interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompaniesByUserManagerWithInvites", reflect.TypeOf((*MockIRepository)(nil).GetCompaniesByUserManagerWithInvites), ctx, userID, userModel)
}

// GetCompany mocks base method.
func (m *MockIRepository) GetCompany(ctx context.Context, companyID string) (*models.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompany", ctx, companyID)
	ret0, _ := ret[0].(*models.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompany indicates an expected call of GetCompany.
func (mr *MockIRepositoryMockRecorder) GetCompany(ctx, companyID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompany", reflect.TypeOf((*MockIRepository)(nil).GetCompany), ctx, companyID)
}

// GetCompanyByExternalID mocks base method.
func (m *MockIRepository) GetCompanyByExternalID(ctx context.Context, companySFID string) (*models.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompanyByExternalID", ctx, companySFID)
	ret0, _ := ret[0].(*models.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompanyByExternalID indicates an expected call of GetCompanyByExternalID.
func (mr *MockIRepositoryMockRecorder) GetCompanyByExternalID(ctx, companySFID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompanyByExternalID", reflect.TypeOf((*MockIRepository)(nil).GetCompanyByExternalID), ctx, companySFID)
}

// GetCompanyByName mocks base method.
func (m *MockIRepository) GetCompanyByName(ctx context.Context, companyName string) (*models.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompanyByName", ctx, companyName)
	ret0, _ := ret[0].(*models.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompanyByName indicates an expected call of GetCompanyByName.
func (mr *MockIRepositoryMockRecorder) GetCompanyByName(ctx, companyName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompanyByName", reflect.TypeOf((*MockIRepository)(nil).GetCompanyByName), ctx, companyName)
}

// GetCompanyBySigningEntityName mocks base method.
func (m *MockIRepository) GetCompanyBySigningEntityName(ctx context.Context, signingEntityName string) (*models.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompanyBySigningEntityName", ctx, signingEntityName)
	ret0, _ := ret[0].(*models.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompanyBySigningEntityName indicates an expected call of GetCompanyBySigningEntityName.
func (mr *MockIRepositoryMockRecorder) GetCompanyBySigningEntityName(ctx, signingEntityName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompanyBySigningEntityName", reflect.TypeOf((*MockIRepository)(nil).GetCompanyBySigningEntityName), ctx, signingEntityName)
}

// GetCompanyInviteRequest mocks base method.
func (m *MockIRepository) GetCompanyInviteRequest(ctx context.Context, companyInviteID string) (*company.Invite, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompanyInviteRequest", ctx, companyInviteID)
	ret0, _ := ret[0].(*company.Invite)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompanyInviteRequest indicates an expected call of GetCompanyInviteRequest.
func (mr *MockIRepositoryMockRecorder) GetCompanyInviteRequest(ctx, companyInviteID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompanyInviteRequest", reflect.TypeOf((*MockIRepository)(nil).GetCompanyInviteRequest), ctx, companyInviteID)
}

// GetCompanyInviteRequests mocks base method.
func (m *MockIRepository) GetCompanyInviteRequests(ctx context.Context, companyID string, status *string) ([]company.Invite, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompanyInviteRequests", ctx, companyID, status)
	ret0, _ := ret[0].([]company.Invite)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompanyInviteRequests indicates an expected call of GetCompanyInviteRequests.
func (mr *MockIRepositoryMockRecorder) GetCompanyInviteRequests(ctx, companyID, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompanyInviteRequests", reflect.TypeOf((*MockIRepository)(nil).GetCompanyInviteRequests), ctx, companyID, status)
}

// GetCompanyUserInviteRequests mocks base method.
func (m *MockIRepository) GetCompanyUserInviteRequests(ctx context.Context, companyID, userID string) (*company.Invite, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompanyUserInviteRequests", ctx, companyID, userID)
	ret0, _ := ret[0].(*company.Invite)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompanyUserInviteRequests indicates an expected call of GetCompanyUserInviteRequests.
func (mr *MockIRepositoryMockRecorder) GetCompanyUserInviteRequests(ctx, companyID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompanyUserInviteRequests", reflect.TypeOf((*MockIRepository)(nil).GetCompanyUserInviteRequests), ctx, companyID, userID)
}

// GetUserInviteRequests mocks base method.
func (m *MockIRepository) GetUserInviteRequests(ctx context.Context, userID string) ([]company.Invite, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserInviteRequests", ctx, userID)
	ret0, _ := ret[0].([]company.Invite)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserInviteRequests indicates an expected call of GetUserInviteRequests.
func (mr *MockIRepositoryMockRecorder) GetUserInviteRequests(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserInviteRequests", reflect.TypeOf((*MockIRepository)(nil).GetUserInviteRequests), ctx, userID)
}

// IsCCLAEnabledForCompany mocks base method.
func (m *MockIRepository) IsCCLAEnabledForCompany(ctx context.Context, companyID string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsCCLAEnabledForCompany", ctx, companyID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsCCLAEnabledForCompany indicates an expected call of IsCCLAEnabledForCompany.
func (mr *MockIRepositoryMockRecorder) IsCCLAEnabledForCompany(ctx, companyID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsCCLAEnabledForCompany", reflect.TypeOf((*MockIRepository)(nil).IsCCLAEnabledForCompany), ctx, companyID)
}

// RejectCompanyAccessRequest mocks base method.
func (m *MockIRepository) RejectCompanyAccessRequest(ctx context.Context, companyInviteID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RejectCompanyAccessRequest", ctx, companyInviteID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RejectCompanyAccessRequest indicates an expected call of RejectCompanyAccessRequest.
func (mr *MockIRepositoryMockRecorder) RejectCompanyAccessRequest(ctx, companyInviteID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RejectCompanyAccessRequest", reflect.TypeOf((*MockIRepository)(nil).RejectCompanyAccessRequest), ctx, companyInviteID)
}

// SearchCompanyByName mocks base method.
func (m *MockIRepository) SearchCompanyByName(ctx context.Context, companyName, nextKey string) (*models.Companies, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchCompanyByName", ctx, companyName, nextKey)
	ret0, _ := ret[0].(*models.Companies)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchCompanyByName indicates an expected call of SearchCompanyByName.
func (mr *MockIRepositoryMockRecorder) SearchCompanyByName(ctx, companyName, nextKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchCompanyByName", reflect.TypeOf((*MockIRepository)(nil).SearchCompanyByName), ctx, companyName, nextKey)
}

// UpdateCompanyAccessList mocks base method.
func (m *MockIRepository) UpdateCompanyAccessList(ctx context.Context, companyID string, companyACL []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCompanyAccessList", ctx, companyID, companyACL)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCompanyAccessList indicates an expected call of UpdateCompanyAccessList.
func (mr *MockIRepositoryMockRecorder) UpdateCompanyAccessList(ctx, companyID, companyACL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCompanyAccessList", reflect.TypeOf((*MockIRepository)(nil).UpdateCompanyAccessList), ctx, companyID, companyACL)
}