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

type Repository interface {
	CalculateAndSaveMetrics() error
}

type repo struct {
	metricTableName string
	dynamoDBClient  *dynamodb.DynamoDB
	stage           string
}

func NewRepository(awsSession *session.Session, stage string) Repository {
	return &repo{
		dynamoDBClient:  dynamodb.New(awsSession),
		metricTableName: fmt.Sprintf("cla-%s-metrics", stage),
		stage:           stage,
	}
}

const (
	IclaSignature = iota
	CclaSignature
	EmployeeSignature
	InvalidSignature
)

type ItemSignature struct {
	SignatureReferenceID   string   `json:"signature_reference_id"`
	SignatureReferenceName string   `json:"signature_reference_name"`
	SignatureACL           []string `json:"signature_acl"`
	SignatureUserCompanyID string   `json:"signature_user_ccla_company_id"`
	SignatureType          string   `json:"signature_type"`
	SignatureReferenceType string   `json:"signature_reference_type"`
	SignatureProjectID     string   `json:"signature_project_id"`
}

type ItemRepository struct {
	RepositoryProjectID string `json:"repository_project_id"`
}

type Metrics struct {
	TotalMetrics            *TotalMetrics            `json:"total_metrics"`
	CompanyMetrics          *CompanyMetrics          `json:"company_metrics"`
	ProjectMetrics          *ProjectMetrics          `json:"project_metrics"`
	ClaManagersDistribution *ClaManagersDistribution `json:"cla_managers_distribution"`
	CalculatedAt            string                   `json:"calculated_at"`
}

type TotalMetrics struct {
	CorporateContributorsCount  int                    `json:"corporate_contributors_count"`
	CorporateContributors       map[string]interface{} `json:"-"`
	IndividualContributorsCount int                    `json:"individual_contributors_count"`
	IndividualContributors      map[string]interface{} `json:"-"`
	ClaManagersCount            int                    `json:"cla_managers_count"`
	ClaManagers                 map[string]interface{} `json:"-"`
	ContributorsCount           int                    `json:"contributors_count"`
	Contributors                map[string]interface{} `json:"-"`
}

type CompanyMetric struct {
	CompanyName                string                 `json:"company_name"`
	ProjectCount               int                    `json:"project_count"`
	CorporateContributorsCount int                    `json:"corporate_contributors_count"`
	CorporateContributors      map[string]interface{} `json:"-"`
	ClaManagersCount           int                    `json:"cla_managers_count"`
	ClaManagers                map[string]interface{} `json:"-"`
}

type ProjectMetric struct {
	CompaniesCount              int                    `json:"companies_count"`
	Companies                   map[string]interface{} `json:"-"`
	ClaManagersCount            int                    `json:"cla_managers_count"`
	ClaManagers                 map[string]interface{} `json:"-"`
	CorporateContributorsCount  int                    `json:"corporate_contributors_count"`
	CorporateContributors       map[string]interface{} `json:"-"`
	IndividualContributorsCount int                    `json:"individual_contributors_count"`
	IndividualContributors      map[string]interface{} `json:"-"`
	TotalContributorsCount      int                    `json:"total_contributors_count"`
	RepositoriesCount           int                    `json:"repositories_count"`
}

func newMetrics() *Metrics {
	return &Metrics{
		TotalMetrics:            newTotalMetrics(),
		CompanyMetrics:          newCompanyMetrics(),
		ProjectMetrics:          newProjectMetrics(),
		ClaManagersDistribution: &ClaManagersDistribution{},
	}
}

func newTotalMetrics() *TotalMetrics {
	return &TotalMetrics{
		CorporateContributorsCount:  0,
		CorporateContributors:       make(map[string]interface{}),
		IndividualContributorsCount: 0,
		IndividualContributors:      make(map[string]interface{}),
		ClaManagersCount:            0,
		ClaManagers:                 make(map[string]interface{}),
		ContributorsCount:           0,
		Contributors:                make(map[string]interface{}),
	}
}

func newCompanyMetric() *CompanyMetric {
	return &CompanyMetric{
		ProjectCount:               0,
		CorporateContributorsCount: 0,
		CorporateContributors:      make(map[string]interface{}),
		ClaManagersCount:           0,
		ClaManagers:                make(map[string]interface{}),
	}
}

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
		Companies:                   make(map[string]interface{}),
		ClaManagersCount:            0,
		ClaManagers:                 make(map[string]interface{}),
		CorporateContributorsCount:  0,
		CorporateContributors:       make(map[string]interface{}),
		IndividualContributorsCount: 0,
		IndividualContributors:      make(map[string]interface{}),
		TotalContributorsCount:      0,
		RepositoriesCount:           0,
	}
}

