package parse

import (
	"context"
	"fmt"

	parser "bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/internal/pkg/parsers/eksanywhere"
	"bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/pkg/cmd/eksa"
	ntnx "bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/pkg/cmd/ntnx"
)

func runValidations(ctx context.Context, nvm *ntnx.NutanixValidationManager, eksaVm *eksa.EksaValidationManager, clusters []*parser.ClusterInfo, machineConfigs []*parser.NutanixMachineConfigInfo) error {
	if err := initClusterValidations(clusters, machineConfigs, eksaVm); err != nil {
		return fmt.Errorf("failed to initialize cluster validations: %w", err)
	}

	for _, machineConfig := range machineConfigs {
		if err := addMachineConfigValidation(machineConfig, nvm, eksaVm); err != nil {
			return fmt.Errorf("failed to add machine config validation for %s: %w", machineConfig.Name, err)
		}
	}

	if err := nvm.Validate(ctx); err != nil {
		return fmt.Errorf("nutanix validation failed: %w", err)
	}

	if err := eksaVm.Validate(ctx); err != nil {
		return fmt.Errorf("eks validation failed: %w", err)
	}

	return nil
}

func initClusterValidations(clusters []*parser.ClusterInfo, machineConfigs []*parser.NutanixMachineConfigInfo, eksaVm *eksa.EksaValidationManager) error {
	for _, cluster := range clusters {
		if err := eksaVm.AddClusterConfigValidator(cluster, machineConfigs); err != nil {
			return fmt.Errorf("failed to add cluster config validation for %s: %w", cluster.Name, err)
		}
	}

	return nil
}

func addMachineConfigValidation(machineConfig *parser.NutanixMachineConfigInfo, nvm *ntnx.NutanixValidationManager, eksaVm *eksa.EksaValidationManager) error {
	if machineConfig.NutanixProjectName != "" {
		nvm.AddProjectValidator(machineConfig.NutanixProjectName)
	}

	if machineConfig.ImageName != "" {
		nvm.AddImageValidator(machineConfig.ImageName)
	}

	if machineConfig.NutanixCluster != "" {
		nvm.AddClusterValidator(machineConfig.NutanixCluster)
	}

	if machineConfig.SubnetName != "" {
		nvm.AddSubnetValidator(machineConfig.SubnetName)
	}

	eksaVm.AddClusterMachineValidator(machineConfig)

	return nil
}
