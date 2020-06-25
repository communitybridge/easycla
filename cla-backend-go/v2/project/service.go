package project

import (
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	v1ProjectParams "github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/project"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	v1Project "github.com/communitybridge/easycla/cla-backend-go/project"
	v2ProjectService "github.com/communitybridge/easycla/cla-backend-go/v2/project-service"
	v2ProjectServiceModels "github.com/communitybridge/easycla/cla-backend-go/v2/project-service/models"
)

// Service interface defines the v2 project service methods
type Service interface {
	GetCLAProjectsByID(foundationSFID string) (*models.EnabledClaList, error)
}

// service
type service struct {
	projectRepo v1Project.ProjectRepository
}

// NewService returns an instance of v2 project service
func NewService(projectRepo v1Project.ProjectRepository) Service {
	return &service{
		projectRepo: projectRepo,
	}
}

func (s *service) GetCLAProjectsByID(foundationSFID string) (*models.EnabledClaList, error) {
	var enabledClas []*models.EnabledCla
	psc := v2ProjectService.GetClient()
	projectDetails, err := psc.GetProject(foundationSFID)
	if err != nil {
		log.Error("Unable to get project details", err)
		return nil, err
	}

	// Filter claProjects and convert to response model
	enabledClas = s.filterClaProjects(projectDetails.Projects)
	log.Warnf("output: %+v", enabledClas)

	// sort by project names
	sort.Slice(enabledClas, func(i, j int) bool {
		return enabledClas[i].ProjectName < enabledClas[j].ProjectName
	})

	return &models.EnabledClaList{List: enabledClas}, nil
}

// filterClaProjects return projects output for which cla_group is present in cla
func (s *service) filterClaProjects(projects []*v2ProjectServiceModels.ProjectOutput) []*models.EnabledCla {
	results := make([]*models.EnabledCla, 0)
	prChan := make(chan *models.EnabledCla)

	for _, v := range projects {
		go func(projectOutput *v2ProjectServiceModels.ProjectOutput) {
			project, err := s.projectRepo.GetProjectsByExternalID(&v1ProjectParams.GetProjectsByExternalIDParams{
				ProjectSFID: projectOutput.ID,
				PageSize:    aws.Int64(1),
			}, false)
			if err != nil {
				log.Warnf("Unable to fetch project details for project with external id %s. error = %s", projectOutput.ID, err)
				prChan <- nil
				return
			}
			if project.ResultCount == 0 {
				prChan <- nil
				return
			}
			// Check if project is cla enabled
			claProject := project.Projects[0]
			if !claProject.ProjectCCLAEnabled && !claProject.ProjectICLAEnabled {
				prChan <- nil
				return
			}

			prChan <- &models.EnabledCla{
				ProjectName:      projectOutput.Name,
				ProjectSfid:      projectOutput.ID,
				ProjectLogo:      projectOutput.ProjectLogo,
				ProjectType:      projectOutput.ProjectType,
				CclaEnabled:      aws.Bool(claProject.ProjectCCLAEnabled),
				IclaEnabled:      aws.Bool(claProject.ProjectICLAEnabled),
				CclaRequiresIcla: aws.Bool(claProject.ProjectCCLARequiresICLA),
			}
		}(v)
	}

	for range projects {
		project := <-prChan
		if project != nil {
			log.Warnf("Adding project to channel: %+v", project)
			results = append(results, project)
		}
	}

	return results
}
