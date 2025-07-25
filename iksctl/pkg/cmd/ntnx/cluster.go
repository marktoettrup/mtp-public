package ntnx

import (
	"context"
	"fmt"
	"strings"

	"bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/internal/pkg/utils"
	clustersApi "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/api"
	clusterClient "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/client"
	config "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
)

func NewclusterValidator(clusterName string, client *clusterClient.ApiClient) *ClusterValidator {
	return &ClusterValidator{
		clusterName: clusterName,
		api:         clustersApi.NewClustersApi(client),
	}
}

func (v *ClusterValidator) Name() string {
	return fmt.Sprintf("cluster/%s", v.clusterName)
}

func (v *ClusterValidator) Validate(ctx context.Context) error {
	cleanclusterName := strings.TrimSpace(strings.ToLower(v.clusterName))
	resource, err := v.GetResource(ctx)
	if err != nil {
		return err
	}

	clusters := resource.([]config.Cluster)

	if len(clusters) == 0 {
		return fmt.Errorf("cluster %s not found", v.clusterName)
	}

	for _, cluster := range clusters {
		if strings.ToLower(strings.TrimSpace(*cluster.Name)) == cleanclusterName {
			return nil
		}
	}

	return fmt.Errorf("cluster %s not found", v.clusterName)
}

func (v *ClusterValidator) GetResource(ctx context.Context) (interface{}, error) {
	cleanImageName := strings.TrimSpace(strings.ToLower(v.clusterName))
	page, limit, filter := utils.GetDefaultFilter(cleanImageName)
	resp, err := v.api.ListClusters(&page, &limit, &filter, nil, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list clusters: %w", err)
	}

	clusters, ok := resp.Data.GetValue().([]config.Cluster)
	if !ok {
		return nil, fmt.Errorf("failed to parse cluster data")
	}

	return clusters, nil
}
