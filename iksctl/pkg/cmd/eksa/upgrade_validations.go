package eksa

import (
	"fmt"
	"strings"

	parser "bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/internal/pkg/parsers/eksanywhere"
	"k8s.io/apimachinery/pkg/util/version"
)

const supportedMinorVersionIncrement = 1

func (v *ClusterConfigValidator) addK8sVersionValidation(cluster *parser.ClusterInfo, machineConfig []*parser.NutanixMachineConfigInfo) error {
	if cluster.KubernetesVersion == "" {
		return fmt.Errorf("kubernetes version is required")
	}

	for _, machine := range machineConfig {
		if machine.ImageName == "" {
			return fmt.Errorf("image name is required for machine config %s", machine.Name)
		}
		if !isImageCompatibleWithK8sVersion(machine.ImageName, cluster.KubernetesVersion) {
			return fmt.Errorf("image %s is not compatible with Kubernetes version %s", machine.ImageName, cluster.KubernetesVersion)
		}
	}

	return nil
}

func isImageCompatibleWithK8sVersion(imageName, k8sVersion string) bool {
	expectedSuffix := strings.Replace(k8sVersion, ".", "-", -1)

	return strings.HasSuffix(imageName, expectedSuffix)
}

func (v *ClusterConfigValidator) addValidateKubeVersionSkew(new, old string) error {
	if new == "" || old == "" {
		return fmt.Errorf("kubernetes version is required for both new and old cluster configurations")
	}

	parsedOldVersion, err := version.ParseGeneric(old)
	if err != nil {
		return fmt.Errorf("parsing old kubernetes version %s: %v", old, err)
	}

	parsedNewVersion, err := version.ParseGeneric(new)
	if err != nil {
		return fmt.Errorf("parsing new kubernetes version %s: %v", new, err)
	}

	if parsedNewVersion.Minor() == parsedOldVersion.Minor() && parsedNewVersion.Major() == parsedOldVersion.Major() {
		return fmt.Errorf("kubernetes version skew validation failed: new version %s is the same as old version %s", parsedNewVersion, parsedOldVersion)
	}

	if err := validateVersionSkew(parsedOldVersion, parsedNewVersion); err != nil {
		return fmt.Errorf("kubernetes version skew validation failed: %w", err)
	}

	return nil
}

func (v *ClusterConfigValidator) addValidateEksaVersionSkew(new, old string) error {
	if new == "" || old == "" {
		return fmt.Errorf("eksa version is required for both new and old cluster configurations")
	}

	oldEksaVersion, err := version.ParseGeneric(old)
	if err != nil {
		return fmt.Errorf("parsing old eksa version %s: %v", old, err)
	}

	newEksaVersion, err := version.ParseGeneric(new)
	if err != nil {
		return fmt.Errorf("parsing new eksa version %s: %v", new, err)
	}

	if err := validateVersionSkew(oldEksaVersion, newEksaVersion); err != nil {
		return fmt.Errorf("eksa version skew validation failed: %w", err)
	}

	return nil
}

func validateVersionSkew(oldVersion, newVersion *version.Version) error {
	if newVersion.LessThan(oldVersion) {
		return fmt.Errorf("version downgrade is not supported (%s) -> (%s)", oldVersion, newVersion)
	}

	newVersionMinor := newVersion.Minor()
	oldVersionMinor := oldVersion.Minor()

	minorVersionDifference := int(newVersionMinor) - int(oldVersionMinor)

	if minorVersionDifference < 0 || minorVersionDifference > supportedMinorVersionIncrement {
		return fmt.Errorf("only +%d minor version skew is supported, detected skew: %d", supportedMinorVersionIncrement, minorVersionDifference)
	}

	return nil
}
