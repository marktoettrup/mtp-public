package ntnx

import (
	"fmt"

	"bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/pkg/clients/ntnxv3Client/ntnxv3models"
)

func (v *ProjectValidator) GetResourceHeadroomSummary(project *ntnxv3models.ProjectIntentResource) (*ResourceHeadroomSummary, error) {
	if project.Status == nil || project.Status.Resources == nil || project.Status.Resources.ResourceDomain == nil {
		return nil, fmt.Errorf("project resource domain information is not available")
	}

	summary := &ResourceHeadroomSummary{
		ProjectName: v.projectName,
	}

	resourceDomain := project.Status.Resources.ResourceDomain
	if resourceDomain.Resources == nil {
		return summary, nil
	}

	for _, resource := range resourceDomain.Resources {
		if resource.ResourceType == nil || resource.Value == nil {
			continue
		}

		usage := ResourceUsage{
			Used:        *resource.Value,
			Limit:       resource.Limit,
			IsUnlimited: resource.Limit <= 0,
		}

		if resource.Units != nil {
			usage.Units = *resource.Units
		}

		if !usage.IsUnlimited {
			usage.Available = usage.Limit - usage.Used
			usage.UsagePercent = float64(usage.Used) / float64(usage.Limit) * 100
		}

		switch *resource.ResourceType {
		case "VCPUS":
			summary.VCPUs = usage
		case "MEMORY":
			summary.Memory = usage
		case "STORAGE":
			summary.Storage = usage
		}
	}

	return summary, nil
}

func (v *ProjectValidator) CheckResourceAvailability(project *ntnxv3models.ProjectIntentResource, request ResourceRequest) (*ResourceAvailabilityResult, error) {
	summary, err := v.GetResourceHeadroomSummary(project)
	if err != nil {
		return nil, fmt.Errorf("failed to get resource summary: %w", err)
	}

	result := &ResourceAvailabilityResult{
		ProjectName: summary.ProjectName,
		Request:     request,
		VCPUs:       v.checkResourceRequest("vCPU", request.VCPUs, summary.VCPUs),
		Memory:      v.checkResourceRequest("Memory", request.Memory, summary.Memory),
		Storage:     v.checkResourceRequest("Storage", request.Storage, summary.Storage),
	}

	result.CanProvision = result.VCPUs.Available && result.Memory.Available && result.Storage.Available

	return result, nil
}

func (v *ProjectValidator) checkResourceRequest(resourceType string, requested int64, current ResourceUsage) ResourceCheck {
	check := ResourceCheck{
		ResourceType: resourceType,
		Requested:    requested,
		Current:      current,
	}

	if current.IsUnlimited {
		check.Available = true
		check.Message = fmt.Sprintf("%s: %d requested (unlimited quota)", resourceType, requested)
		return check
	}

	if current.Limit == 0 {
		check.Available = false
		check.Message = fmt.Sprintf("%s: %d requested but no quota allocated", resourceType, requested)
		return check
	}

	check.Available = current.Available >= requested

	if check.Available {
		remainingAfter := current.Available - requested
		check.Message = fmt.Sprintf("%s: %d requested, %d available (%d %s remaining after provision)",
			resourceType, requested, current.Available, remainingAfter, current.Units)
	} else {
		shortage := requested - current.Available
		check.Message = fmt.Sprintf("%s: %d requested, only %d available (shortage: %d %s)",
			resourceType, requested, current.Available, shortage, current.Units)
	}

	return check
}

func (v *ProjectValidator) ValidateResourceThresholds(summary *ResourceHeadroomSummary, warningThreshold, criticalThreshold float64) error {
	var errors []string
	var warnings []string

	// Check vCPU
	if !summary.VCPUs.IsUnlimited && summary.VCPUs.UsagePercent > criticalThreshold {
		errors = append(errors, fmt.Sprintf("vCPU usage (%.1f%%) exceeds critical threshold (%.1f%%)", summary.VCPUs.UsagePercent, criticalThreshold))
	} else if !summary.VCPUs.IsUnlimited && summary.VCPUs.UsagePercent > warningThreshold {
		warnings = append(warnings, fmt.Sprintf("vCPU usage (%.1f%%) exceeds warning threshold (%.1f%%)", summary.VCPUs.UsagePercent, warningThreshold))
	}

	// Check Memory
	if !summary.Memory.IsUnlimited && summary.Memory.UsagePercent > criticalThreshold {
		errors = append(errors, fmt.Sprintf("Memory usage (%.1f%%) exceeds critical threshold (%.1f%%)", summary.Memory.UsagePercent, criticalThreshold))
	} else if !summary.Memory.IsUnlimited && summary.Memory.UsagePercent > warningThreshold {
		warnings = append(warnings, fmt.Sprintf("Memory usage (%.1f%%) exceeds warning threshold (%.1f%%)", summary.Memory.UsagePercent, warningThreshold))
	}

	// Check Storage
	if !summary.Storage.IsUnlimited && summary.Storage.UsagePercent > criticalThreshold {
		errors = append(errors, fmt.Sprintf("Storage usage (%.1f%%) exceeds critical threshold (%.1f%%)", summary.Storage.UsagePercent, criticalThreshold))
	} else if !summary.Storage.IsUnlimited && summary.Storage.UsagePercent > warningThreshold {
		warnings = append(warnings, fmt.Sprintf("Storage usage (%.1f%%) exceeds warning threshold (%.1f%%)", summary.Storage.UsagePercent, warningThreshold))
	}

	// Print warnings
	for _, warning := range warnings {
		fmt.Printf("WARNING: %s\n", warning)
	}

	// Return error if any critical thresholds are exceeded
	if len(errors) > 0 {
		return fmt.Errorf("critical resource thresholds exceeded: %v", errors)
	}

	return nil
}
