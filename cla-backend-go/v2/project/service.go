package project

import (
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
	v1ProjectParams "github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/project"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	v2ProjectService "github.com/communitybridge/easycla/cla-backend-go/v2/project-service"
	v2ProjectServiceModels "github.com/communitybridge/easycla/cla-backend-go/v2/project-service/models"
)

// Service interface defines the v2 project service methods
type Service interface {
	GetCLAProjectsByID(projectSfdcID string) (*models.EnabledClaList, error)
}

//Repo interface defines Project repository methods in v1
type Repo interface {
	GetProjectsByExternalID(params *v1ProjectParams.GetProjectsByExternalIDParams, loadRepoDetails bool) (*v1Models.Projects, error)
}

// service
type service struct {
	projectRepo Repo
}

// NewService returns an instance of v2 project service
func NewService(projectRepo Repo) Service {
	return &service{
		projectRepo: projectRepo,
	}
}
func (s *service) GetCLAProjectsByID(projectSfdcID string) (*models.EnabledClaList, error) {
	var enabledClas []*models.EnabledCla
	psc := v2ProjectService.GetClient()
	projectDetails, err := psc.GetProject(projectSfdcID)
	enabledClas = make([]*models.EnabledCla, 0, len(projectDetails.Projects))
	if err != nil {
		log.Error("Unable to get project details", err)
		return nil, err
	}

	// Filter claProjects from project service Projects
	subProjects := s.filterClaProjects(projectDetails.Projects)
	for _, subProject := range subProjects {
		sp := &models.EnabledCla{
			ProjectName: subProject.Name,
			ProjectSfid: subProject.ID,
			ProjectLogo: subProject.ProjectLogo,
			ProjectType: subProject.ProjectType,
		}
		enabledClas = append(enabledClas, sp)
	}

	// sort by project names
	sort.Slice(enabledClas, func(i, j int) bool {
		return enabledClas[i].ProjectName < enabledClas[j].ProjectName
	})
	return &models.EnabledClaList{List: enabledClas}, nil

}

// return projects output for which cla_group is present in cla
func (s *service) filterClaProjects(projects []*v2ProjectServiceModels.ProjectOutput) []*v2ProjectServiceModels.ProjectOutput {
	results := make([]*v2ProjectServiceModels.ProjectOutput, 0)
	prChan := make(chan *v2ProjectServiceModels.ProjectOutput)
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
			prChan <- projectOutput
		}(v)
	}
	for range projects {
		project := <-prChan
		if project != nil {
			results = append(results, project)
		}
	}
	return results
}
