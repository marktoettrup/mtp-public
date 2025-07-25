package parser

import (
	"fmt"
	"io"
	"os"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"

	externaltypes "bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/pkg/api/external/v1alpha1"
)

type Parser struct {
	decoder runtime.Decoder
}

func New() *Parser {
	decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

	return &Parser{
		decoder: decoder,
	}
}

func (p *Parser) ParseFile(filePath string) ([]*unstructured.Unstructured, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	return p.ParseReader(file)
}

func (p *Parser) ParseReader(reader io.Reader) ([]*unstructured.Unstructured, error) {
	var objects []*unstructured.Unstructured

	decoder := utilyaml.NewYAMLToJSONDecoder(reader)

	for {
		var obj unstructured.Unstructured
		err := decoder.Decode(&obj)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Warning: failed to decode YAML object: %v\n", err)
			continue
		}

		if len(obj.Object) == 0 {
			continue
		}

		objects = append(objects, &obj)
	}

	return objects, nil
}

func (p *Parser) GetClusters(objects []*unstructured.Unstructured) []*unstructured.Unstructured {
	var clusters []*unstructured.Unstructured

	for _, obj := range objects {
		if obj.GroupVersionKind() == externaltypes.ClusterGVK {
			clusters = append(clusters, obj)
		}
	}

	return clusters
}

func (p *Parser) GetNutanixMachineConfigs(objects []*unstructured.Unstructured) []*unstructured.Unstructured {
	var configs []*unstructured.Unstructured

	for _, obj := range objects {
		if obj.GroupVersionKind() == externaltypes.NutanixMachineConfigGVK {
			configs = append(configs, obj)
		}
	}

	return configs
}

func (p *Parser) GetNutanixDatacenterConfigs(objects []*unstructured.Unstructured) []*unstructured.Unstructured {
	var configs []*unstructured.Unstructured

	for _, obj := range objects {
		if obj.GroupVersionKind() == externaltypes.NutanixDatacenterConfigGVK {
			configs = append(configs, obj)
		}
	}

	return configs
}

func (p *Parser) GetFluxConfigs(objects []*unstructured.Unstructured) []*unstructured.Unstructured {
	var configs []*unstructured.Unstructured

	for _, obj := range objects {
		if obj.GroupVersionKind() == externaltypes.FluxConfigGVK {
			configs = append(configs, obj)
		}
	}

	return configs
}

func (p *Parser) ValidateSupportedCRDs(objects []*unstructured.Unstructured) error {
	supportedGVKs := externaltypes.GetSupportedGVKs()

	for _, obj := range objects {
		gvk := obj.GroupVersionKind()
		if !supportedGVKs[gvk] {
			return fmt.Errorf("unsupported CRD type detected: %s/%s %s", gvk.Group, gvk.Version, gvk.Kind)
		}
	}

	return nil
}

func (p *Parser) ParseStructure(filePath string, logger *logr.Logger) ([]*ClusterInfo, []*NutanixMachineConfigInfo, error) {
	objects, err := p.ParseFile(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse file: %w", err)
	}

	if len(objects) == 0 {
		return nil, nil, fmt.Errorf("no objects found in file %s", filePath)
	}

	// Validate that only supported CRD types are present
	if err := p.ValidateSupportedCRDs(objects); err != nil {
		return nil, nil, fmt.Errorf("unsupported CRD types detected: %w", err)
	}

	var clusterInfos []*ClusterInfo
	clusters := p.GetClusters(objects)
	for _, cluster := range clusters {
		info, err := p.ExtractClusterInfo(cluster)
		if err != nil {
			logger.Error(err, "Failed to extract cluster info", "cluster", cluster.GetName())
			continue
		}

		clusterInfos = append(clusterInfos, &info)
	}

	var nutanixMachineConfigs []*NutanixMachineConfigInfo
	configs := p.GetNutanixMachineConfigs(objects)
	for _, config := range configs {
		info, err := p.ExtractNutanixMachineConfigInfo(config)
		if err != nil {
			logger.Error(err, "Failed to extract machine config info", "config", config.GetName())
			continue
		}

		nutanixMachineConfigs = append(nutanixMachineConfigs, &info)
	}

	return clusterInfos, nutanixMachineConfigs, nil
}
