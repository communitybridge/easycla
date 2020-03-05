package metrics

import (
	"fmt"
	"time"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
)

// Repository provides methods for calculation,storage and retrieval of metrics
type Repository interface {
	CalculateAndSaveMetrics() error
	GetClaManagerDistribution() (*ClaManagersDistribution, error)
	GetTotalCountMetrics() (*TotalCountMetrics, error)
	GetCompanyMetrics() ([]*CompanyMetric, error)
	GetProjectMetrics() ([]*ProjectMetric, error)
	GetCompanyMetric(companyID string) (*CompanyMetric, error)
	GetProjectMetric(projectID string) (*ProjectMetric, error)
}

type repo struct {
	metricTableName string
	dynamoDBClient  *dynamodb.DynamoDB
	stage           string
}

// NewRepository creates new metrics repository
func NewRepository(awsSession *session.Session, stage string) Repository {
	return &repo{
		dynamoDBClient:  dynamodb.New(awsSession),
		metricTableName: fmt.Sprintf("cla-%s-metrics", stage),
		stage:           stage,
	}
}

// SignatureType constants
const (
	IclaSignature = iota
	CclaSignature
	EmployeeSignature
	InvalidSignature
)

// MetricType constants and ID constants
const (
	MetricTypeTotalCount             = "total_count"
	MetricTypeCompany                = "company"
	MetricTypeProject                = "project"
	MetricTypeClaManagerDistribution = "cla_manager_distribution"

	IDTotalCount             = "total_count"
	IDClaManagerDistribution = "cla_manager_distribution"
)

// ItemSignature represent item of signature table
type ItemSignature struct {
	SignatureReferenceID   string   `json:"signature_reference_id"`
	SignatureReferenceName string   `json:"signature_reference_name"`
	SignatureACL           []string `json:"signature_acl"`
	SignatureUserCompanyID string   `json:"signature_user_ccla_company_id"`
	SignatureType          string   `json:"signature_type"`
	SignatureReferenceType string   `json:"signature_reference_type"`
	SignatureProjectID     string   `json:"signature_project_id"`
}

// ItemRepository represent item of repositories table
type ItemRepository struct {
	RepositoryProjectID string `json:"repository_project_id"`
}

// Metrics contain all metrics related to easycla
type Metrics struct {
	TotalCountMetrics       *TotalCountMetrics       `json:"total_metrics"`
	CompanyMetrics          *CompanyMetrics          `json:"company_metrics"`
	ProjectMetrics          *ProjectMetrics          `json:"project_metrics"`
	ClaManagersDistribution *ClaManagersDistribution `json:"cla_managers_distribution"`
	CalculatedAt            string                   `json:"calculated_at"`
}

// TotalCountMetrics contains all metrics related to total count
type TotalCountMetrics struct {
	CorporateContributorsCount  int `json:"corporate_contributors_count"`
	IndividualContributorsCount int `json:"individual_contributors_count"`
	ClaManagersCount            int `json:"cla_managers_count"`
	ContributorsCount           int `json:"contributors_count"`

	ProjectsCount     int64 `json:"projects_count"`
	RepositoriesCount int64 `json:"repositories_count"`
	CompaniesCount    int64 `json:"companies_count"`

	CreatedAt              string `json:"created_at"`
	corporateContributors  map[string]interface{}
	individualContributors map[string]interface{}
	claManagers            map[string]interface{}
	contributors           map[string]interface{}
}

// CompanyMetric contains all metrics related with particular company
type CompanyMetric struct {
	ID                         string `json:"id"`
	CompanyName                string `json:"company_name"`
	ProjectCount               int    `json:"project_count"`
	CorporateContributorsCount int    `json:"corporate_contributors_count"`
	ClaManagersCount           int    `json:"cla_managers_count"`
	CreatedAt                  string `json:"created_at"`
	corporateContributors      map[string]interface{}
	claManagers                map[string]interface{}
}

