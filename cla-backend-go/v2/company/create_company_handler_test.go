// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package company

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LF-Engineering/lfx-kit/auth"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/linuxfoundation/easycla/cla-backend-go/gen/v2/models"
	"github.com/linuxfoundation/easycla/cla-backend-go/gen/v2/restapi/operations"
	v2CompanyOps "github.com/linuxfoundation/easycla/cla-backend-go/gen/v2/restapi/operations/company"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testCompanyName = "Acme Inc"
	testCompanySite = "https://acme.example.invalid"
)

type fakeCreateCompanyService struct {
	Service
	calls    int
	seen     string
	seenUser string
	err      error
}

func (f *fakeCreateCompanyService) CreateCompany(_ context.Context, params *v2CompanyOps.CreateCompanyParams) (*models.CompanyOutput, error) {
	f.calls++
	f.seen = params.Input.UserEmail.String()
	if params.XUSERNAME != nil {
		f.seenUser = *params.XUSERNAME
	}
	if f.err != nil {
		return nil, f.err
	}
	return &models.CompanyOutput{CompanyID: "0014100000Created", CompanyName: *params.Input.CompanyName}, nil
}

// the company can only be created for the signed-in user (lfx-self-serve-ops#343)
func TestCreateCompanyHandlerBindsTheCaller(t *testing.T) {
	var lastService *fakeCreateCompanyService
	call := func(t *testing.T, authUser *auth.User, bodyEmail string, headerUser *string) (int, int) {
		api := operations.NewEasyclaAPI(nil)
		service := &fakeCreateCompanyService{}
		lastService = service
		Configure(api, service, nil, "")
		require.NotNil(t, api.CompanyCreateCompanyHandler)
		email := strfmt.Email(bodyEmail)
		name, site := testCompanyName, testCompanySite
		responder := api.CompanyCreateCompanyHandler.Handle(v2CompanyOps.CreateCompanyParams{
			HTTPRequest: httptest.NewRequest(http.MethodPost, "/v4/user/user-1/company", nil),
			UserID:      "user-1",
			XUSERNAME:   headerUser,
			Input:       &models.CompanyInput{CompanyName: &name, CompanyWebsite: &site, UserEmail: &email},
		}, authUser)
		recorder := httptest.NewRecorder()
		responder.WriteResponse(recorder, runtime.JSONProducer())
		return recorder.Code, service.calls
	}
	caller := func() *auth.User {
		return &auth.User{UserName: "jane", Email: "jane@example.invalid", ACL: auth.ACL{Allowed: true}}
	}
	status, calls := call(t, caller(), "mallory@example.invalid", nil)
	assert.Equal(t, http.StatusForbidden, status)
	assert.Equal(t, 0, calls, "no organization or role may be created for another person's email")
	status, calls = call(t, caller(), "Jane@Example.invalid", nil)
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, 1, calls)
	assert.Equal(t, "jane", lastService.seenUser, "the service receives the principal's username even without an X-USERNAME header")
	principal := caller()
	other := "mallory"
	status, calls = call(t, principal, "Jane@Example.invalid", &other)
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, 1, calls)
	assert.Equal(t, principal.UserName, lastService.seenUser, "the username the service sees is the authenticated principal's, not a free header value")
	status, calls = call(t, &auth.User{UserName: "ghost", ACL: auth.ACL{Allowed: true}}, "ghost@example.invalid", nil)
	assert.Equal(t, http.StatusForbidden, status, "a principal without an email cannot be bound")
	assert.Equal(t, 0, calls)
	status, calls = call(t, nil, "nobody@example.invalid", nil)
	assert.Equal(t, http.StatusForbidden, status, "a missing principal is refused, not dereferenced")
	assert.Equal(t, 0, calls)
}

// a user ID that is not the caller's own is a 403, not a 400 (lfx-self-serve-ops#343)
func TestCreateCompanyHandlerRefusesAnotherUsersID(t *testing.T) {
	api := operations.NewEasyclaAPI(nil)
	service := &fakeCreateCompanyService{err: fmt.Errorf("%w: u-victim", ErrCompanyUserMismatch)}
	Configure(api, service, nil, "")
	email := strfmt.Email("jane@example.invalid")
	name, site := testCompanyName, testCompanySite
	responder := api.CompanyCreateCompanyHandler.Handle(v2CompanyOps.CreateCompanyParams{
		HTTPRequest: httptest.NewRequest(http.MethodPost, "/v4/user/u-victim/company", nil),
		UserID:      "u-victim",
		Input:       &models.CompanyInput{CompanyName: &name, CompanyWebsite: &site, UserEmail: &email},
	}, &auth.User{UserName: "jane", Email: "jane@example.invalid", ACL: auth.ACL{Allowed: true}})
	recorder := httptest.NewRecorder()
	responder.WriteResponse(recorder, runtime.JSONProducer())
	assert.Equal(t, http.StatusForbidden, recorder.Code)
	assert.Equal(t, 1, service.calls)
}
