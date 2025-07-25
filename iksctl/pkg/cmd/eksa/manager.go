package eksa

import (
	"context"
	"fmt"

	parser "bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/internal/pkg/parsers/eksanywhere"
	"bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/internal/pkg/validation"
	"github.com/go-logr/logr"
)

func NewEksaValidationManager(logger *logr.Logger) *EksaValidationManager {
	return &EksaValidationManager{
		validationManager: validation.NewValidationManager(logger),
	}
}

func (evm *EksaValidationManager) Setup(ctx context.Context) error {
	// Setup logic for EksaValidationManager can be added here if needed.
	// Currently, it just initializes the validation manager.
	return nil
}

func (evm *EksaValidationManager) SetupWithVersionValidations(ctx context.Context, eksaVersion, k8sVersion string) error {
	evm.eksaVersion = eksaVersion
	evm.k8sVersion = k8sVersion

	return nil
}

func (evm *EksaValidationManager) AddClusterConfigValidator(clusterConfig *parser.ClusterInfo, machineConfigs []*parser.NutanixMachineConfigInfo) error {
	validator := NewClusterConfigValidator(clusterConfig, machineConfigs)
	evm.validationManager.AddValidator(validator)

	if evm.eksaVersion != "" && evm.k8sVersion != "" {
		if err := validator.addValidateEksaVersionSkew(evm.eksaVersion, clusterConfig.EksaVersion); err != nil {
			return fmt.Errorf("failed to add EKS Anywhere version skew validation: %w", err)
		}
		if err := validator.addValidateKubeVersionSkew(evm.k8sVersion, clusterConfig.KubernetesVersion); err != nil {
			return fmt.Errorf("failed to add Kubernetes version skew validation: %w", err)
		}
	}

	return nil
}

func (evm *EksaValidationManager) AddClusterMachineValidator(machineConfig *parser.NutanixMachineConfigInfo) error {
	validator := NewMachineConfigValidator(machineConfig)
	evm.validationManager.AddValidator(validator)

	return nil
}

func (evm *EksaValidationManager) Validate(ctx context.Context) error {
	return evm.validationManager.ValidateAll(ctx)
}
