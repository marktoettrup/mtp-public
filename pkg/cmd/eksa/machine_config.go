package eksa

import (
	"context"
	"fmt"

	parser "bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/internal/pkg/parsers/eksanywhere"
	"bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/internal/pkg/utils"
)

func NewMachineConfigValidator(machineConfig *parser.NutanixMachineConfigInfo) *MachineConfigValidator {
	return &MachineConfigValidator{
		machineConfig: machineConfig,
	}
}

func (v *MachineConfigValidator) Name() string {
	return "machine-config/" + v.machineConfig.Name
}

func (v *MachineConfigValidator) Validate(ctx context.Context) error {
	if v.machineConfig == nil {
		return fmt.Errorf("machine config is nil")
	}

	if err := v.validateCPUConfiguration(); err != nil {
		return fmt.Errorf("CPU configuration validation failed: %w", err)
	}

	if err := v.validateMemoryConfiguration(); err != nil {
		return fmt.Errorf("memory configuration validation failed: %w", err)
	}

	if err := v.validateDiskConfiguration(); err != nil {
		return fmt.Errorf("disk configuration validation failed: %w", err)
	}

	if err := v.validateResourceReferences(); err != nil {
		return fmt.Errorf("resource reference validation failed: %w", err)
	}

	return nil
}

func (v *MachineConfigValidator) GetResource(ctx context.Context) (interface{}, error) {
	return v.machineConfig, nil
}

func (v *MachineConfigValidator) validateCPUConfiguration() error {
	if v.machineConfig.VCPUSockets < 2 {
		return fmt.Errorf("vcpuSockets must be at least 2, got %d", v.machineConfig.VCPUSockets)
	}

	if v.machineConfig.VCPUsPerSocket < 1 {
		return fmt.Errorf("vcpusPerSocket must be at least 1, got %d", v.machineConfig.VCPUsPerSocket)
	}

	// Ensure minimum 2 vCPUs total
	totalVCPUs := v.machineConfig.VCPUSockets * v.machineConfig.VCPUsPerSocket
	if totalVCPUs < 2 {
		return fmt.Errorf("total vCPUs must be at least 2, got %d (sockets: %d, per socket: %d)",
			totalVCPUs, v.machineConfig.VCPUSockets, v.machineConfig.VCPUsPerSocket)
	}

	return nil
}

func (v *MachineConfigValidator) validateMemoryConfiguration() error {
	if v.machineConfig.MemorySize == "" {
		return fmt.Errorf("memorySize is required")
	}

	if gib, err := utils.ParseMemoryToGB(v.machineConfig.MemorySize); err != nil {
		return fmt.Errorf("invalid memorySize format '%s': %w", v.machineConfig.MemorySize, err)
	} else if gib < 4 {
		return fmt.Errorf("memorySize must be at least 4 GiB, got %d GiB", gib)
	}

	return nil
}

func (v *MachineConfigValidator) validateDiskConfiguration() error {
	if v.machineConfig.SystemDiskSize == "" {
		return fmt.Errorf("systemDiskSize is required")
	}

	if gib, err := utils.ParseStorageToGB(v.machineConfig.SystemDiskSize); err != nil {
		return fmt.Errorf("invalid systemDiskSize format '%s': %w", v.machineConfig.SystemDiskSize, err)
	} else if gib < 40 {
		return fmt.Errorf("systemDiskSize must be at least 40 GiB, got %d GiB", gib)
	}

	return nil
}

func (v *MachineConfigValidator) validateResourceReferences() error {
	if v.machineConfig.NutanixCluster == "" {
		return fmt.Errorf("nutanix cluster reference is required")
	}

	if v.machineConfig.ImageName == "" {
		return fmt.Errorf("image name reference is required")
	}

	if v.machineConfig.SubnetName == "" {
		return fmt.Errorf("subnet name reference is required")
	}

	if v.machineConfig.NutanixProjectName == "" {
		return fmt.Errorf("nutanix project name reference is required")
	}

	return nil
}
