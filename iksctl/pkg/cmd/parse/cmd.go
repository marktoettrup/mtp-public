package parse

import (
	"fmt"

	parser "bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/internal/pkg/parsers/eksanywhere"
	"bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/internal/pkg/utils"
	"bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/pkg/cmd/eksa"
	ntnx "bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/pkg/cmd/ntnx"
	ctxhelpers "bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/pkg/context"
	"github.com/spf13/cobra"
)

const (
	flagNtnxHost         = "host"
	flagNtnxPort         = "port"
	flagNtnxUsername     = "username"
	flagNtnxPassword     = "password"
	flagNtnxInsecure     = "insecure"
	flagFile             = "file"
	upgradeK8sVersionTo  = "upgrade-k8s-version-to"
	upgradeEksaVersionTo = "upgrade-eks-version-to"
)

func ValidateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate a Nutanix cluster",
		Long:  "Validate a Nutanix cluster",
		RunE:  runParseAndValidate,
	}

	cmd.Flags().StringP(flagFile, "f", "", "Path to the YAML manifest file to parse (required)")
	cmd.Flags().String(flagNtnxHost, "", "Nutanix Prism Central hostname or IP address")
	cmd.Flags().Int(flagNtnxPort, 9440, "Nutanix Prism Central port")
	cmd.Flags().String(flagNtnxUsername, "", "Nutanix Prism Central username")
	cmd.Flags().String(flagNtnxPassword, "", "Nutanix Prism Central password")
	cmd.Flags().Bool(flagNtnxInsecure, false, "Skip TLS verification")
	cmd.Flags().String(upgradeK8sVersionTo, "", "Upgrade Kubernetes version")
	cmd.Flags().String(upgradeEksaVersionTo, "", "Upgrade EKS Anywhere version")

	if err := utils.MarkRequiredFlags(cmd, flagFile, flagNtnxHost, flagNtnxUsername, flagNtnxPassword); err != nil {
		panic(err) // This should never happen during command setup
	}

	return cmd
}

func runParseAndValidate(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	logger, err := ctxhelpers.GetLogger(ctx)
	if err != nil {
		return fmt.Errorf("failed to get logger: %w", err)
	}

	p := parser.New()

	filePath, _ := cmd.Flags().GetString(flagFile)

	clusterInfos, nutanixMachineConfigInfos, err := p.ParseStructure(filePath, logger)
	if err != nil {
		return fmt.Errorf("failed to parse file %s: %w", filePath, err)
	}

	if len(clusterInfos) == 0 && len(nutanixMachineConfigInfos) == 0 {
		logger.Info("No EKS Anywhere resources found in the provided file", "file", filePath)
	}

	logger.Info("Parsing successful, Starting validation of Nutanix resources", "file", filePath)

	host, _ := cmd.Flags().GetString(flagNtnxHost)
	port, _ := cmd.Flags().GetInt(flagNtnxPort)
	username, _ := cmd.Flags().GetString(flagNtnxUsername)
	password, _ := cmd.Flags().GetString(flagNtnxPassword)
	insecure, _ := cmd.Flags().GetBool(flagNtnxInsecure)
	upgradeK8s, _ := cmd.Flags().GetString(upgradeK8sVersionTo)
	upgradeEksa, _ := cmd.Flags().GetString(upgradeEksaVersionTo)

	if (upgradeK8s != "" && upgradeEksa == "") || (upgradeK8s == "" && upgradeEksa != "") {
		return fmt.Errorf("both --upgrade-k8s-version-to and --upgrade-eks-version-to must be specified together")
	}

	config := ntnx.ConnectionConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Insecure: insecure,
	}

	nvm := ntnx.NewNutanixValidationManager(config, logger)
	if err := nvm.Setup(ctx); err != nil {
		return fmt.Errorf("failed to setup Nutanix validation manager: %w", err)
	}

	eksaVm := eksa.NewEksaValidationManager(logger)
	if upgradeK8s == "" && upgradeEksa == "" {
		if err := eksaVm.Setup(ctx); err != nil {
			return fmt.Errorf("failed to setup EKS validation manager: %w", err)
		}
	} else {
		if err := eksaVm.SetupWithVersionValidations(ctx, upgradeEksa, upgradeK8s); err != nil {
			return fmt.Errorf("failed to setup EKS validation manager with version validations: %w", err)
		}
	}

	if err := runValidations(ctx, nvm, eksaVm, clusterInfos, nutanixMachineConfigInfos); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parse",
		Short: "Parse EKS Anywhere manifests",
		Long:  "Parse and extract information from EKS Anywhere Cluster and NutanixMachineConfig manifests",
	}

	// Add subcommands
	cmd.AddCommand(AnalyzeCommand())
	cmd.AddCommand(ValidateCommand())

	return cmd
}

func AnalyzeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "analyze",
		Short: "Analyze EKS Anywhere manifests",
		Long:  "Parse and analyze EKS Anywhere Cluster and NutanixMachineConfig manifests",
		RunE:  runParse,
	}

	cmd.Flags().StringP(flagFile, "f", "", "Path to the YAML manifest file to parse (required)")
	err := cmd.MarkFlagRequired(flagFile)
	if err != nil {
		fmt.Printf("Error marking flag required: %v\n", err)
	}

	return cmd
}
