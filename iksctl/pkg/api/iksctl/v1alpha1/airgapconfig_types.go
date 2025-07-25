package v1alpha1

import (
	"fmt"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
type AirgapConfig struct {
	metav1.TypeMeta `json:",inline"`

	Charts  []Chart `json:"charts,omitempty"`
	OCIRepo OCIRepo `json:"ociRepo"`
}

type Chart struct {
	Repo        string                `json:"repo,omitempty"`
	Name        string                `json:"name,omitempty"`
	Version     string                `json:"version,omitempty"`
	Values      *apiextensionsv1.JSON `json:"values,omitempty"`
	ExtraImages []string              `json:"extraImages,omitempty"`
}

func (c Chart) TarBallName() string {
	return fmt.Sprintf("%s-%s.tgz", c.Name, c.Version)
}

type OCIRepo struct {
	Host             string `json:"host,omitempty"`
	PWDBCredentialID int    `json:"pwdbCredentialID"`
	Repository       string `json:"repository"`
}

func (r OCIRepo) URL() string {
	return fmt.Sprintf("oci://%s/%s", r.Host, r.Repository)
}
