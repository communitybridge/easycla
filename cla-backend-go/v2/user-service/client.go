package user_service

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/communitybridge/easycla/cla-backend-go/token"
	"github.com/communitybridge/easycla/cla-backend-go/v2/user-service/client"
	"github.com/communitybridge/easycla/cla-backend-go/v2/user-service/client/bulk"
	"github.com/communitybridge/easycla/cla-backend-go/v2/user-service/client/user"
	"github.com/communitybridge/easycla/cla-backend-go/v2/user-service/models"
	runtimeClient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// errors
var (
	ErrUserNotFound = errors.New("user not found")
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

// GetUserByUsername returns user by lfUsername
func (usc *Client) GetUserByUsername(lfUsername string) (*models.User, error) {
	users, err := usc.GetUsersByUsernames([]string{lfUsername})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, ErrUserNotFound
	}
	return users[0], nil
}

// SearchUsers returns a single user based on firstName, lastName and email parameters
func (usc *Client) SearchUsers(firstName string, lastName string, email string) (*models.User, error) {
	params := &user.SearchUsersParams{
		Firstname: &firstName,
		Lastname:  &lastName,
		Email:     &email,
	}
	tok, err := token.GetToken()
	if err != nil {
		return nil, err
	}
	clientAuth := runtimeClient.BearerToken(tok)
	result, err := usc.cl.User.SearchUsers(params, clientAuth)
	if err != nil {
		return nil, err
	}
	users := result.Payload.Data

	if len(users) == 0 {
		return nil, ErrUserNotFound
	}

	return users[0], nil
}

// SearchUserByEmail search user by email
func (usc *Client) SearchUserByEmail(email string) (*models.User, error) {
	params := &user.SearchUsersParams{
		Email:   &email,
		Context: context.Background(),
	}
	tok, err := token.GetToken()
	if err != nil {
		return nil, err
	}
	clientAuth := runtimeClient.BearerToken(tok)
	result, err := usc.cl.User.SearchUsers(params, clientAuth)
	if err != nil {
		return nil, err
	}
	users := result.Payload.Data

	if len(users) == 0 {
		return nil, ErrUserNotFound
	}
	return users[0], nil
}

// ConvertToContact converts user to contact from lead
func (usc *Client) ConvertToContact(userSFID string) error {
	params := &user.ConvertToContactParams{
		SalesforceID: userSFID,
		Context:      context.Background(),
	}
	tok, err := token.GetToken()
	if err != nil {
		return err
	}
	clientAuth := runtimeClient.BearerToken(tok)
	_, _, err = usc.cl.User.ConvertToContact(params, clientAuth) //nolint
	if err != nil {
		return err
	}
	return nil
}

// GetUser returns user from user-service
func (usc *Client) GetUser(userSFID string) (*models.User, error) {
	params := &user.GetUserParams{
		SalesforceID: userSFID,
		Context:      context.Background(),
	}
	tok, err := token.GetToken()
	if err != nil {
		return nil, err
	}
	clientAuth := runtimeClient.BearerToken(tok)
	result, err := usc.cl.User.GetUser(params, clientAuth)
	if err != nil {
		return nil, err
	}
	return result.Payload, nil
}
