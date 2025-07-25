package ntnx

import (
	"context"
	"fmt"

	"bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/internal/pkg/validation"
	"bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/pkg/clients/ntnxv3Client/ntnxv3ClientProjects"
	"github.com/go-logr/logr"
	clusterClient "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/client"
	networkingClient "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/client"
	imagesClient "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/client"
)

func NewNutanixValidationManager(config ConnectionConfig, logger *logr.Logger) *NutanixValidationManager {
	return &NutanixValidationManager{
		connectionConfig:  config,
		validationManager: validation.NewValidationManager(logger),
	}
}

func (nvm *NutanixValidationManager) Setup(ctx context.Context) error {
	nvm.networkingClient = networkingClient.NewApiClient()
	nvm.networkingClient.Host = nvm.connectionConfig.Host
	nvm.networkingClient.Port = nvm.connectionConfig.Port
	nvm.networkingClient.Username = nvm.connectionConfig.Username
	nvm.networkingClient.Password = nvm.connectionConfig.Password
	nvm.networkingClient.VerifySSL = !nvm.connectionConfig.Insecure

	nvm.imagesClient = imagesClient.NewApiClient()
	nvm.imagesClient.Host = nvm.connectionConfig.Host
	nvm.imagesClient.Port = nvm.connectionConfig.Port
	nvm.imagesClient.Username = nvm.connectionConfig.Username
	nvm.imagesClient.Password = nvm.connectionConfig.Password
	nvm.imagesClient.VerifySSL = !nvm.connectionConfig.Insecure

	nvm.clusterClient = clusterClient.NewApiClient()
	nvm.clusterClient.Host = nvm.connectionConfig.Host
	nvm.clusterClient.Port = nvm.connectionConfig.Port
	nvm.clusterClient.Username = nvm.connectionConfig.Username
	nvm.clusterClient.Password = nvm.connectionConfig.Password
	nvm.clusterClient.VerifySSL = !nvm.connectionConfig.Insecure

	protocol := "https"
	if nvm.connectionConfig.Insecure {
		protocol = "http"
	}

	projectsClient := ntnxv3ClientProjects.NewClientWithBasicAuth(
		fmt.Sprintf("%s:%d", nvm.connectionConfig.Host, nvm.connectionConfig.Port),
		"/api/nutanix/v3",
		protocol,
		nvm.connectionConfig.Username,
		nvm.connectionConfig.Password,
	)
	nvm.projectsClient = &projectsClient

	return nil
}

func (nvm *NutanixValidationManager) AddImageValidator(imageName string) *NutanixValidationManager {
	imageValidator := NewImageValidator(imageName, nvm.imagesClient)
	nvm.validationManager.AddValidator(imageValidator)
	return nvm
}

func (nvm *NutanixValidationManager) AddSubnetValidator(subnetName string) *NutanixValidationManager {
	subnetValidator := NewSubnetValidator(subnetName, nvm.networkingClient)
	nvm.validationManager.AddValidator(subnetValidator)
	return nvm
}

func (nvm *NutanixValidationManager) AddProjectValidator(projectName string) *NutanixValidationManager {
	projectValidator := NewProjectValidator(projectName, *nvm.projectsClient)
	nvm.validationManager.AddValidator(projectValidator)
	return nvm
}

func (nvm *NutanixValidationManager) AddClusterValidator(clusterName string) *NutanixValidationManager {
	clusterValidator := NewclusterValidator(clusterName, nvm.clusterClient)
	nvm.validationManager.AddValidator(clusterValidator)
	return nvm
}

func (nvm *NutanixValidationManager) Validate(ctx context.Context) error {
	return nvm.validationManager.ValidateAll(ctx)
}
