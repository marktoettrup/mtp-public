package ntnx

import (
	"context"
	"fmt"
	"strings"

	"bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/pkg/clients/ntnxv3Client/ntnxv3ClientProjects"
	"bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/pkg/clients/ntnxv3Client/ntnxv3models"
)

func NewProjectValidator(projectName string, client ntnxv3ClientProjects.ProjectsService) *ProjectValidator {
	return &ProjectValidator{
		projectName: projectName,
		api:         client,
	}
}

func (v *ProjectValidator) Name() string {
	return "project/" + v.projectName
}

func (v *ProjectValidator) Validate(ctx context.Context) error {
	cleanProjectName := strings.TrimSpace(strings.ToLower(v.projectName))
	project, err := v.GetResource(ctx)
	if err != nil {
		return fmt.Errorf("error fetching project resource: %w", err)
	}

	ntnxProject, ok := project.(*ntnxv3models.ProjectIntentResource)
	if !ok {
		return fmt.Errorf("expected ProjectIntentResource, got %T", project)
	}

	if strings.ToLower(strings.TrimSpace(ntnxProject.Metadata.Name)) == cleanProjectName {
		return nil
	}

	return fmt.Errorf("project %s not found", v.projectName)
}

func (v *ProjectValidator) GetResource(ctx context.Context) (interface{}, error) {
	listParams := getProjectListParams(ctx, v.projectName)
	projectList, err := v.api.PostProjectsList(listParams, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching project list: %w", err)
	}

	for _, project := range projectList.Payload.Entities {
		if project.Status != nil && project.Status.Name != nil && *project.Status.Name == v.projectName {
			return project, nil
		}
	}

	return nil, fmt.Errorf("project '%s' not found", v.projectName)
}

func (v *ProjectValidator) ValidateProject(ntnxProject *ntnxv3models.ProjectIntentResource, projectName string) (bool, error) {
	if ntnxProject.Status != nil && ntnxProject.Status.Name != nil && *ntnxProject.Status.Name == projectName {
		fmt.Printf("Found project: %s\n", projectName)

		summary, err := v.GetResourceHeadroomSummary(ntnxProject)
		if err != nil {
			return true, fmt.Errorf("failed to get resource summary: %w", err)
		}

		v.PrintResourceHeadroomSummary(summary)

		if err := v.ValidateResourceThresholds(summary, 80.0, 95.0); err != nil {
			fmt.Printf("WARNING: %v\n", err)
		}

		fmt.Printf("Project %s validation completed\n", projectName)
		return true, nil
	}

	return false, nil
}

func (v *ProjectValidator) ValidateProjectForWorkload(ctx context.Context, request ResourceRequest) error {
	listParams := getProjectListParams(ctx, v.projectName)
	projectList, err := v.api.PostProjectsList(listParams, nil)
	if err != nil {
		return fmt.Errorf("error fetching project list: %w", err)
	}

	for _, project := range projectList.Payload.Entities {
		if project.Status != nil && project.Status.Name != nil && *project.Status.Name == v.projectName {
			fmt.Printf("Found project: %s\n", v.projectName)

			result, err := v.CheckResourceAvailability(project, request)
			if err != nil {
				return fmt.Errorf("failed to check resource availability: %w", err)
			}

			v.PrintResourceAvailabilityResult(result)

			if !result.CanProvision {
				return fmt.Errorf("insufficient resources in project '%s' to provision requested workload", v.projectName)
			}

			fmt.Printf("Project %s has sufficient resources for the requested workload\n", v.projectName)
			return nil
		}
	}

	return fmt.Errorf("project '%s' not found", v.projectName)
}

func getProjectListParams(ctx context.Context, projectName string) *ntnxv3ClientProjects.PostProjectsListParams {
	return ntnxv3ClientProjects.NewPostProjectsListParams().
		WithContext(ctx).
		WithGetEntitiesRequest(&ntnxv3models.ProjectListMetadata{
			Filter:    "",
			Offset:    func(i int32) *int32 { return &i }(0),
			Length:    100,
			SortOrder: "ASCENDING",
		})
}