// ProjectMetric contains all metrics related with particular project
type ProjectMetric struct {
	ID                          string `json:"id"`
	CompaniesCount              int    `json:"companies_count"`
	ClaManagersCount            int    `json:"cla_managers_count"`
	CorporateContributorsCount  int    `json:"corporate_contributors_count"`
	IndividualContributorsCount int    `json:"individual_contributors_count"`
	TotalContributorsCount      int    `json:"total_contributors_count"`
	RepositoriesCount           int    `json:"repositories_count"`
	CreatedAt                   string `json:"created_at"`
	companies                   map[string]interface{}
	claManagers                 map[string]interface{}
	corporateContributors       map[string]interface{}
	individualContributors      map[string]interface{}
}

// ClaManagersDistribution tells distribution of number of cla mangers associated with company
type ClaManagersDistribution struct {
	OneClaManager        int    `json:"one_cla_manager"`
	TwoClaManager        int    `json:"two_cla_manager"`
	ThreeClaManager      int    `json:"three_cla_manager"`
	FourOrMoreClaManager int    `json:"four_or_more_cla_manager"`
	CreatedAt            string `json:"created_at"`
}

func newMetrics() *Metrics {
	return &Metrics{
		TotalCountMetrics:       newTotalCountMetrics(),
		CompanyMetrics:          newCompanyMetrics(),
		ProjectMetrics:          newProjectMetrics(),
		ClaManagersDistribution: &ClaManagersDistribution{},
	}
}

func newTotalCountMetrics() *TotalCountMetrics {
	return &TotalCountMetrics{
		CorporateContributorsCount:  0,
		corporateContributors:       make(map[string]interface{}),
		IndividualContributorsCount: 0,
		individualContributors:      make(map[string]interface{}),
		ClaManagersCount:            0,
		claManagers:                 make(map[string]interface{}),
		ContributorsCount:           0,
		contributors:                make(map[string]interface{}),
	}
}

func newCompanyMetric() *CompanyMetric {
	return &CompanyMetric{
		ProjectCount:               0,
		CorporateContributorsCount: 0,
		corporateContributors:      make(map[string]interface{}),
		ClaManagersCount:           0,
		claManagers:                make(map[string]interface{}),
	}
}

// CompanyMetrics contain collection of all company metrics
type CompanyMetrics struct {
	CompanyMetrics map[string]*CompanyMetric
}

func newCompanyMetrics() *CompanyMetrics {
	return &CompanyMetrics{
		CompanyMetrics: make(map[string]*CompanyMetric),
	}
}

func newProjectMetric() *ProjectMetric {
	return &ProjectMetric{
		CompaniesCount:              0,
		companies:                   make(map[string]interface{}),
		ClaManagersCount:            0,
		claManagers:                 make(map[string]interface{}),
		CorporateContributorsCount:  0,
		corporateContributors:       make(map[string]interface{}),
		IndividualContributorsCount: 0,
		individualContributors:      make(map[string]interface{}),
		TotalContributorsCount:      0,
		RepositoriesCount:           0,
	}
}

// ProjectMetrics contain collection of all project metrics
type ProjectMetrics struct {
	ProjectMetrics map[string]*ProjectMetric
}

func newProjectMetrics() *ProjectMetrics {
	return &ProjectMetrics{
		ProjectMetrics: make(map[string]*ProjectMetric),
	}
}

func increaseCountIfNotPresent(cacheData map[string]interface{}, count *int, key string) {
	if _, ok := cacheData[key]; !ok {
		cacheData[key] = nil
		*count++
	}
}

func signatureType(sig *ItemSignature) int {
	if sig.SignatureType == "ccla" && sig.SignatureReferenceType == "company" {
		return CclaSignature
	}
	if sig.SignatureType == "cla" {
		if sig.SignatureUserCompanyID != "" {
			return EmployeeSignature
		}
		return IclaSignature
	}
	return InvalidSignature
}

func (m *Metrics) processSignature(sig *ItemSignature) {
	sigType := signatureType(sig)
	if sigType == InvalidSignature {
		log.Printf("Warn: invalid signature: %v\n", sig)
		return
	}
	m.CompanyMetrics.processSignature(sig, sigType)
	m.TotalCountMetrics.processSignature(sig, sigType)
	m.ProjectMetrics.processSignature(sig, sigType)
}

