// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package contractgroup

import (
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/contract_group"
	"github.com/communitybridge/easycla/cla-backend-go/user"

	"github.com/go-openapi/runtime/middleware"
)

// Configure sets up the middleware handlers for the contract group service
func Configure(api *operations.ClaAPI, service Service) {
	//Create a Contract Group
	api.ContractGroupCreateContractGroupHandler = contract_group.CreateContractGroupHandlerFunc(func(contractGroup contract_group.CreateContractGroupParams, claUser *user.CLAUser) middleware.Responder {

		contract, err := service.CreateContractGroup(contractGroup.HTTPRequest.Context(), contractGroup.ProjectSfdcID, contractGroup.Body)
		if err != nil {
			return contract_group.NewCreateContractGroupBadRequest().WithPayload(errorResponse(err))
		}

		return contract_group.NewCreateContractGroupOK().WithPayload(contract)
	})

	//Get Contract Group by Project ID
	api.ContractGroupGetContractGroupsHandler = contract_group.GetContractGroupsHandlerFunc(func(params contract_group.GetContractGroupsParams, claUser *user.CLAUser) middleware.Responder {

		contractGroups, err := service.GetContractGroups(params.HTTPRequest.Context(), params.ProjectSfdcID)
		if err != nil {
			return contract_group.NewGetContractGroupsBadRequest().WithPayload(errorResponse(err))
		}

		return contract_group.NewGetContractGroupsOK().WithPayload(contractGroups)
	})

	//Create Contract Template
	api.ContractGroupCreateContractTemplateHandler = contract_group.CreateContractTemplateHandlerFunc(func(params contract_group.CreateContractTemplateParams, claUser *user.CLAUser) middleware.Responder {

		contractTemplate, err := service.CreateContractTemplate(params.HTTPRequest.Context(), params.Body, params.ContractGroupID)
		if err != nil {
			return contract_group.NewCreateContractGroupBadRequest().WithPayload(errorResponse(err))
		}

		return contract_group.NewCreateContractGroupOK().WithPayload(contractTemplate)
	})

	//Adding Github Org name
	api.ContractGroupAddGitHubOrgHandler = contract_group.AddGitHubOrgHandlerFunc(func(params contract_group.AddGitHubOrgParams, claUser *user.CLAUser) middleware.Responder {

		githubOrg, err := service.CreateGitHubOrganization(params.HTTPRequest.Context(), params.ContractGroupID, params.Body)
		if err != nil {
			return contract_group.NewAddGitHubOrgBadRequest().WithPayload(errorResponse(err))
		}

		return contract_group.NewAddGitHubOrgOK().WithPayload(githubOrg)
	})

	//Adding Gerrit Instance
	api.ContractGroupAddGerritInstanceHandler = contract_group.AddGerritInstanceHandlerFunc(func(gerritInstance contract_group.AddGerritInstanceParams, claUser *user.CLAUser) middleware.Responder {

		gerritInstanceResponse, err := service.CreateGerritInstance(gerritInstance.HTTPRequest.Context(), gerritInstance.ProjectSfdcID, gerritInstance.ContractGroupID, userID, gerritInstance.Body)
		if err != nil {
			return contract_group.NewAddGerritInstanceBadRequest().WithPayload(errorResponse(err))
		}

		return contract_group.NewAddGerritInstanceOK().WithPayload(gerritInstanceResponse)
	})

	//Deleting Gerrit Instance
	api.ContractGroupDeleteGerritInstanceHandler = contract_group.DeleteGerritInstanceHandlerFunc(func(gerritInstance contract_group.DeleteGerritInstanceParams, claUser *user.CLAUser) middleware.Responder {

		err := service.DeleteGerritInstance(gerritInstance.HTTPRequest.Context(), gerritInstance.ProjectSfdcID, gerritInstance.ContractGroupID, gerritInstance.GerritInstanceID)
		if err != nil {
			return contract_group.NewAddGerritInstanceBadRequest().WithPayload(errorResponse(err))
		}

		return contract_group.NewAddGerritInstanceOK()
	})

	//Getting Contract Group Signatures
	api.ContractGroupGetContractGroupSignaturesHandler = contract_group.GetContractGroupSignaturesHandlerFunc(func(params contract_group.GetContractGroupSignaturesParams, claUser *user.CLAUser) middleware.Responder {

		contractGroupSignatures, err := service.GetContractGroupSignatures(params.HTTPRequest.Context(), params.ProjectSfdcID, params.ContractGroupID)
		if err != nil {
			return contract_group.NewGetContractGroupSignaturesBadRequest().WithPayload(errorResponse(err))
		}

		return contract_group.NewGetContractGroupSignaturesOK().WithPayload(&contractGroupSignatures)
	})
}

// codedResponse interface
type codedResponse interface {
	Code() string
}

// errorResponse is a helper function to wrap the specified error in a error response model
func errorResponse(err error) *models.ErrorResponse {
	code := ""
	if e, ok := err.(codedResponse); ok {
		code = e.Code()
	}

	e := models.ErrorResponse{
		Code:    code,
		Message: err.Error(),
	}

	return &e
}
