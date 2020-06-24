package signatures

import (
	"bytes"
	"fmt"

	"github.com/jinzhu/copier"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/communitybridge/easycla/cla-backend-go/company"
	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
	v1Project "github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/project"
	v1Signatures "github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/signatures"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/project"
	"github.com/communitybridge/easycla/cla-backend-go/signatures"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
	"github.com/sirupsen/logrus"
)

// constants
const (
	// used when we want to query all data from dependent service.
	HugePageSize      = int64(10000)
	CclaSignatureType = "ccla"
)

type service struct {
	v1ProjectService   project.Service
	v1CompanyService   company.IService
	v1SignatureService signatures.SignatureService
}

// Service contains method of v2 signature service
type Service interface {
	GetProjectCompanySignatures(companySFID string, projectSFID string) (*models.Signatures, error)
	GetProjectIclaSignaturesCsv(claGroupID string) ([]byte, error)
	GetProjectIclaSignatures(claGroupID string, searchTerm *string) (*models.IclaSignatures, error)
}

// NewService creates instance of v2 signature service
func NewService(v1ProjectService project.Service,
	v1CompanyService company.IService,
	v1SignatureService signatures.SignatureService) *service {
	return &service{
		v1ProjectService:   v1ProjectService,
		v1CompanyService:   v1CompanyService,
		v1SignatureService: v1SignatureService,
	}
}

func (s *service) GetProjectCompanySignatures(companySFID string, projectSFID string) (*models.Signatures, error) {
	companyModel, err := s.v1CompanyService.GetCompanyByExternalID(companySFID)
	if err != nil {
		return nil, err
	}
	projects, err := s.v1ProjectService.GetProjectsByExternalID(&v1Project.GetProjectsByExternalIDParams{
		PageSize:    aws.Int64(HugePageSize),
		ProjectSFID: projectSFID,
	})
	if err != nil {
		return nil, err
	}
	projectIDs := utils.NewStringSet()
	for _, p := range projects.Projects {
		projectIDs.Add(p.ProjectID)
	}
	sigs, err := s.v1SignatureService.GetCompanySignatures(v1Signatures.GetCompanySignaturesParams{
		CompanyID:     companyModel.CompanyID,
		PageSize:      aws.Int64(HugePageSize),
		SignatureType: aws.String(CclaSignatureType),
	})
	if err != nil {
		return nil, err
	}
	filteredSigs := &v1Models.Signatures{
		LastKeyScanned: sigs.LastKeyScanned,
		ResultCount:    0,
		Signatures:     nil,
	}
	for _, sig := range sigs.Signatures {
		if projectIDs.Include(sig.ProjectID) {
			filteredSigs.Signatures = append(filteredSigs.Signatures, sig)
			filteredSigs.ResultCount++
		}
	}
	return v2SignaturesReplaceCompanyID(filteredSigs, companyModel.CompanyID, companySFID)
}

func iclaSigCsvLine(sig *v1Models.IclaSignature) string {
	var dateTime string
	t, err := utils.ParseDateTime(sig.SignedOn)
	if err != nil {
		log.WithFields(logrus.Fields{"signature_id": sig.SignatureID, "signature_created": sig.SignedOn}).
			Error("invalid time format present for signatures")
	} else {
		dateTime = t.Format("Jan 2,2006")
	}
	return fmt.Sprintf("\n%s,%s,%s,%s,\"%s\"", sig.GithubUsername, sig.LfUsername, sig.UserName, sig.UserEmail, dateTime)
}

func (s service) GetProjectIclaSignaturesCsv(claGroupID string) ([]byte, error) {
	var b bytes.Buffer
	result, err := s.v1SignatureService.GetClaGroupICLASignatures(claGroupID, nil)
	if err != nil {
		return nil, err
	}
	b.WriteString(`Github ID,LF_ID,Name,Email,Date Signed`)
	for _, sig := range result.List {
		b.WriteString(iclaSigCsvLine(sig))
	}
	return b.Bytes(), nil
}

func (s service) GetProjectIclaSignatures(claGroupID string, searchTerm *string) (*models.IclaSignatures, error) {
	var out models.IclaSignatures
	result, err := s.v1SignatureService.GetClaGroupICLASignatures(claGroupID, searchTerm)
	if err != nil {
		return nil, err
	}
	err = copier.Copy(&out, result)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
