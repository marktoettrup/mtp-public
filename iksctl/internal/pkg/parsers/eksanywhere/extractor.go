package parser

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func (p *Parser) ExtractClusterInfo(cluster *unstructured.Unstructured) (ClusterInfo, error) {
	info := ClusterInfo{
		Name:      cluster.GetName(),
		Namespace: cluster.GetNamespace(),
	}

	// Extract datacenterRef
	datacenterRefKind, found, err := unstructured.NestedString(cluster.Object, "spec", "datacenterRef", "kind")
	if err != nil {
		return info, fmt.Errorf("failed to extract datacenterRef kind: %w", err)
	}
	if found {
		info.DatacenterRef.Kind = datacenterRefKind
	}

	datacenterRefName, found, err := unstructured.NestedString(cluster.Object, "spec", "datacenterRef", "name")
	if err != nil {
		return info, fmt.Errorf("failed to extract datacenterRef name: %w", err)
	}
	if found {
		info.DatacenterRef.Name = datacenterRefName
	}

	// Extract gitOpsRef
	gitOpsRefKind, found, err := unstructured.NestedString(cluster.Object, "spec", "gitOpsRef", "kind")
	if err != nil {
		return info, fmt.Errorf("failed to extract gitOpsRef kind: %w", err)
	}
	if found {
		info.GitOpsRef.Kind = gitOpsRefKind
	}

	gitOpsRefName, found, err := unstructured.NestedString(cluster.Object, "spec", "gitOpsRef", "name")
	if err != nil {
		return info, fmt.Errorf("failed to extract gitOpsRef name: %w", err)
	}
	if found {
		info.GitOpsRef.Name = gitOpsRefName
	}

	// Extract controlPlane machineGroupRef
	controlPlaneMachineRefKind, found, err := unstructured.NestedString(cluster.Object, "spec", "controlPlaneConfiguration", "machineGroupRef", "kind")
	if err != nil {
		return info, fmt.Errorf("failed to extract controlPlane machineGroupRef kind: %w", err)
	}
	if found {
		info.ControlPlaneMachineRef.Kind = controlPlaneMachineRefKind
	}

	controlPlaneMachineRefName, found, err := unstructured.NestedString(cluster.Object, "spec", "controlPlaneConfiguration", "machineGroupRef", "name")
	if err != nil {
		return info, fmt.Errorf("failed to extract controlPlane machineGroupRef name: %w", err)
	}
	if found {
		info.ControlPlaneMachineRef.Name = controlPlaneMachineRefName
	}

	controlPlaneCount, found, err := unstructured.NestedInt64(cluster.Object, "spec", "controlPlaneConfiguration", "count")
	if err != nil {
		return info, fmt.Errorf("failed to extract controlPlane count: %w", err)
	}
	if found {
		info.ControlPlaneCount = controlPlaneCount
	}

	workerGroups, found, err := unstructured.NestedSlice(cluster.Object, "spec", "workerNodeGroupConfigurations")
	if err != nil {
		return info, fmt.Errorf("failed to extract worker node groups: %w", err)
	}
	if found {
		for i, group := range workerGroups {
			if groupMap, ok := group.(map[string]interface{}); ok {
				machineRefKind, found, err := unstructured.NestedString(groupMap, "machineGroupRef", "kind")
				if err != nil {
					return info, fmt.Errorf("failed to extract worker group %d machineGroupRef kind: %w", i, err)
				}

				machineRefName, found2, err := unstructured.NestedString(groupMap, "machineGroupRef", "name")
				if err != nil {
					return info, fmt.Errorf("failed to extract worker group %d machineGroupRef name: %w", i, err)
				}

				workerCount, found3, err := unstructured.NestedInt64(groupMap, "count")
				if err != nil {
					return info, fmt.Errorf("failed to extract worker group %d count: %w", i, err)
				}

				if found && found2 && found3 {
					info.WorkerMachineRefs = append(info.WorkerMachineRefs, WorkerGroupRef{
						ObjectRef: ObjectRef{
							Kind: machineRefKind,
							Name: machineRefName,
						},
						Count: workerCount,
					})
				}
			}
		}
	}

	endpoint, found, err := unstructured.NestedString(cluster.Object, "spec", "controlPlaneConfiguration", "endpoint", "host")
	if err != nil {
		return info, fmt.Errorf("failed to extract cluster endpoint: %w", err)
	}

	if found {
		info.Endpoint = endpoint
	}

	version, found, err := unstructured.NestedString(cluster.Object, "spec", "kubernetesVersion")
	if err != nil {
		return info, fmt.Errorf("failed to extract kubernetes version: %w", err)
	}
	if found {
		info.KubernetesVersion = version
	}

	eksaVersion, found, err := unstructured.NestedString(cluster.Object, "spec", "eksaVersion")
	if err != nil {
		return info, fmt.Errorf("failed to extract EKS Anywhere version: %w", err)
	}
	if found {
		info.EksaVersion = eksaVersion
	}

	// Extract cluster network configuration
	networkInfo, err := p.extractClusterNetworkInfo(cluster)
	if err != nil {
		return info, fmt.Errorf("failed to extract cluster network info: %w", err)
	}
	info.ClusterNetwork = networkInfo

	return info, nil
}

