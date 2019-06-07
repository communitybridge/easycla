package health

import (
	"context"
	"time"

	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/gen/models"
	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/gen/restapi/operations"
)

// service handles async log of audit event
type service struct {
	//	db         *sqlx.DB
	gitHash    string
	buildStamp string
}

// New is a simple helper function to create a service instance
func New(GitHash, BuildStamp string) service {
	return service{
		//	db:         db,
		gitHash:    GitHash,
		buildStamp: BuildStamp,
	}
}

func (s service) HealthCheck(ctx context.Context, in operations.HealthCheckParams) (*models.Health, error) {
	t := time.Now()

	//var pong string
	//	dbErr := s.db.Get(&pong, "select 'pong'")

	duration := time.Since(t)

	hs := models.HealthStatus{TimeStamp: time.Now().String(), Healthy: true, Name: "CLA", Duration: duration.String()}
	// if dbErr != nil {
	// 	hs.Healthy = false
	// 	hs.Error = dbErr.Error()
	// }

	response := models.Health{
		Status:         "healthy",
		TimeStamp:      time.Now().String(),
		Githash:        s.gitHash,
		BuildTimeStamp: s.buildStamp,
		Healths:        []*models.HealthStatus{&hs},
	}

	return &response, nil
}
