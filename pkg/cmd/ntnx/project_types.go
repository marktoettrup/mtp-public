package ntnx

type ResourceHeadroomSummary struct {
	ProjectName string
	VCPUs       ResourceUsage
	Memory      ResourceUsage
	Storage     ResourceUsage
}

type ResourceUsage struct {
	Used         int64
	Limit        int64
	Available    int64
	UsagePercent float64
	Units        string
	IsUnlimited  bool
}

type ResourceRequest struct {
	VCPUs   int64
	Memory  int64
	Storage int64
}

type ResourceAvailabilityResult struct {
	ProjectName  string
	Request      ResourceRequest
	VCPUs        ResourceCheck
	Memory       ResourceCheck
	Storage      ResourceCheck
	CanProvision bool
}

type ResourceCheck struct {
	ResourceType string
	Requested    int64
	Current      ResourceUsage
	Available    bool
	Message      string
}

func NewResourceRequest(vcpus, memoryGB, storageGB int64) ResourceRequest {
	return ResourceRequest{
		VCPUs:   vcpus,
		Memory:  memoryGB,
		Storage: storageGB,
	}
}

func NewResourceRequestFromNodeSpec(nodeCount, vcpuPerNode, memoryGBPerNode, storageGBPerNode int64) ResourceRequest {
	return ResourceRequest{
		VCPUs:   nodeCount * vcpuPerNode,
		Memory:  nodeCount * memoryGBPerNode,
		Storage: nodeCount * storageGBPerNode,
	}
}