func (cm *TotalCountMetrics) processSignature(sig *ItemSignature, sigType int) {
	switch sigType {
	case CclaSignature:
		for _, acl := range sig.SignatureACL {
			increaseCountIfNotPresent(cm.claManagers, &cm.ClaManagersCount, acl)
		}
	case EmployeeSignature:
		userID := sig.SignatureReferenceID
		increaseCountIfNotPresent(cm.corporateContributors, &cm.CorporateContributorsCount, userID)
		increaseCountIfNotPresent(cm.contributors, &cm.ContributorsCount, userID)
	case IclaSignature:
		userID := sig.SignatureReferenceID
		increaseCountIfNotPresent(cm.individualContributors, &cm.IndividualContributorsCount, userID)
		increaseCountIfNotPresent(cm.contributors, &cm.ContributorsCount, userID)
	}
}

func (cm *CompanyMetrics) processSignature(sig *ItemSignature, sigType int) {
	switch sigType {
	case CclaSignature:
		companyName := sig.SignatureReferenceName
		companyID := sig.SignatureReferenceID
		m, ok := cm.CompanyMetrics[companyID]
		if !ok {
			m = newCompanyMetric()
			cm.CompanyMetrics[companyID] = m
		}
		if m.CompanyName == "" {
			m.CompanyName = companyName
		}
		m.ProjectCount++
		for _, acl := range sig.SignatureACL {
			increaseCountIfNotPresent(m.claManagers, &m.ClaManagersCount, acl)
		}
	case EmployeeSignature:
		companyID := sig.SignatureUserCompanyID
		m, ok := cm.CompanyMetrics[companyID]
		if !ok {
			m = newCompanyMetric()
			cm.CompanyMetrics[companyID] = m
		}
		userID := sig.SignatureReferenceID
		increaseCountIfNotPresent(m.corporateContributors, &m.CorporateContributorsCount, userID)
	}
}

func (pm *ProjectMetrics) processSignature(sig *ItemSignature, sigType int) {
	projectID := sig.SignatureProjectID
	m, ok := pm.ProjectMetrics[projectID]
	if !ok {
		m = newProjectMetric()
		pm.ProjectMetrics[projectID] = m
	}
	switch sigType {
	case CclaSignature:
		companyID := sig.SignatureReferenceID
		increaseCountIfNotPresent(m.companies, &m.CompaniesCount, companyID)
		for _, acl := range sig.SignatureACL {
			increaseCountIfNotPresent(m.claManagers, &m.ClaManagersCount, acl)
		}
	case EmployeeSignature:
		userID := sig.SignatureReferenceID
		increaseCountIfNotPresent(m.corporateContributors, &m.CorporateContributorsCount, userID)
	case IclaSignature:
		userID := sig.SignatureReferenceID
		increaseCountIfNotPresent(m.individualContributors, &m.IndividualContributorsCount, userID)
	}
	m.TotalContributorsCount = m.IndividualContributorsCount + m.CorporateContributorsCount
}

func (pm *ProjectMetrics) processRepositories(repo *ItemRepository) {
	projectID := repo.RepositoryProjectID
	m, ok := pm.ProjectMetrics[projectID]
	if !ok {
		m = newProjectMetric()
		pm.ProjectMetrics[projectID] = m
	}
	m.RepositoriesCount++
}

func calculateClaManagerDistribution(cm *CompanyMetrics) *ClaManagersDistribution {
	var cmd ClaManagersDistribution
	for _, companyMetric := range cm.CompanyMetrics {
		switch companyMetric.ClaManagersCount {
		case 1:
			cmd.OneClaManager++
		case 2:
			cmd.TwoClaManager++
		case 3:
			cmd.ThreeClaManager++
		default:
			if companyMetric.ClaManagersCount >= 4 {
				cmd.FourOrMoreClaManager++
			}
		}
	}
	return &cmd
}

