package ntnx

import (
	"fmt"

	ctxhelpers "bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/pkg/context"
	"github.com/spf13/cobra"
)

const (
	flagNtnxHost          = "host"
	flagNtnxPort          = "port"
	flagNtnxUsername      = "username"
	flagNtnxPassword      = "password"
	flagNtnxPasswordKeyID = "password-key-id"
	flagNtnxInsecure      = "insecure"
	flagNtnxSubnet        = "subnet"
	flagNtnxImage         = "image"
)

var vm NutanixValidationManager

func NtnxCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ntnx",
		Short: "Nutanix utilities",
		Long:  "Commands for interacting with Nutanix Prism Central",
	}

	cmd.AddCommand(ValidatesubnetCommand())
	cmd.AddCommand(ValidateimageCommand())
	cmd.AddCommand(ValidateProjectCommand())
	cmd.AddCommand(ValidateClusterCommand())

	return cmd
}

func markRequiredFlags(cmd *cobra.Command, flags ...string) error {
	for _, flag := range flags {
		if err := cmd.MarkFlagRequired(flag); err != nil {
			return fmt.Errorf("failed to mark flag '%s' as required: %w", flag, err)
		}
	}
	return nil
}

func ValidatesubnetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate-subnet [subnet name]",
		Short: "Validate a Nutanix subnet",
		Long:  "Validate a Nutanix subnet",
		RunE:  vm.RunValidateSubnet,
		Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	}

	// Add flags
	cmd.Flags().String(flagNtnxHost, "", "Nutanix Prism Central hostname or IP address")
	cmd.Flags().Int(flagNtnxPort, 9440, "Nutanix Prism Central port")
	cmd.Flags().String(flagNtnxUsername, "", "Nutanix Prism Central username")
	cmd.Flags().String(flagNtnxPassword, "", "Nutanix Prism Central password")
	cmd.Flags().Int(flagNtnxPasswordKeyID, 0, "Password ID in PWDB for Nutanix credentials")
	cmd.Flags().Bool(flagNtnxInsecure, false, "Skip TLS verification")

	// Mark required flags
	if err := markRequiredFlags(cmd, flagNtnxHost, flagNtnxUsername, flagNtnxPassword); err != nil {
		panic(err) // This should never happen during command setup
	}

	return cmd
}

func ValidateimageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate-image [image name]",
		Short: "Validate a Nutanix image",
		Long:  "Validate a Nutanix image",
		RunE:  vm.RunValidateImage,
		Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	}

	cmd.Flags().String(flagNtnxHost, "", "Nutanix Prism Central hostname or IP address")
	cmd.Flags().Int(flagNtnxPort, 9440, "Nutanix Prism Central port")
	cmd.Flags().String(flagNtnxUsername, "", "Nutanix Prism Central username")
	cmd.Flags().String(flagNtnxPassword, "", "Nutanix Prism Central password")
	cmd.Flags().Int(flagNtnxPasswordKeyID, 0, "Password ID in PWDB for Nutanix credentials")
	cmd.Flags().Bool(flagNtnxInsecure, false, "Skip TLS verification")

	if err := markRequiredFlags(cmd, flagNtnxHost, flagNtnxUsername, flagNtnxPassword); err != nil {
		panic(err) // This should never happen during command setup
	}

	return cmd
}

func ValidateProjectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate-project [project name]",
		Short: "Validate a Nutanix project",
		Long:  "Validate a Nutanix project",
		RunE:  vm.RunValidateProject,
		Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	}

	cmd.Flags().String(flagNtnxHost, "", "Nutanix Prism Central hostname or IP address")
	cmd.Flags().Int(flagNtnxPort, 9440, "Nutanix Prism Central port")
	cmd.Flags().String(flagNtnxUsername, "", "Nutanix Prism Central username")
	cmd.Flags().String(flagNtnxPassword, "", "Nutanix Prism Central password")
	cmd.Flags().Int(flagNtnxPasswordKeyID, 0, "Password ID in PWDB for Nutanix credentials")
	cmd.Flags().Bool(flagNtnxInsecure, false, "Skip TLS verification")

	if err := markRequiredFlags(cmd, flagNtnxHost, flagNtnxUsername, flagNtnxPassword); err != nil {
		panic(err) // This should never happen during command setup
	}

	return cmd
}

func ValidateClusterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate-cluster [cluster name]",
		Short: "Validate a Nutanix cluster",
		Long:  "Validate a Nutanix cluster",
		RunE:  vm.RunValidateCluster,
		Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	}

	cmd.Flags().String(flagNtnxHost, "", "Nutanix Prism Central hostname or IP address")
	cmd.Flags().Int(flagNtnxPort, 9440, "Nutanix Prism Central port")
	cmd.Flags().String(flagNtnxUsername, "", "Nutanix Prism Central username")
	cmd.Flags().String(flagNtnxPassword, "", "Nutanix Prism Central password")
	cmd.Flags().Int(flagNtnxPasswordKeyID, 0, "Password ID in PWDB for Nutanix credentials")
	cmd.Flags().Bool(flagNtnxInsecure, false, "Skip TLS verification")

	if err := markRequiredFlags(cmd, flagNtnxHost, flagNtnxUsername, flagNtnxPassword); err != nil {
		panic(err) // This should never happen during command setup
	}

	return cmd
}

func (vm *NutanixValidationManager) RunValidateSubnet(cmd *cobra.Command, args []string) error {
	subnetName := args[0]

	// Get logger from context
	logger, err := ctxhelpers.GetLogger(cmd.Context())
	if err != nil {
		return fmt.Errorf("failed to get logger: %w", err)
	}

	host, _ := cmd.Flags().GetString(flagNtnxHost)
	port, _ := cmd.Flags().GetInt(flagNtnxPort)
	username, _ := cmd.Flags().GetString(flagNtnxUsername)
	password, _ := cmd.Flags().GetString(flagNtnxPassword)
	passwordKeyID, _ := cmd.Flags().GetInt(flagNtnxPasswordKeyID)
	insecure, _ := cmd.Flags().GetBool(flagNtnxInsecure)

	// Handle PWDB password retrieval if password key ID is provided
	if passwordKeyID > 0 {
		// TODO: Implement PWDB password retrieval
		// For now, we'll use the password flag
		if password == "" {
			return fmt.Errorf("password is required when not using PWDB")
		}
	}

	config := ConnectionConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Insecure: insecure,
	}

	manager := NewNutanixValidationManager(config, logger)

	if err := manager.Setup(cmd.Context()); err != nil {
		return fmt.Errorf("failed to setup validation manager: %w", err)
	}

	manager.AddSubnetValidator(subnetName)

	fmt.Printf("Validating subnet: %s\n", subnetName)
	if err := manager.Validate(cmd.Context()); err != nil {
		return fmt.Errorf("subnet validation failed: %w", err)
	}

	fmt.Printf("✅ Subnet '%s' validation passed\n", subnetName)
	return nil
}

func (vm *NutanixValidationManager) RunValidateImage(cmd *cobra.Command, args []string) error {
	imageName := args[0]

	// Get logger from context
	logger, err := ctxhelpers.GetLogger(cmd.Context())
	if err != nil {
		return fmt.Errorf("failed to get logger: %w", err)
	}

	host, _ := cmd.Flags().GetString(flagNtnxHost)
	port, _ := cmd.Flags().GetInt(flagNtnxPort)
	username, _ := cmd.Flags().GetString(flagNtnxUsername)
	password, _ := cmd.Flags().GetString(flagNtnxPassword)
	passwordKeyID, _ := cmd.Flags().GetInt(flagNtnxPasswordKeyID)
	insecure, _ := cmd.Flags().GetBool(flagNtnxInsecure)

	// Handle PWDB password retrieval if password key ID is provided
	if passwordKeyID > 0 {
		// TODO: Implement PWDB password retrieval
		// For now, we'll use the password flag
		if password == "" {
			return fmt.Errorf("password is required when not using PWDB")
		}
	}

	config := ConnectionConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Insecure: insecure,
	}

	manager := NewNutanixValidationManager(config, logger)

	if err := manager.Setup(cmd.Context()); err != nil {
		return fmt.Errorf("failed to setup validation manager: %w", err)
	}

	manager.AddImageValidator(imageName)

	fmt.Printf("Validating image: %s\n", imageName)
	if err := manager.Validate(cmd.Context()); err != nil {
		return fmt.Errorf("image validation failed: %w", err)
	}

	fmt.Printf("✅ Image '%s' validation passed\n", imageName)
	return nil
}

