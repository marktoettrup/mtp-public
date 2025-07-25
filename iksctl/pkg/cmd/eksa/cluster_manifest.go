package eksa

import (
	"context"
	"fmt"
	"net"

	parser "bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/internal/pkg/parsers/eksanywhere"
)

func NewClusterConfigValidator(clusterConfig *parser.ClusterInfo, machineConfigs []*parser.NutanixMachineConfigInfo) *ClusterConfigValidator {
	return &ClusterConfigValidator{
		clusterConfig:  clusterConfig,
		machineConfigs: machineConfigs,
	}
}

func (v *ClusterConfigValidator) Name() string {
	return "cluster-config/" + v.clusterConfig.Name
}

func (v *ClusterConfigValidator) Validate(ctx context.Context) error {
	if v.machineConfigs == nil {
		return fmt.Errorf("machine configs are nil")
	}

	if v.clusterConfig == nil {
		return fmt.Errorf("cluster config is nil")
	}

	if err := v.validateClusterManifest(); err != nil {
		return fmt.Errorf("cluster manifest validation failed: %w", err)
	}

	if err := v.addK8sVersionValidation(v.clusterConfig, v.machineConfigs); err != nil {
		return fmt.Errorf("Kubernetes version validation failed: %w", err)
	}

	if err := v.validateControlPlaneEndpointAvailable(); err != nil {
		return fmt.Errorf("control plane endpoint validation failed: %w", err)
	}

	if err := v.addCniPluginValidations(v.clusterConfig); err != nil {
		return fmt.Errorf("CNI plugin validation failed: %w", err)
	}

	return nil
}

func (v *ClusterConfigValidator) GetResource(ctx context.Context) (interface{}, error) {
	return v.clusterConfig, nil
}

func (v *ClusterConfigValidator) isPrivateIPv4(ip net.IP) bool {
	for _, privateNet := range privateNetworks {
		_, network, err := net.ParseCIDR(privateNet)
		if err != nil {
			continue
		}

		if network.Contains(ip) {
			return true
		}
	}

	return false
}

func (v *ClusterConfigValidator) validateClusterManifest() error {
	if v.clusterConfig == nil {
		return fmt.Errorf("cluster info is nil")
	}

	if v.clusterConfig.ControlPlaneCount < 3 {
		return fmt.Errorf("control plane must have at least 3 machines, got %d", v.clusterConfig.ControlPlaneCount)
	}

	machineConfigMap := make(map[string]*parser.NutanixMachineConfigInfo)
	for _, config := range v.machineConfigs {
		machineConfigMap[config.Name] = config
	}

	if v.clusterConfig.DatacenterRef.Name != "" {
		if v.clusterConfig.DatacenterRef.Kind != "NutanixDatacenterConfig" {
			return fmt.Errorf("unsupported datacenterRef kind: %s, expected NutanixDatacenterConfig", v.clusterConfig.DatacenterRef.Kind)
		}
	}

	if v.clusterConfig.GitOpsRef.Name != "" {
		if v.clusterConfig.GitOpsRef.Kind != "FluxConfig" {
			return fmt.Errorf("unsupported gitOpsRef kind: %s, expected FluxConfig", v.clusterConfig.GitOpsRef.Kind)
		}
	}

	if v.clusterConfig.ControlPlaneMachineRef.Name != "" {
		if v.clusterConfig.ControlPlaneMachineRef.Kind != "NutanixMachineConfig" {
			return fmt.Errorf("unsupported controlPlane machineGroupRef kind: %s, expected NutanixMachineConfig", v.clusterConfig.ControlPlaneMachineRef.Kind)
		}

		_, exists := machineConfigMap[v.clusterConfig.ControlPlaneMachineRef.Name]
		if !exists {
			return fmt.Errorf("control plane machine config %s not found", v.clusterConfig.ControlPlaneMachineRef.Name)
		}
	}

	for i, workerRef := range v.clusterConfig.WorkerMachineRefs {
		if workerRef.Kind != "NutanixMachineConfig" {
			return fmt.Errorf("unsupported worker node %d machineGroupRef kind: %s, expected NutanixMachineConfig", i, workerRef.Kind)
		}

		_, exists := machineConfigMap[workerRef.Name]
		if !exists {
			return fmt.Errorf("worker machine config %s not found", workerRef.Name)
		}
	}

	return nil
}
