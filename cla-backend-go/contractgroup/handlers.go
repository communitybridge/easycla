package contractgroup

import (
	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/gen/models"
	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/gen/restapi/operations"
	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/gen/restapi/operations/contract_group"

	"github.com/go-openapi/runtime/middleware"
)

func Configure(api *operations.ClaAPI, service Service) {
	//Create a Contract Group
	api.ContractGroupCreateContractGroupHandler = contract_group.CreateContractGroupHandlerFunc(func(contractGroup contract_group.CreateContractGroupParams) middleware.Responder {

		contract, err := service.CreateContractGroup(contractGroup.HTTPRequest.Context(), contractGroup.ProjectSfdcID, contractGroup.Body)
		if err != nil {
			return contract_group.NewCreateContractGroupBadRequest().WithPayload(errorResponse(err))
		}

		return contract_group.NewCreateContractGroupOK().WithPayload(contract)
	})

	//Get Contract Group by Project ID
	api.ContractGroupGetContractGroupsHandler = contract_group.GetContractGroupsHandlerFunc(func(params contract_group.GetContractGroupsParams) middleware.Responder {

		contractGroups, err := service.GetContractGroups(params.HTTPRequest.Context(), params.ProjectSfdcID)
		if err != nil {
			return contract_group.NewGetContractGroupsBadRequest().WithPayload(errorResponse(err))
		}

		return contract_group.NewGetContractGroupsOK().WithPayload(contractGroups)
	})

	//Create Contract Template
	api.ContractGroupCreateContractTemplateHandler = contract_group.CreateContractTemplateHandlerFunc(func(params contract_group.CreateContractTemplateParams) middleware.Responder {

		contractTemplate, err := service.CreateContractTemplate(params.HTTPRequest.Context(), params.Body, params.ContractGroupID)
		if err != nil {
			return contract_group.NewCreateContractGroupBadRequest().WithPayload(errorResponse(err))
		}

		return contract_group.NewCreateContractGroupOK().WithPayload(contractTemplate)
	})

	//Adding Github Org name
	api.ContractGroupAddGitHubOrgHandler = contract_group.AddGitHubOrgHandlerFunc(func(params contract_group.AddGitHubOrgParams) middleware.Responder {

		githubOrg, err := service.CreateGitHubOrganization(params.HTTPRequest.Context(), params.ContractGroupID, params.Body)
		if err != nil {
			return contract_group.NewAddGitHubOrgBadRequest().WithPayload(errorResponse(err))
		}

		return contract_group.NewAddGitHubOrgOK().WithPayload(githubOrg)
	})

	//Adding Gerrit Instance
	api.ContractGroupAddGerritInstanceHandler = contract_group.AddGerritInstanceHandlerFunc(func(gerritInstance contract_group.AddGerritInstanceParams) middleware.Responder {

		gerritInstanceResponse, err := service.CreateGerritInstance(gerritInstance.HTTPRequest.Context(), gerritInstance.ProjectSfdcID, gerritInstance.ContractGroupID, userID, gerritInstance.Body)
		if err != nil {
			return contract_group.NewAddGerritInstanceBadRequest().WithPayload(errorResponse(err))
		}

		return contract_group.NewAddGerritInstanceOK().WithPayload(gerritInstanceResponse)
	})

	//Deleting Gerrit Instance
	api.ContractGroupDeleteGerritInstanceHandler = contract_group.DeleteGerritInstanceHandlerFunc(func(gerritInstance contract_group.DeleteGerritInstanceParams) middleware.Responder {

		err := service.DeleteGerritInstance(gerritInstance.HTTPRequest.Context(), gerritInstance.ProjectSfdcID, gerritInstance.ContractGroupID, gerritInstance.GerritInstanceID)
		if err != nil {
			return contract_group.NewAddGerritInstanceBadRequest().WithPayload(errorResponse(err))
		}

		return contract_group.NewAddGerritInstanceOK()
	})

	//Getting Contract Group Signatures
	api.ContractGroupGetContractGroupSignaturesHandler = contract_group.GetContractGroupSignaturesHandlerFunc(func(params contract_group.GetContractGroupSignaturesParams) middleware.Responder {

		contractGroupSignatures, err := service.GetContractGroupSignatures(params.HTTPRequest.Context(), params.ProjectSfdcID, params.ContractGroupID)
		if err != nil {
			return contract_group.NewGetContractGroupSignaturesBadRequest().WithPayload(errorResponse(err))
		}

		return contract_group.NewGetContractGroupSignaturesOK().WithPayload(&contractGroupSignatures)
	})
}

type codedResponse interface {
	Code() string
}

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