func (repo *repo) processSignaturesTable(metrics *Metrics) error {
	log.Println("processing signatures table")
	filter := expression.Name("signature_signed").Equal(expression.Value(true))
	projection := expression.NamesList(
		expression.Name("signature_reference_id"),
		expression.Name("signature_reference_name"), // Added to support simplified UX queries
		expression.Name("signature_acl"),
		expression.Name("signature_user_ccla_company_id"), // reference to the company
		expression.Name("signature_type"),                 // ccla or cla
		expression.Name("signature_reference_type"),       // user or company
		expression.Name("signature_project_id"),           // project id
	)
	expr, err := expression.NewBuilder().WithFilter(filter).WithProjection(projection).Build()
	if err != nil {
		log.Warnf("error building expression for metric scan, error: %v", err)
		return err
	}

	signatureTableName := fmt.Sprintf("cla-%s-signatures", repo.stage)
	// Assemble the query input parameters
	scanInput := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(signatureTableName),
	}

	for {
		results, err := repo.dynamoDBClient.Scan(scanInput)
		if err != nil {
			log.Warnf("error retrieving signatures, error: %v", err)
			return err
		}

		var sigs []*ItemSignature

		err = dynamodbattribute.UnmarshalListOfMaps(results.Items, &sigs)
		if err != nil {
			log.Warnf("error unmarshalling signatures from database. error: %v", err)
			return err
		}

		for _, sig := range sigs {
			metrics.processSignature(sig)
		}

		if len(results.LastEvaluatedKey) != 0 {
			scanInput.ExclusiveStartKey = results.LastEvaluatedKey
		} else {
			break
		}
	}
	return nil
}

func (repo *repo) processRepositoriesTable(metrics *Metrics) error {
	log.Println("processing repositories table")
	projection := expression.NamesList(
		expression.Name("repository_project_id"),
	)
	expr, err := expression.NewBuilder().WithProjection(projection).Build()
	if err != nil {
		log.Warnf("error building expression for metric scan, error: %v", err)
		return err
	}

	repositoriesTableName := fmt.Sprintf("cla-%s-repositories", repo.stage)
	// Assemble the query input parameters
	scanInput := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(repositoriesTableName),
	}

	for {
		results, err := repo.dynamoDBClient.Scan(scanInput)
		if err != nil {
			log.Warnf("error retrieving repositories, error: %v", err)
			return err
		}

		var repos []*ItemRepository

		err = dynamodbattribute.UnmarshalListOfMaps(results.Items, &repos)
		if err != nil {
			log.Warnf("error unmarshalling repositories from database. error: %v", err)
			return err
		}

		for _, r := range repos {
			metrics.TotalCountMetrics.RepositoriesCount++
			metrics.ProjectMetrics.processRepositories(r)
		}

		if len(results.LastEvaluatedKey) != 0 {
			scanInput.ExclusiveStartKey = results.LastEvaluatedKey
		} else {
			break
		}
	}
	return nil
}

func (repo *repo) processProjectsTable(metrics *Metrics) error {
	projectTableName := fmt.Sprintf("cla-%s-projects", repo.stage)
	projectCount, err := repo.getItemCount(projectTableName)
	if err != nil {
		return err
	}
	metrics.TotalCountMetrics.ProjectsCount = projectCount
	return nil
}

func (repo *repo) processCompaniesTable(metrics *Metrics) error {
	companiesTableName := fmt.Sprintf("cla-%s-companies", repo.stage)
	companiesCount, err := repo.getItemCount(companiesTableName)
	if err != nil {
		return err
	}
	metrics.TotalCountMetrics.CompaniesCount = companiesCount
	return nil
}

func (repo *repo) calculateMetrics() (*Metrics, error) {
	metrics := newMetrics()
	t := time.Now()
	err := repo.processSignaturesTable(metrics)
	if err != nil {
		return nil, err
	}
	err = repo.processRepositoriesTable(metrics)
	if err != nil {
		return nil, err
	}
	err = repo.processProjectsTable(metrics)
	if err != nil {
		return nil, err
	}
	err = repo.processCompaniesTable(metrics)
	if err != nil {
		return nil, err
	}
	metrics.ClaManagersDistribution = calculateClaManagerDistribution(metrics.CompanyMetrics)
	_, metrics.CalculatedAt = utils.CurrentTime()
	log.Println("GetMetrics took time", time.Since(t).String())
	return metrics, nil
}

