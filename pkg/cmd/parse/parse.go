package parse

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"

	parser "bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/internal/pkg/parsers/eksanywhere"
	ctxhelpers "bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/pkg/context"
)

func runParse(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	logger, err := ctxhelpers.GetLogger(ctx)
	if err != nil {
		return fmt.Errorf("failed to get logger: %w", err)
	}

	filePath, err := cmd.Flags().GetString(flagFile)
	if err != nil {
		return fmt.Errorf("failed to get file flag: %w", err)
	}

	p := parser.New()

	clusterInfos, nutanixMachineConfigInfos, err := p.ParseStructure(filePath, logger)
	if err != nil {
		return fmt.Errorf("failed to parse file %s: %w", filePath, err)
	}

	if len(clusterInfos) == 0 && len(nutanixMachineConfigInfos) == 0 {
		logger.Info("No EKS Anywhere resources found in the provided file", "file", filePath)
	}

	DisplayNutanixResources(clusterInfos, nutanixMachineConfigInfos, logger)

	return nil
}

func DisplayNutanixResources(clusters []*parser.ClusterInfo, machineConfigs []*parser.NutanixMachineConfigInfo, logger *logr.Logger) {
	fmt.Println("=== Nutanix Resources Summary ===")

	if len(clusters) > 0 {
		fmt.Printf("Clusters: %d\n", len(clusters))
		displayClusterSummary(clusters, logger)
	}

	if len(machineConfigs) > 0 {
		fmt.Printf("Machine Configs: %d\n", len(machineConfigs))
		displayMachineConfigSummary(machineConfigs, logger)
	}
}

func displayClusterSummary(clusters []*parser.ClusterInfo, logger *logr.Logger) {
	if len(clusters) == 0 {
		logger.Info("No clusters found")
		return
	}

	logger.Info("Cluster summary", "count", len(clusters))
	for _, cluster := range clusters {
		logger.Info("Cluster details",
			"name", cluster.Name,
			"namespace", cluster.Namespace,
		)
	}
}

func displayMachineConfigSummary(machineConfigs []*parser.NutanixMachineConfigInfo, logger *logr.Logger) {
	if len(machineConfigs) == 0 {
		logger.Info("No machine configs found")
		return
	}

	logger.Info("Machine config summary", "count", len(machineConfigs))
	for _, config := range machineConfigs {
		logger.Info("Machine config details",
			"name", config.Name,
			"namespace", config.Namespace,
			"vcpus", config.VCPUs,
			"memorySize", config.MemorySize,
			"nutanixCluster", config.NutanixCluster,
			"imageName", config.ImageName,
			"subnetName", config.SubnetName,
			"nutanixProjectName", config.NutanixProjectName,
		)
	}
}
