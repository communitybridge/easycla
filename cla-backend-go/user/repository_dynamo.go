package user

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type RepositoryDynamo struct {
	Stage              string
	DynamoDBClient     *dynamodb.DynamoDB
	senderEmailAddress string
}

type User struct {
	UserID             string   `json:"user_id"`
	LFEmail            string   `json:"lf_email"`
	LFUsername         string   `json:"lf_username"`
	DateCreated        string   `json:"date_created"`
	DateModified       string   `json:"date_modified"`
	UserName           string   `json:"user_name"`
	Version            string   `json:"version"`
	UserEmails         []string `json:"user_emails"`
	UserGithubID       string   `json:"user_github_id"`
	UserCompanyID      string   `json:"user_company_id"`
	UserGithubUsername string   `json:"user_github_username"`
}

func NewDynamoRepository(awsSession *session.Session, stage, senderEmailAddress string) RepositoryDynamo {
	return RepositoryDynamo{
		Stage:              stage,
		DynamoDBClient:     dynamodb.New(awsSession),
		senderEmailAddress: senderEmailAddress,
	}
}

func (repo RepositoryDynamo) GetUserAndProfilesByLFID(lfidUsername string) (CLAUser, error) {
	tableName := fmt.Sprintf("cla-%s-users", repo.Stage)

	input := &dynamodb.QueryInput{
		KeyConditions: map[string]*dynamodb.Condition{
			"lf_username": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(lfidUsername),
					},
				},
			},
		},
		TableName: aws.String(tableName),
		IndexName: aws.String("lf-username-index"),
	}
	result, err := repo.DynamoDBClient.Query(input)

	if err != nil {
		fmt.Println("Unable to retrieve data from users")
		return CLAUser{}, err
	}

	users := []User{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		fmt.Println(err.Error())
		return CLAUser{}, err
	}

	if len(users) < 1 {
		fmt.Println(fmt.Sprintf("No user has been found with the given LFID: %s", lfidUsername))
		return CLAUser{}, err
	}

	claUser := CLAUser{
		UserID:     users[0].UserID,
		Name:       users[0].UserName,
		LFEmail:    users[0].LFEmail,
		LFUsername: users[0].LFUsername,
	}

	return claUser, nil
}

func (repo RepositoryDynamo) GetUserProjectIDs(userID string) ([]string, error) {
	tableName := fmt.Sprintf("cla-%s-user-permissions", repo.Stage)
	result, err := repo.DynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(userID),
			},
		},
	})

	if err != nil {
		fmt.Println(err.Error())
		return []string{}, err
	}

	projects, ok := result.Item["projects"]
	if !ok {
		projects = &dynamodb.AttributeValue{}
	}

	//take off pointer from []*string
	return aws.StringValueSlice(projects.SS), nil
}

func (repo RepositoryDynamo) GetClaManagerCorporateClaIDs(userID string) ([]string, error) {
	return []string{}, nil
}

func (repo RepositoryDynamo) GetUserCompanyIDs(userID string) ([]string, error) {
	return []string{}, nil
}

func (repo RepositoryDynamo) GetUser(userID string) (User, error) {
	tableName := fmt.Sprintf("cla-%s-users", repo.Stage)
	userAV, err := repo.DynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"user_id": {
				S: aws.String(userID),
			},
		},
	})

	user := User{}
	err = dynamodbattribute.UnmarshalMap(userAV.Item, &user)
	if err != nil {
		return User{}, err
	}

	return user, err
}
