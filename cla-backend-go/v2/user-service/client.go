package user_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/communitybridge/easycla/cla-backend-go/v2/user-service/client/staff"

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
	cl       *client.UserService
	apiKey   string
	apiGwURL string
}

var (
	userServiceClient *Client
)

// InitClient initializes the user_service client
func InitClient(APIGwURL string, apiKey string) {
	APIGwURL = strings.ReplaceAll(APIGwURL, "https://", "")
	userServiceClient = &Client{
		apiKey:   apiKey,
		apiGwURL: APIGwURL,
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

	query := fmt.Sprintf("email=%s&firstname=%s&lastname=%s", email, firstName, lastName)
	url := fmt.Sprintf("https://%s/user-service/v1/users/search?%s", usc.apiGwURL, query)
	tok, err := token.GetToken()
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("X-API-KEY", usc.apiKey)
	request.Header.Set("Authorization", "Bearer "+tok)
	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	userList, err := getUsers(data)
	if err != nil {
		return nil, err
	}
	for _, userItem := range userList {
		for _, userEmail := range userItem.Emails {
			if *userEmail.EmailAddress == email {
				if userItem.FirstName == firstName && userItem.LastName == lastName {
					return userItem, nil
				}
			}
		}
	}

	return nil, errors.New("user not found")
}

func getUsers(body []byte) ([]*models.User, error) {
	var users = new(models.UserList)
	err := json.Unmarshal(body, &users)
	if err != nil {
		return nil, err
	}
	return users.Data, err
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

// GetStaff returns staff details from user-service
func (usc *Client) GetStaff(userSFID string) (*models.Staff, error) {
	params := &staff.GetStaffParams{
		SalesforceID: userSFID,
		Context:      context.Background(),
	}
	tok, err := token.GetToken()
	if err != nil {
		return nil, err
	}
	clientAuth := runtimeClient.BearerToken(tok)
	result, err := usc.cl.Staff.GetStaff(params, clientAuth)
	if err != nil {
		return nil, err
	}
	return result.Payload, nil
}
