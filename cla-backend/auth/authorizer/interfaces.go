// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later
package authorizer

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// Interfaces defines an interface to interact with the Interfaces layer
type Interfaces interface {
	Handler(context.Context, events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error)
}

type interfacesContainer struct {
	usecases Usecases
}

// NewInterfaces creates an InterfacesInteractor
func NewInterfaces(usecases Usecases) Interfaces {
	return &interfacesContainer{
		usecases: usecases,
	}
}

// Handler authorizes a request
func (ic *interfacesContainer) Handler(
	ctx context.Context,
	event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := event.AuthorizationToken
	parsedToken, err := ic.usecases.ValidateToken(token)
	if err != nil {
		fields := map[string]interface{}{}
		log.Print(fields, err.Error())
		return events.APIGatewayCustomAuthorizerResponse{}, err
	}

	return generateAuthResponse(event.MethodArn, &parsedToken), nil
}

func generateAuthResponse(methodArn string, token *TokenInfo) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: token.Subject}

	authResponse.PolicyDocument = generatePolicy(methodArn)
	authResponse.Context = map[string]interface{}{
		"subject":       token.Subject,
		"email":         token.Email,
		"emailVerified": token.EmailVerified,
	}
	return authResponse
}

func generatePolicy(methodARN string) events.APIGatewayCustomAuthorizerPolicy {
	resource := getWildcardFunctionARN(methodARN)
	return events.APIGatewayCustomAuthorizerPolicy{
		Version: "2012-10-17",
		Statement: []events.IAMPolicyStatement{
			{
				Action:   []string{"execute-api:Invoke"},
				Effect:   "Allow",
				Resource: []string{resource},
			},
		},
	}
}

func getWildcardFunctionARN(methodARN string) string {
	paths := strings.Split(methodARN, "/")
	return paths[0] + "/" + paths[1] + "/*"
}