func (repo *repo) saveMetrics(metrics *Metrics) error {
	t := time.Now()
	err := repo.saveTotalMetris(metrics.TotalCountMetrics)
	if err != nil {
		return err
	}
	err = repo.saveCompaniesMetrics(metrics.CompanyMetrics)
	if err != nil {
		return err
	}
	err = repo.saveProjectMetrics(metrics.ProjectMetrics)
	if err != nil {
		return err
	}
	err = repo.saveClaManagerDistribution(metrics.ClaManagersDistribution)
	if err != nil {
		return err
	}
	log.Printf("save metrics took :%s \n", time.Since(t).String())
	return nil
}

func addIDTypeTime(item map[string]*dynamodb.AttributeValue, id string, metricType string) {
	_, ctime := utils.CurrentTime()
	utils.AddStringAttribute(item, "id", id)
	utils.AddStringAttribute(item, "metric_type", metricType)
	utils.AddStringAttribute(item, "created_at", ctime)
}

func (repo *repo) saveTotalMetris(tm *TotalCountMetrics) error {
	log.Println("saving total count metrics")
	av, err := dynamodbattribute.MarshalMap(tm)
	if err != nil {
		return err
	}
	addIDTypeTime(av, IDTotalCount, MetricTypeTotalCount)
	_, err = repo.dynamoDBClient.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(repo.metricTableName),
	})
	if err != nil {
		log.Error("cannot put total_metrics in dynamodb", err)
		return err
	}
	return nil
}

func (repo *repo) saveCompaniesMetrics(companyMetrics *CompanyMetrics) error {
	t := time.Now()
	log.Println("saving company_metrics")
	for id, cm := range companyMetrics.CompanyMetrics {
		av, err := dynamodbattribute.MarshalMap(cm)
		if err != nil {
			return err
		}
		addIDTypeTime(av, id, MetricTypeCompany)
		_, err = repo.dynamoDBClient.PutItem(&dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(repo.metricTableName),
		})
		if err != nil {
			log.Printf("cannot put company_metric in dynamodb, metric = %v, error = %s\n", cm, err.Error())
			return err
		}
	}
	log.Printf("saving company_metrics took :%s \n", time.Since(t).String())
	return nil
}

func (repo *repo) saveClaManagerDistribution(cmd *ClaManagersDistribution) error {
	log.Println("saving cla_managers_distribution")
	av, err := dynamodbattribute.MarshalMap(cmd)
	if err != nil {
		return err
	}
	addIDTypeTime(av, IDClaManagerDistribution, MetricTypeClaManagerDistribution)
	_, err = repo.dynamoDBClient.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(repo.metricTableName),
	})
	if err != nil {
		log.Printf("cannot put cla_managers_distribution in dynamodb, metric = %v, error = %s\n", cmd, err.Error())
		return err
	}
	return nil
}

func (repo *repo) saveProjectMetrics(projectMetrics *ProjectMetrics) error {
	t := time.Now()
	log.Println("saving project_metrics")
	for id, cm := range projectMetrics.ProjectMetrics {
		av, err := dynamodbattribute.MarshalMap(cm)
		if err != nil {
			return err
		}
		addIDTypeTime(av, id, MetricTypeProject)
		_, err = repo.dynamoDBClient.PutItem(&dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(repo.metricTableName),
		})
		if err != nil {
			log.Printf("cannot put project_metric in dynamodb, metric = %v, error = %s\n", cm, err.Error())
			return err
		}
	}
	log.Printf("saving project_metrics took :%s \n", time.Since(t).String())
	return nil
}

func (repo *repo) CalculateAndSaveMetrics() error {
	m, err := repo.calculateMetrics()
	if err != nil {
		return err
	}
	err = repo.saveMetrics(m)
	if err != nil {
		return err
	}
	return nil
}

