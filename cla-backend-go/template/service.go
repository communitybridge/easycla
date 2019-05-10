package template

import (
	"context"

	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/gen/models"
)

type Service interface {
	GetTemplates(ctx context.Context) ([]models.Template, error)
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
