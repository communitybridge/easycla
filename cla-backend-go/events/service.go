package events

import (
	"context"
	"strconv"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/events"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
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
	const defaultPageSize int64 = 50
	var pageSize = defaultPageSize

	if params.PageSize != nil {
		var err error
		pageSize, err = strconv.ParseInt(*params.PageSize, 10, 64)
		if err != nil {
			log.Warnf("error parsing pageSize parameter to int64 - using default size of %d - error: %v",
				defaultPageSize, err)
		}

		// Make sure it's positive
		if pageSize < 1 {
			log.Warnf("invalid page size of %d - must be a positive value - using default size of %d",
				pageSize, defaultPageSize)
			pageSize = defaultPageSize
		}
	}
	return s.repo.SearchEvents(ctx, params, pageSize)
}
