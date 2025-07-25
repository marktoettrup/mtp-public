package ntnx

import "fmt"

func (v *ProjectValidator) PrintResourceHeadroomSummary(summary *ResourceHeadroomSummary) {
	fmt.Printf("\n=== Resource Headroom Summary for Project: %s ===\n", summary.ProjectName)

	v.printResourceUsage("vCPU", summary.VCPUs)
	v.printResourceUsage("Memory", summary.Memory)
	v.printResourceUsage("Storage", summary.Storage)

	fmt.Printf("=== End of Summary ===\n\n")
}

func (v *ProjectValidator) printResourceUsage(resourceType string, usage ResourceUsage) {
	if usage.Used == 0 && usage.Limit == 0 {
		fmt.Printf("%s: No data available\n", resourceType)
		return
	}

	if usage.IsUnlimited {
		fmt.Printf("%s: %d %s (unlimited)\n", resourceType, usage.Used, usage.Units)
	} else {
		status := "✓"
		if usage.UsagePercent > 95 {
			status = "✗ CRITICAL"
		} else if usage.UsagePercent > 80 {
			status = "⚠ WARNING"
		}

		fmt.Printf("%s: %d/%d %s (%.1f%% used, %d %s available) %s\n",
			resourceType, usage.Used, usage.Limit, usage.Units,
			usage.UsagePercent, usage.Available, usage.Units, status)
	}
}

func (v *ProjectValidator) PrintResourceAvailabilityResult(result *ResourceAvailabilityResult) {
	fmt.Printf("\n=== Resource Availability Check for Project: %s ===\n", result.ProjectName)
	fmt.Printf("Requested Resources: %d vCPUs, %d GB Memory, %d GB Storage\n\n",
		result.Request.VCPUs, result.Request.Memory, result.Request.Storage)

	v.printResourceCheck(result.VCPUs)
	v.printResourceCheck(result.Memory)
	v.printResourceCheck(result.Storage)

	fmt.Printf("\n")
	if result.CanProvision {
		fmt.Printf("✅ RESULT: Resources are AVAILABLE - workload can be provisioned\n")
	} else {
		fmt.Printf("❌ RESULT: Insufficient resources - workload CANNOT be provisioned\n")
	}
	fmt.Printf("=== End of Availability Check ===\n\n")
}

func (v *ProjectValidator) printResourceCheck(check ResourceCheck) {
	status := "✅"
	if !check.Available {
		status = "❌"
	}
	fmt.Printf("%s %s\n", status, check.Message)
}
