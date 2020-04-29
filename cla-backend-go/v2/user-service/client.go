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

// Client is client for user_service
type Client struct {
	cl *client.UserService
}

var (
	userServiceClient *Client
)

// InitClient initializes the user_service client
func InitClient(APIGwURL string) {
	APIGwURL = strings.ReplaceAll(APIGwURL, "https://", "")
	userServiceClient = &Client{
		cl: client.NewHTTPClientWithConfig(strfmt.Default, &client.TransportConfig{
			Host:     APIGwURL,
			BasePath: "user-service/v1",
			Schemes:  []string{"https"},
		}),
	}
}

// GetClient return user_service client
func GetClient() *Client {
	return userServiceClient
}

// GetUsersByUsernames search users by lf username
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
	result, err := usc.cl.Bulk.SearchBulk(params, clientAuth)
	if err != nil {
		return nil, err
	}
	return result.Payload.Data, nil
}