type ProjectMetrics struct {
	ProjectMetrics map[string]*ProjectMetric
}

func newProjectMetrics() *ProjectMetrics {
	return &ProjectMetrics{
		ProjectMetrics: make(map[string]*ProjectMetric),
	}
}

type ClaManagersDistribution struct {
	OneClaManager        int `json:"one_cla_manager"`
	TwoClaManager        int `json:"two_cla_manager"`
	ThreeClaManager      int `json:"three_cla_manager"`
	FourOrMoreClaManager int `json:"four_or_more_cla_manager"`
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
		} else {
			return IclaSignature
		}
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
	m.TotalMetrics.processSignature(sig, sigType)
	m.ProjectMetrics.processSignature(sig, sigType)
}

func (cm *TotalMetrics) processSignature(sig *ItemSignature, sigType int) {
	switch sigType {
	case CclaSignature:
		for _, acl := range sig.SignatureACL {
			increaseCountIfNotPresent(cm.ClaManagers, &cm.ClaManagersCount, acl)
		}
	case EmployeeSignature:
		userID := sig.SignatureReferenceID
		increaseCountIfNotPresent(cm.CorporateContributors, &cm.CorporateContributorsCount, userID)
		increaseCountIfNotPresent(cm.Contributors, &cm.ContributorsCount, userID)
	case IclaSignature:
		userID := sig.SignatureReferenceID
		increaseCountIfNotPresent(cm.IndividualContributors, &cm.IndividualContributorsCount, userID)
		increaseCountIfNotPresent(cm.Contributors, &cm.ContributorsCount, userID)
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
			increaseCountIfNotPresent(m.ClaManagers, &m.ClaManagersCount, acl)
		}
	case EmployeeSignature:
		companyID := sig.SignatureUserCompanyID
		m, ok := cm.CompanyMetrics[companyID]
		if !ok {
			m = newCompanyMetric()
			cm.CompanyMetrics[companyID] = m
		}
		userID := sig.SignatureReferenceID
		increaseCountIfNotPresent(m.CorporateContributors, &m.CorporateContributorsCount, userID)
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
		increaseCountIfNotPresent(m.Companies, &m.CompaniesCount, companyID)
		for _, acl := range sig.SignatureACL {
			increaseCountIfNotPresent(m.ClaManagers, &m.ClaManagersCount, acl)
		}
	case EmployeeSignature:
		userID := sig.SignatureReferenceID
		increaseCountIfNotPresent(m.CorporateContributors, &m.CorporateContributorsCount, userID)
	case IclaSignature:
		userID := sig.SignatureReferenceID
		increaseCountIfNotPresent(m.IndividualContributors, &m.IndividualContributorsCount, userID)
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

func getClaManagerDistribution(cm *CompanyMetrics) *ClaManagersDistribution {
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
	metrics.ClaManagersDistribution = getClaManagerDistribution(metrics.CompanyMetrics)
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
			metrics.ProjectMetrics.processRepositories(r)
		}

		if len(results.LastEvaluatedKey) != 0 {
			scanInput.ExclusiveStartKey = results.LastEvaluatedKey
		} else {
			break
		}
	}
	metrics.ClaManagersDistribution = getClaManagerDistribution(metrics.CompanyMetrics)
	return nil
}

func (repo *repo) getMetrics() (*Metrics, error) {
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
	_, metrics.CalculatedAt = utils.CurrentTime()
	log.Println("GetMetrics took time", time.Since(t).String())
	return metrics, nil
}

func (repo *repo) saveMetrics(metrics *Metrics) error {
	t := time.Now()
	err := repo.saveTotalMetris(metrics.TotalMetrics)
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

func addIdTypeTime(item map[string]*dynamodb.AttributeValue, id string, metricType string) {
	_, ctime := utils.CurrentTime()
	utils.AddStringAttribute(item, "id", id)
	utils.AddStringAttribute(item, "metric_type", metricType)
	utils.AddStringAttribute(item, "created_at", ctime)
}

func (repo *repo) saveTotalMetris(tm *TotalMetrics) error {
	log.Println("saving total_metrics")
	av, err := dynamodbattribute.MarshalMap(tm)
	if err != nil {
		return err
	}
	addIdTypeTime(av, "total_metrics", "total_metrics")
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
		addIdTypeTime(av, id, "company_metric")
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
	addIdTypeTime(av, "cla_managers_distribution", "cla_managers_distribution")
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
		addIdTypeTime(av, id, "project_metric")
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
	m, err := repo.getMetrics()
	if err != nil {
		return err
	}
	err = repo.saveMetrics(m)
	if err != nil {
		return err
	}
	return nil
}