func (vm *NutanixValidationManager) RunValidateProject(cmd *cobra.Command, args []string) error {
	projectName := args[0]

	logger, err := ctxhelpers.GetLogger(cmd.Context())
	if err != nil {
		return fmt.Errorf("failed to get logger: %w", err)
	}

	host, _ := cmd.Flags().GetString(flagNtnxHost)
	port, _ := cmd.Flags().GetInt(flagNtnxPort)
	username, _ := cmd.Flags().GetString(flagNtnxUsername)
	password, _ := cmd.Flags().GetString(flagNtnxPassword)
	passwordKeyID, _ := cmd.Flags().GetInt(flagNtnxPasswordKeyID)
	insecure, _ := cmd.Flags().GetBool(flagNtnxInsecure)

	// Handle PWDB password retrieval if password key ID is provided
	if passwordKeyID > 0 {
		// TODO: Implement PWDB password retrieval
		// For now, we'll use the password flag
		if password == "" {
			return fmt.Errorf("password is required when not using PWDB")
		}
	}

	config := ConnectionConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Insecure: insecure,
	}

	manager := NewNutanixValidationManager(config, logger)

	if err := manager.Setup(cmd.Context()); err != nil {
		return fmt.Errorf("failed to setup validation manager: %w", err)
	}

	manager.AddProjectValidator(projectName)

	fmt.Printf("Validating project: %s\n", projectName)
	if err := manager.Validate(cmd.Context()); err != nil {
		return fmt.Errorf("project validation failed: %w", err)
	}

	fmt.Printf("✅ Project '%s' validation passed\n", projectName)
	return nil
}

func (vm *NutanixValidationManager) RunValidateCluster(cmd *cobra.Command, args []string) error {
	clusterName := args[0]

	logger, err := ctxhelpers.GetLogger(cmd.Context())
	if err != nil {
		return fmt.Errorf("failed to get logger: %w", err)
	}

	host, _ := cmd.Flags().GetString(flagNtnxHost)
	port, _ := cmd.Flags().GetInt(flagNtnxPort)
	username, _ := cmd.Flags().GetString(flagNtnxUsername)
	password, _ := cmd.Flags().GetString(flagNtnxPassword)
	passwordKeyID, _ := cmd.Flags().GetInt(flagNtnxPasswordKeyID)
	insecure, _ := cmd.Flags().GetBool(flagNtnxInsecure)

	// Handle PWDB password retrieval if password key ID is provided
	if passwordKeyID > 0 {
		// TODO: Implement PWDB password retrieval
		// For now, we'll use the password flag
		if password == "" {
			return fmt.Errorf("password is required when not using PWDB")
		}
	}

	config := ConnectionConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Insecure: insecure,
	}

	manager := NewNutanixValidationManager(config, logger)

	if err := manager.Setup(cmd.Context()); err != nil {
		return fmt.Errorf("failed to setup validation manager: %w", err)
	}

	manager.AddClusterValidator(clusterName)

	fmt.Printf("Validating cluster: %s\n", clusterName)
	if err := manager.Validate(cmd.Context()); err != nil {
		return fmt.Errorf("project validation failed: %w", err)
	}

	fmt.Printf("✅ Cluster '%s' validation passed\n", clusterName)
	return nil
}
