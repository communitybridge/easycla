// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package company

import (
	"context"
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

type fakeCreateCompanyService struct {
	Service
	calls int
	seen  string
}

func (f *fakeCreateCompanyService) CreateCompany(_ context.Context, params *v2CompanyOps.CreateCompanyParams) (*models.CompanyOutput, error) {
	f.calls++
	f.seen = params.Input.UserEmail.String()
	return &models.CompanyOutput{CompanyID: "0014100000Created", CompanyName: *params.Input.CompanyName}, nil
}

// the company can only be created for the signed-in user (lfx-self-serve-ops#343)
func TestCreateCompanyHandlerBindsTheCaller(t *testing.T) {
	call := func(t *testing.T, authUser *auth.User, bodyEmail string) (int, int) {
		api := operations.NewEasyclaAPI(nil)
		service := &fakeCreateCompanyService{}
		Configure(api, service, nil, "")
		require.NotNil(t, api.CompanyCreateCompanyHandler)
		email := strfmt.Email(bodyEmail)
		name, site := "Acme Inc", "https://acme.example.invalid"
		responder := api.CompanyCreateCompanyHandler.Handle(v2CompanyOps.CreateCompanyParams{
			HTTPRequest: httptest.NewRequest(http.MethodPost, "/v4/user/user-1/company", nil),
			UserID:      "user-1",
			Input:       &models.CompanyInput{CompanyName: &name, CompanyWebsite: &site, UserEmail: &email},
		}, authUser)
		recorder := httptest.NewRecorder()
		responder.WriteResponse(recorder, runtime.JSONProducer())
		return recorder.Code, service.calls
	}
	caller := &auth.User{UserName: "jane", Email: "jane@example.invalid", ACL: auth.ACL{Allowed: true}}
	status, calls := call(t, caller, "mallory@example.invalid")
	assert.Equal(t, http.StatusForbidden, status)
	assert.Equal(t, 0, calls, "no organization or role may be created for another person's email")
	status, calls = call(t, caller, "Jane@Example.invalid")
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, 1, calls)
	status, calls = call(t, &auth.User{UserName: "ghost", ACL: auth.ACL{Allowed: true}}, "ghost@example.invalid")
	assert.Equal(t, http.StatusForbidden, status, "a principal without an email cannot be bound")
	assert.Equal(t, 0, calls)
}
