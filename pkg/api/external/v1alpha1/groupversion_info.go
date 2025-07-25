// Package v1alpha1 contains external API Schema definitions for EKS Anywhere types
// +kubebuilder:object:generate=false
package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	ClusterGroupVersion = schema.GroupVersion{
		Group:   "anywhere.eks.amazonaws.com",
		Version: "v1alpha1",
	}

	NutanixMachineConfigGroupVersion = schema.GroupVersion{
		Group:   "anywhere.eks.amazonaws.com",
		Version: "v1alpha1",
	}

	NutanixDatacenterConfigGroupVersion = schema.GroupVersion{
		Group:   "anywhere.eks.amazonaws.com",
		Version: "v1alpha1",
	}

	FluxConfigGroupVersion = schema.GroupVersion{
		Group:   "anywhere.eks.amazonaws.com",
		Version: "v1alpha1",
	}
)

var (
	ClusterGVK = ClusterGroupVersion.WithKind("Cluster")

	NutanixMachineConfigGVK = NutanixMachineConfigGroupVersion.WithKind("NutanixMachineConfig")

	NutanixDatacenterConfigGVK = NutanixDatacenterConfigGroupVersion.WithKind("NutanixDatacenterConfig")

	FluxConfigGVK = FluxConfigGroupVersion.WithKind("FluxConfig")
)

func GetSupportedGVKs() map[schema.GroupVersionKind]bool {
	return map[schema.GroupVersionKind]bool{
		ClusterGVK:                 true,
		NutanixMachineConfigGVK:    true,
		NutanixDatacenterConfigGVK: true,
		FluxConfigGVK:              true,
	}
}
