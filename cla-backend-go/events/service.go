package events

import (
	"context"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/events"
)

// Service interface defines methods of event service
type Service interface {
	CreateEvent(event models.Event) error
	SearchEvents(ctx context.Context, params *events.SearchEventsParams) (*models.EventList, error)
}

type service struct {
	repo Repository
}

// NewService creates new instance of event service
func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) CreateEvent(event models.Event) error {
	return s.repo.CreateEvent(&event)
}

func (s *service) SearchEvents(ctx context.Context, params *events.SearchEventsParams) (*models.EventList, error) {
	return s.repo.SearchEvents(ctx, params)
}