func (p *Parser) ExtractNutanixMachineConfigInfo(config *unstructured.Unstructured) (NutanixMachineConfigInfo, error) {
	info := NutanixMachineConfigInfo{
		Name:      config.GetName(),
		Namespace: config.GetNamespace(),
	}

	vcpusPerSocket, found, err := unstructured.NestedInt64(config.Object, "spec", "vcpusPerSocket")
	if err != nil {
		return info, fmt.Errorf("failed to extract vcpusPerSocket: %w", err)
	}
	if found {
		info.VCPUsPerSocket = vcpusPerSocket
	}

	vcpuSockets, found, err := unstructured.NestedInt64(config.Object, "spec", "vcpuSockets")
	if err != nil {
		return info, fmt.Errorf("failed to extract vcpuSockets: %w", err)
	}
	if found {
		info.VCPUSockets = vcpuSockets
	}

	info.VCPUs = info.VCPUsPerSocket * info.VCPUSockets

	memorySize, found, err := unstructured.NestedString(config.Object, "spec", "memorySize")
	if err != nil {
		return info, fmt.Errorf("failed to extract memory size: %w", err)
	}
	if found {
		info.MemorySize = memorySize
	}

	systemDiskSize, found, err := unstructured.NestedString(config.Object, "spec", "systemDiskSize")
	if err != nil {
		return info, fmt.Errorf("failed to extract disk size: %w", err)
	}
	if found {
		info.SystemDiskSize = systemDiskSize
	}

	nutanixCluster, found, err := unstructured.NestedString(config.Object, "spec", "cluster", "name")
	if err != nil {
		return info, fmt.Errorf("failed to extract nutanix cluster: %w", err)
	}
	if found {
		info.NutanixCluster = nutanixCluster
	}

	imageName, found, err := unstructured.NestedString(config.Object, "spec", "image", "name")
	if err != nil {
		return info, fmt.Errorf("failed to extract image name: %w", err)
	}
	if found {
		info.ImageName = imageName
	}

	subnetName, found, err := unstructured.NestedString(config.Object, "spec", "subnet", "name")
	if err != nil {
		return info, fmt.Errorf("failed to extract subnet name: %w", err)
	}
	if found {
		info.SubnetName = subnetName
	}

	nutanixProjectName, found, err := unstructured.NestedString(config.Object, "spec", "project", "name")
	if err != nil {
		return info, fmt.Errorf("failed to extract nutanix project name: %w", err)
	}
	if found {
		info.NutanixProjectName = nutanixProjectName
	}

	return info, nil
}

func (p *Parser) GetClusterSpec(cluster *unstructured.Unstructured) (map[string]interface{}, error) {
	spec, found, err := unstructured.NestedMap(cluster.Object, "spec")
	if err != nil {
		return nil, fmt.Errorf("failed to extract cluster spec: %w", err)
	}
	if !found {
		return nil, fmt.Errorf("cluster spec not found")
	}
	return spec, nil
}

func (p *Parser) GetWorkerNodeGroups(cluster *unstructured.Unstructured) ([]map[string]interface{}, error) {
	workerGroups, found, err := unstructured.NestedSlice(cluster.Object, "spec", "workerNodeGroupConfigurations")
	if err != nil {
		return nil, fmt.Errorf("failed to extract worker node groups: %w", err)
	}
	if !found {
		return []map[string]interface{}{}, nil
	}

	var groups []map[string]interface{}
	for _, group := range workerGroups {
		if groupMap, ok := group.(map[string]interface{}); ok {
			groups = append(groups, groupMap)
		}
	}
	return groups, nil
}

func (p *Parser) GetControlPlaneConfig(cluster *unstructured.Unstructured) (map[string]interface{}, error) {
	controlPlane, found, err := unstructured.NestedMap(cluster.Object, "spec", "controlPlaneConfiguration")
	if err != nil {
		return nil, fmt.Errorf("failed to extract control plane config: %w", err)
	}
	if !found {
		return nil, fmt.Errorf("control plane configuration not found")
	}
	return controlPlane, nil
}

func (p *Parser) extractClusterNetworkInfo(cluster *unstructured.Unstructured) (ClusterNetworkInfo, error) {
	info := ClusterNetworkInfo{}

	// Extract CNI configuration
	cniInfo, err := p.extractCNIConfigInfo(cluster)
	if err != nil {
		return info, fmt.Errorf("failed to extract CNI config: %w", err)
	}
	info.CNIConfig = cniInfo

	// Extract pods network configuration
	podsInfo, err := p.extractPodsNetworkInfo(cluster)
	if err != nil {
		return info, fmt.Errorf("failed to extract pods network config: %w", err)
	}
	info.Pods = podsInfo

	// Extract services network configuration
	servicesInfo, err := p.extractServicesNetworkInfo(cluster)
	if err != nil {
		return info, fmt.Errorf("failed to extract services network config: %w", err)
	}
	info.Services = servicesInfo

	return info, nil
}