func (repo *repo) GetClaManagerDistribution() (*ClaManagersDistribution, error) {
	var out ClaManagersDistribution
	err := repo.getMetricByID(IDClaManagerDistribution, MetricTypeClaManagerDistribution, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (repo *repo) GetTotalCountMetrics() (*TotalCountMetrics, error) {
	var out TotalCountMetrics
	err := repo.getMetricByID(IDTotalCount, MetricTypeTotalCount, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (repo *repo) GetCompanyMetrics() ([]*CompanyMetric, error) {

	keyCondition := expression.Key("metric_type").Equal(expression.Value(MetricTypeCompany))

	expr, err := expression.NewBuilder().WithKeyCondition(keyCondition).Build()
	if err != nil {
		log.Warnf("error building expression for company metric scan, error: %v", err)
		return nil, err
	}

	// Assemble the query input parameters
	queryInput := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		TableName:                 aws.String(repo.metricTableName),
	}

	var companyMetrics []*CompanyMetric
	for {
		results, errQuery := repo.dynamoDBClient.Query(queryInput)
		if errQuery != nil {
			log.Warnf("error retrieving company metrics, error: %v", errQuery)
			return nil, errQuery
		}

		var companyMetricsTmp []*CompanyMetric

		err := dynamodbattribute.UnmarshalListOfMaps(results.Items, &companyMetricsTmp)
		if err != nil {
			log.Warnf("error unmarshalling company metrics from database. error: %v", err)
			return nil, err
		}
		companyMetrics = append(companyMetrics, companyMetricsTmp...)

		if len(results.LastEvaluatedKey) != 0 {
			queryInput.ExclusiveStartKey = results.LastEvaluatedKey
		} else {
			break
		}
	}
	return companyMetrics, nil
}

func (repo *repo) GetProjectMetrics() ([]*ProjectMetric, error) {

	keyCondition := expression.Key("metric_type").Equal(expression.Value(MetricTypeProject))

	expr, err := expression.NewBuilder().WithKeyCondition(keyCondition).Build()
	if err != nil {
		log.Warnf("error building expression for company metric scan, error: %v", err)
		return nil, err
	}

	// Assemble the query input parameters
	queryInput := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		TableName:                 aws.String(repo.metricTableName),
	}

	var projectMetrics []*ProjectMetric
	for {
		results, errQuery := repo.dynamoDBClient.Query(queryInput)
		if errQuery != nil {
			log.Warnf("error retrieving project metrics, error: %v", errQuery)
			return nil, errQuery
		}

		var projectMetricsTmp []*ProjectMetric

		err := dynamodbattribute.UnmarshalListOfMaps(results.Items, &projectMetricsTmp)
		if err != nil {
			log.Warnf("error unmarshalling project metrics from database. error: %v", err)
			return nil, err
		}
		projectMetrics = append(projectMetrics, projectMetricsTmp...)

		if len(results.LastEvaluatedKey) != 0 {
			queryInput.ExclusiveStartKey = results.LastEvaluatedKey
		} else {
			break
		}
	}
	return projectMetrics, nil
}
func (repo *repo) GetCompanyMetric(companyID string) (*CompanyMetric, error) {
	var out CompanyMetric
	err := repo.getMetricByID(companyID, MetricTypeCompany, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (repo *repo) GetProjectMetric(projectID string) (*ProjectMetric, error) {
	var out ProjectMetric
	err := repo.getMetricByID(projectID, MetricTypeProject, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (repo *repo) getMetricByID(id string, metricType string, out interface{}) error {
	result, err := repo.dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(repo.metricTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
			"metric_type": {
				S: aws.String(metricType),
			},
		},
	})
	if err != nil {
		return err
	}
	if len(result.Item) == 0 {
		return fmt.Errorf("metric with id:%s metric_type:%s not found", id, metricType)
	}
	err = dynamodbattribute.UnmarshalMap(result.Item, out)
	if err != nil {
		return err
	}
	return nil
}

func (repo *repo) getItemCount(tableName string) (int64, error) {
	// How many total records do we have - may not be up-to-date as this value is updated only periodically
	describeTableInput := &dynamodb.DescribeTableInput{
		TableName: &tableName,
	}
	describeTableResult, err := repo.dynamoDBClient.DescribeTable(describeTableInput)
	if err != nil {
		log.Warnf("error retrieving total record count of table %s, error: %v", tableName, err)
		return 0, err
	}
	return *describeTableResult.Table.ItemCount, nil
}
