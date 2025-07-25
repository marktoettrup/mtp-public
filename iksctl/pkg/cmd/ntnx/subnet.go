package ntnx

import (
	"context"
	"fmt"
	"strings"

	utils "bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/internal/pkg/utils"
	subnetsApi "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/api"
	networkingClient "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/client"
	config "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
)

type Subnet struct {
	Name                   string `json:"name"`
	NetworkId              int    `json:"networkId"`
	VirtualSwitchReference string `json:"virtualSwitchReference"`
	ClusterReference       string `json:"clusterReference"`
	SubnetType             string `json:"subnetType"`
}

func NewSubnetValidator(subnetName string, client *networkingClient.ApiClient) *SubnetValidator {
	return &SubnetValidator{
		subnetName: subnetName,
		api:        subnetsApi.NewSubnetsApi(client),
	}
}

func (v *SubnetValidator) Name() string {
	return fmt.Sprintf("subnet/%s", v.subnetName)
}

func (v *SubnetValidator) Validate(ctx context.Context) error {
	cleanSubnetName := strings.TrimSpace(strings.ToLower(v.subnetName))
	resource, err := v.GetResource(ctx)
	if err != nil {
		return err
	}

	subnets, ok := resource.([]config.Subnet)
	if !ok {
		return fmt.Errorf("expected []config.Subnet, got %T", resource)
	}

	if len(subnets) == 0 {
		return fmt.Errorf("subnet %s not found", v.subnetName)
	}

	for _, subnet := range subnets {
		if strings.ToLower(strings.TrimSpace(*subnet.Name)) == cleanSubnetName {
			return nil
		}
	}

	return fmt.Errorf("image %s not found", v.subnetName)
}

func (v *SubnetValidator) GetResource(ctx context.Context) (interface{}, error) {
	cleanSubnetName := strings.TrimSpace(strings.ToLower(v.subnetName))
	page, limit, filter := utils.GetDefaultFilter(cleanSubnetName)
	resp, err := v.api.ListSubnets(&page, &limit, &filter, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	subnets, ok := resp.Data.GetValue().([]config.Subnet)
	if !ok {
		return nil, fmt.Errorf("failed to parse subnet data")
	}
	return subnets, nil
}