func (p *Parser) extractCNIConfigInfo(cluster *unstructured.Unstructured) (CNIConfigInfo, error) {
	info := CNIConfigInfo{}

	// Extract Cilium configuration
	ciliumInfo, err := p.extractCiliumConfigInfo(cluster)
	if err != nil {
		return info, fmt.Errorf("failed to extract Cilium config: %w", err)
	}
	info.Cilium = ciliumInfo

	return info, nil
}

func (p *Parser) extractCiliumConfigInfo(cluster *unstructured.Unstructured) (CiliumConfigInfo, error) {
	info := CiliumConfigInfo{}

	// Extract policyEnforcementMode
	policyMode, found, err := unstructured.NestedString(cluster.Object, "spec", "clusterNetwork", "cniConfig", "cilium", "policyEnforcementMode")
	if err != nil {
		return info, fmt.Errorf("failed to extract Cilium policyEnforcementMode: %w", err)
	}
	if found {
		info.PolicyEnforcementMode = policyMode
	}

	// Extract egressMasqueradeInterfaces
	egressInterfaces, found, err := unstructured.NestedString(cluster.Object, "spec", "clusterNetwork", "cniConfig", "cilium", "egressMasqueradeInterfaces")
	if err != nil {
		return info, fmt.Errorf("failed to extract Cilium egressMasqueradeInterfaces: %w", err)
	}
	if found {
		info.EgressMasqueradeInterfaces = egressInterfaces
	}

	// Extract skipUpgrade
	skipUpgrade, found, err := unstructured.NestedBool(cluster.Object, "spec", "clusterNetwork", "cniConfig", "cilium", "skipUpgrade")
	if err != nil {
		return info, fmt.Errorf("failed to extract Cilium skipUpgrade: %w", err)
	}
	if found {
		info.SkipUpgrade = skipUpgrade
	}

	// Extract routingMode
	routingMode, found, err := unstructured.NestedString(cluster.Object, "spec", "clusterNetwork", "cniConfig", "cilium", "routingMode")
	if err != nil {
		return info, fmt.Errorf("failed to extract Cilium routingMode: %w", err)
	}
	if found {
		info.RoutingMode = routingMode
	}

	// Extract ipv4NativeRoutingCIDR
	ipv4CIDR, found, err := unstructured.NestedString(cluster.Object, "spec", "clusterNetwork", "cniConfig", "cilium", "ipv4NativeRoutingCIDR")
	if err != nil {
		return info, fmt.Errorf("failed to extract Cilium ipv4NativeRoutingCIDR: %w", err)
	}
	if found {
		info.IPv4NativeRoutingCIDR = ipv4CIDR
	}

	// Extract ipv6NativeRoutingCIDR
	ipv6CIDR, found, err := unstructured.NestedString(cluster.Object, "spec", "clusterNetwork", "cniConfig", "cilium", "ipv6NativeRoutingCIDR")
	if err != nil {
		return info, fmt.Errorf("failed to extract Cilium ipv6NativeRoutingCIDR: %w", err)
	}
	if found {
		info.IPv6NativeRoutingCIDR = ipv6CIDR
	}

	return info, nil
}

func (p *Parser) extractPodsNetworkInfo(cluster *unstructured.Unstructured) (NetworkRangeInfo, error) {
	info := NetworkRangeInfo{}

	// Extract pods CIDR blocks
	cidrBlocks, found, err := unstructured.NestedStringSlice(cluster.Object, "spec", "clusterNetwork", "pods", "cidrBlocks")
	if err != nil {
		return info, fmt.Errorf("failed to extract pods cidrBlocks: %w", err)
	}
	if found {
		info.CIDRBlocks = cidrBlocks
	}

	return info, nil
}

func (p *Parser) extractServicesNetworkInfo(cluster *unstructured.Unstructured) (NetworkRangeInfo, error) {
	info := NetworkRangeInfo{}

	// Extract services CIDR blocks
	cidrBlocks, found, err := unstructured.NestedStringSlice(cluster.Object, "spec", "clusterNetwork", "services", "cidrBlocks")
	if err != nil {
		return info, fmt.Errorf("failed to extract services cidrBlocks: %w", err)
	}
	if found {
		info.CIDRBlocks = cidrBlocks
	}

	return info, nil
}

func (p *Parser) GetClusterNetworkConfig(cluster *unstructured.Unstructured) (map[string]interface{}, error) {
	clusterNetwork, found, err := unstructured.NestedMap(cluster.Object, "spec", "clusterNetwork")
	if err != nil {
		return nil, fmt.Errorf("failed to extract cluster network config: %w", err)
	}
	if !found {
		return nil, fmt.Errorf("cluster network configuration not found")
	}
	return clusterNetwork, nil
}
