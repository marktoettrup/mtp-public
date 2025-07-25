# EKS Anywhere Nutanix Cluster Validation Reference

This document provides a comprehensive overview of all validations performed by the `iksctl` tool for EKS Anywhere clusters running on Nutanix infrastructure.

## Overview

The `iksctl` tool performs multiple layers of validation to ensure that EKS Anywhere cluster configurations are correct, secure, and deployable on Nutanix infrastructure. Validations are categorized into several areas:

## Validation Categories

### 1. Cluster Manifest Validation

#### Cluster Resource Structure

- **Single Cluster Requirement**: Exactly one `Cluster` resource must be present in the manifest
- **Control Plane Configuration**: Control plane configuration must be present and valid
- **Worker Node Groups**: Worker node group configurations must be properly defined

#### Resource Reference Validation

- **Control Plane Machine Reference**:
  - Must reference a valid `NutanixMachineConfig` resource
  - Referenced machine config must exist in the manifest
  - Machine config kind must be exactly `NutanixMachineConfig`
- **Worker Node Machine References**:
  - Each worker group must reference a valid `NutanixMachineConfig` resource
  - All referenced machine configs must exist in the manifest
  - Machine config kind must be exactly `NutanixMachineConfig`

#### Datacenter and GitOps References

- **Datacenter Reference**: If specified, must be of kind `NutanixDatacenterConfig`
- **GitOps Reference**: If specified, must be of kind `FluxConfig`

### 2. Machine Configuration Validation

#### CPU Configuration

- **Minimum Requirements**:
  - At least 2 vCPU sockets required
  - At least 1 vCPU per socket required
  - Total vCPUs calculated as `vcpusPerSocket × vcpuSockets`
- **Maximum Limits**:
  - vCPU sockets: Maximum 64
  - vCPUs per socket: Maximum 128  
  - Total vCPUs: Maximum 512
- **Consistency Check**: Total vCPUs must match the product of sockets and vCPUs per socket

#### Memory Configuration

- **Format Validation**: Memory size must be specified with valid units (Gi, Mi, Ki)
- **Required Field**: Memory size specification is mandatory
- **Parsing Validation**: Memory values must be parseable integers with supported units

#### Storage Configuration

- **Format Validation**: System disk size must be specified with valid units (Gi, Mi, Ki)
- **Required Field**: System disk size specification is mandatory  
- **Parsing Validation**: Storage values must be parseable integers with supported units

#### Resource References

All Nutanix resource references are mandatory:

- **Nutanix Cluster Reference**: Must specify target Nutanix cluster name
- **Image Reference**: Must specify valid VM image name
- **Subnet Reference**: Must specify target subnet name
- **Project Reference**: Must specify Nutanix project name

### 3. Network and CNI Validation

#### Cluster Network Configuration

- **Pods CIDR**: Must be exactly `10.128.0.0/18` with exactly one CIDR block
- **Services CIDR**: Must be exactly `10.128.32.0/18` with exactly one CIDR block
- **DNS Configuration**: ClusterIP must be `10.128.32.10`

#### CNI Configuration (Cilium)

- **Policy Enforcement**: Cilium policy enforcement mode validation
- **Native Routing**: Native routing CIDR validation when specified
- **CNI Provider**: Must be Cilium for supported configurations

#### Control Plane Endpoint Validation

- **IP Address Format**: Must be a valid IPv4 address
- **Network Restrictions**:
  - Cannot be unspecified (0.0.0.0)
  - Cannot be loopback address
  - Cannot be link-local address  
  - Cannot be multicast address
  - Cannot be in reserved IP ranges
- **Private Network Requirement**: Must be in RFC 1918 private address space
- **Gateway Restrictions**: Cannot be common gateway addresses (.1 or .254)
- **High Octet Requirement**: Last octet must be ≥ 221

### 4. Version Validation

#### Kubernetes Version Validation

- **Version Format**: Must be a valid semantic version (e.g., "1.30", "1.29")
- **Required Field**: Kubernetes version is mandatory for all clusters
- **Image Compatibility**: VM images must be compatible with the specified Kubernetes version
  - Image names must end with the Kubernetes version (with dots replaced by hyphens)
  - Example: Image `IKS_2024.0.1_EKSA_v0.20.11_BUILD_202505251644_KUBE_1-30` is compatible with Kubernetes version `1.30`

#### EKS-A Version Validation

- **Version Format**: Must be a valid semantic version (e.g., "v0.22.1", "v0.21.0")
- **Required Field**: EKS-A version is mandatory for all clusters

#### Version Upgrade Validation

- **Coordinated Upgrades**: Both Kubernetes and EKS-A version upgrades must be specified together
  - If `--upgrade-k8s-version-to` is specified, `--upgrade-eks-version-to` must also be specified
  - If `--upgrade-eks-version-to` is specified, `--upgrade-k8s-version-to` must also be specified
- **Version Skew Policy**:
  - **No Downgrades**: Version downgrades are not supported for either Kubernetes or EKS-A
  - **Minor Version Increment**: Only +1 minor version upgrades are supported
  - **Same Version Prevention**: Upgrading to the same version is not allowed
- **Parsing Validation**: Both old and new versions must be valid semantic versions

### 5. Nutanix Resource Validation

#### Infrastructure Resource Availability

- **Project Validation**: Validates that specified Nutanix projects exist and are accessible
- **Cluster Validation**: Verifies that target Nutanix clusters are available
- **Image Validation**: Confirms that specified VM images exist in the environment
- **Subnet Validation**: Validates that specified subnets are available and accessible

## Validation Command Usage

### Main Validation Command

```bash
# Validate complete cluster manifest
iksctl parse validate -f cluster-manifest.yaml \
  --host prism.example.com \
  --username admin \
  --password secret
```

### Cluster Upgrade Validation

```bash
# Validate cluster upgrade with version changes
iksctl parse validate -f cluster-manifest.yaml \
  --host prism.example.com \
  --username admin \
  --password secret \
  --upgrade-k8s-version-to "1.30" \
  --upgrade-eks-version-to "v0.22.1"
```

### Individual Resource Validation

```bash
# Validate specific Nutanix resources
iksctl ntnx validate-project "My Project"
iksctl ntnx validate-cluster "My Cluster" 
iksctl ntnx validate-image "My Image"
iksctl ntnx validate-subnet "My Subnet"
```
