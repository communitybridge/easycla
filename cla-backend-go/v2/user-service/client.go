package user_service

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/communitybridge/easycla/cla-backend-go/token"
	"github.com/communitybridge/easycla/cla-backend-go/v2/user-service/client"
	"github.com/communitybridge/easycla/cla-backend-go/v2/user-service/client/bulk"
	"github.com/communitybridge/easycla/cla-backend-go/v2/user-service/models"
	runtimeClient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

type Client struct {
	cl        *client.UserService
	dummyXACL string
}

var (
	userServiceClient *Client
)

func InitClient(ApiGwURL string) {
	ApiGwURL = strings.ReplaceAll(ApiGwURL, "https://", "")
	userServiceClient = &Client{
		cl: client.NewHTTPClientWithConfig(strfmt.Default, &client.TransportConfig{
			Host:     ApiGwURL,
			BasePath: "user-service/v1",
			Schemes:  []string{"https"},
		}),
		dummyXACL: "ewogICAgImFsbG93ZWQiOnRydWUsCiAgICAiaXNBZG1pbiI6IHRydWUsICAgIAogICAgInJlc291cmNlIjoidmlld19wcm9qZWN0IiwKICAgICJzY29wZXMiOlsgCiAgICBdIAp9",
	}
}

func GetClient() *Client {
	return userServiceClient
}

func (usc *Client) GetUsersByUsernames(lfUsernames []string) ([]*models.User, error) {
	params := bulk.NewSearchBulkParams()
	params.SearchBulk = &models.SearchBulk{
		List: lfUsernames,
		Type: aws.String("username"),
	}
	tok, err := token.GetToken()
	if err != nil {
		return nil, err
	}
	clientAuth := runtimeClient.BearerToken(tok)
	params.XACL = usc.dummyXACL
	result, err := usc.cl.Bulk.SearchBulk(params, clientAuth)
	if err != nil {
		return nil, err
	}
	return result.Payload.Data, nil
}
