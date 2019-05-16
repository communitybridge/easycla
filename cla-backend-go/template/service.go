package template

import (
	"context"

	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/gen/models"
)

type Service interface {
	GetTemplates(ctx context.Context) ([]models.Template, error)
	AddContractGroupTemplates(ctx context.Context, contractGroupID string) error
}

type service struct {
	templateRepo Repository
}

func NewService(templateRepo Repository) service {
	return service{
		templateRepo: templateRepo,
	}
}

func (s service) GetTemplates(ctx context.Context) ([]models.Template, error) {
	templates, err := s.templateRepo.GetTemplates(ctx)
	if err != nil {
		return nil, err
	}

	return templates, nil
}

func (s service) AddContractGroupTemplates(ctx context.Context, contractGroupID string) error {
	templates, err := s.templateRepo.GetTemplates(ctx)
	if err != nil {
		return err
	}

	template := templates[0]

	err = s.templateRepo.AddContractGroupTemplates(ctx, contractGroupID, template)
	if err != nil {
		return err
	}

	return nil

}
