package parser

type ClusterInfo struct {
	Name                   string
	Namespace              string
	Endpoint               string
	KubernetesVersion      string
	EksaVersion            string
	DatacenterRef          ObjectRef
	GitOpsRef              ObjectRef
	ControlPlaneMachineRef ObjectRef
	ControlPlaneCount      int64
	WorkerMachineRefs      []WorkerGroupRef
	ClusterNetwork         ClusterNetworkInfo
}

type ObjectRef struct {
	Kind string
	Name string
}

type WorkerGroupRef struct {
	ObjectRef
	Count int64
}

type NutanixMachineConfigInfo struct {
	Name               string
	Namespace          string
	VCPUsPerSocket     int64
	VCPUSockets        int64
	VCPUs              int64 // Calculated field: VCPUsPerSocket * VCPUSockets
	MemorySize         string
	SystemDiskSize     string
	NutanixCluster     string
	ImageName          string
	SubnetName         string
	NutanixProjectName string
}

type ClusterNetworkInfo struct {
	CNIConfig CNIConfigInfo
	Pods      NetworkRangeInfo
	Services  NetworkRangeInfo
}

type CNIConfigInfo struct {
	Cilium CiliumConfigInfo
}

type CiliumConfigInfo struct {
	PolicyEnforcementMode      string
	EgressMasqueradeInterfaces string
	SkipUpgrade                bool
	RoutingMode                string
	IPv4NativeRoutingCIDR      string
	IPv6NativeRoutingCIDR      string
}

type NetworkRangeInfo struct {
	CIDRBlocks []string
}
