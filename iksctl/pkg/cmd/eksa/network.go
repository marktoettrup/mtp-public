package eksa

import (
	"fmt"
	"net"
	"strings"

	parser "bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/internal/pkg/parsers/eksanywhere"
)

var privateNetworks = []string{
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
}

func (v *ClusterConfigValidator) validateNetworkConfiguration(network *parser.ClusterNetworkInfo) error {
	if len(network.Pods.CIDRBlocks) == 0 {
		return fmt.Errorf("pods cidrBlocks are required")
	}
	if len(network.Pods.CIDRBlocks) != 1 {
		return fmt.Errorf("exactly 1 pod CIDR block is required, got %d", len(network.Pods.CIDRBlocks))
	}

	podsCIDR := network.Pods.CIDRBlocks[0]
	if podsCIDR != "10.128.0.0/18" {
		return fmt.Errorf("pods CIDR must be '10.128.0.0/18', got '%s'", podsCIDR)
	}

	if len(network.Services.CIDRBlocks) == 0 {
		return fmt.Errorf("services cidrBlocks are required")
	}
	if len(network.Services.CIDRBlocks) != 1 {
		return fmt.Errorf("exactly 1 service CIDR block is required, got %d", len(network.Services.CIDRBlocks))
	}

	servicesCIDR := network.Services.CIDRBlocks[0]
	if servicesCIDR != "10.128.32.0/18" {
		return fmt.Errorf("services CIDR must be '10.128.32.0/18', got '%s'", servicesCIDR)
	}

	if err := v.validateBasicCIDR(podsCIDR, "pods CIDR"); err != nil {
		return err
	}
	if err := v.validateBasicCIDR(servicesCIDR, "services CIDR"); err != nil {
		return err
	}

	return nil
}

func (v *ClusterConfigValidator) validateControlPlaneEndpointAvailable() error {
	if v.clusterConfig.Endpoint == "" {
		return fmt.Errorf("cluster controlPlaneConfiguration.Endpoint.Host is required")
	}

	parsedIP := net.ParseIP(v.clusterConfig.Endpoint)
	if parsedIP == nil {
		return fmt.Errorf("cluster controlPlaneConfiguration.Endpoint.Host must be a valid IP address, got: %s", v.clusterConfig.Endpoint)
	}

	if err := v.validateEndpointIP(parsedIP); err != nil {
		return fmt.Errorf("cluster controlPlaneConfiguration.Endpoint.Host IP validation failed: %w", err)
	}

	return nil
}

func (v *ClusterConfigValidator) validateEndpointIP(ip net.IP) error {
	ipv4 := ip.To4()
	if ipv4 == nil {
		return fmt.Errorf("endpoint must be an IPv4 address, got: %s", ip.String())
	}

	if ip.IsUnspecified() {
		return fmt.Errorf("endpoint cannot be unspecified address (0.0.0.0)")
	}

	if ip.IsLoopback() {
		return fmt.Errorf("endpoint cannot be loopback address (%s)", ip.String())
	}

	if ip.IsLinkLocalUnicast() {
		return fmt.Errorf("endpoint cannot be link-local address (%s)", ip.String())
	}

	if ip.IsMulticast() {
		return fmt.Errorf("endpoint cannot be multicast address (%s)", ip.String())
	}

	if !v.isPrivateIPv4(ipv4) {
		return fmt.Errorf("endpoint IP must be in a private network range (10.0.0.0/8, 172.16.0.0/12, or 192.168.0.0/16), got %s", ip.String())
	}

	if ipv4[3] == 1 || ipv4[3] == 254 || ipv4[3] == 255 {
		return fmt.Errorf("endpoint IP cannot be a gateway address (last octet .1 or .254, .255), got %s", ip.String())
	}

	if ipv4[3] < 221 {
		return fmt.Errorf("endpoint IP must have last octet of 221 or higher, got %d in %s", ipv4[3], ip.String())
	}

	return nil
}

func (v *ClusterConfigValidator) validateIPv4Network(cidr, fieldName string) error {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return fmt.Errorf("%s is not a valid CIDR: %w", fieldName, err)
	}

	if ipNet.IP.To4() == nil {
		return fmt.Errorf("%s '%s' is not an IPv4 network", fieldName, cidr)
	}

	return nil
}

func (v *ClusterConfigValidator) validateIPv6Network(cidr, fieldName string) error {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return fmt.Errorf("%s is not a valid CIDR: %w", fieldName, err)
	}

	if ipNet.IP.To4() != nil {
		return fmt.Errorf("%s '%s' is not an IPv6 network", fieldName, cidr)
	}

	if ipNet.IP.To16() == nil {
		return fmt.Errorf("%s '%s' is not a valid IPv6 network", fieldName, cidr)
	}

	return nil
}

func (v *ClusterConfigValidator) addCniPluginValidations(cluster *parser.ClusterInfo) error {
	if err := v.validateNetworkConfiguration(&cluster.ClusterNetwork); err != nil {
		return fmt.Errorf("network configuration validation failed: %w", err)
	}

	if err := v.validateCiliumConfig(&cluster.ClusterNetwork.CNIConfig.Cilium); err != nil {
		return fmt.Errorf("Cilium configuration validation failed: %w", err)
	}

	return nil
}

func (v *ClusterConfigValidator) validateBasicCIDR(cidr, fieldName string) error {
	_, _, err := net.ParseCIDR(cidr)
	if err != nil {
		return fmt.Errorf("%s '%s' is not a valid CIDR: %w", fieldName, cidr, err)
	}
	return nil
}

func (v *ClusterConfigValidator) validateCiliumConfig(cilium *parser.CiliumConfigInfo) error {
	if cilium.PolicyEnforcementMode != "" {
		validModes := []string{"default", "always", "never"}
		if !v.isValidOption(cilium.PolicyEnforcementMode, validModes) {
			return fmt.Errorf("invalid policyEnforcementMode '%s', must be one of: %s",
				cilium.PolicyEnforcementMode, strings.Join(validModes, ", "))
		}
	}

	if cilium.RoutingMode != "" {
		validModes := []string{"default", "direct"}
		if !v.isValidOption(cilium.RoutingMode, validModes) {
			return fmt.Errorf("invalid routingMode '%s', must be one of: %s",
				cilium.RoutingMode, strings.Join(validModes, ", "))
		}
	}

	if cilium.IPv4NativeRoutingCIDR != "" {
		if err := v.validateBasicCIDR(cilium.IPv4NativeRoutingCIDR, "ipv4NativeRoutingCIDR"); err != nil {
			return err
		}

		if err := v.validateIPv4Network(cilium.IPv4NativeRoutingCIDR, "ipv4NativeRoutingCIDR"); err != nil {
			return err
		}

		if cilium.RoutingMode != "direct" {
			return fmt.Errorf("ipv4NativeRoutingCIDR can only be used when routingMode is set to 'direct'")
		}
	}

	if cilium.IPv6NativeRoutingCIDR != "" {
		if err := v.validateBasicCIDR(cilium.IPv6NativeRoutingCIDR, "ipv6NativeRoutingCIDR"); err != nil {
			return err
		}

		if err := v.validateIPv6Network(cilium.IPv6NativeRoutingCIDR, "ipv6NativeRoutingCIDR"); err != nil {
			return err
		}

		if cilium.RoutingMode != "direct" {
			return fmt.Errorf("ipv6NativeRoutingCIDR can only be used when routingMode is set to 'direct'")
		}
	}

	return nil
}

func (v *ClusterConfigValidator) isValidOption(value string, validOptions []string) bool {
	for _, option := range validOptions {
		if value == option {
			return true
		}
	}
	return false
}
