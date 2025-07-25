package eksa

import (
	parser "bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/internal/pkg/parsers/eksanywhere"
	"bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/internal/pkg/validation"
)

type EksaValidationManager struct {
	validationManager       *validation.ValidationManager
	eksaVersion, k8sVersion string
}

type ClusterConfigValidator struct {
	clusterConfig  *parser.ClusterInfo
	machineConfigs []*parser.NutanixMachineConfigInfo
}

type MachineConfigValidator struct {
	machineConfig *parser.NutanixMachineConfigInfo
}
