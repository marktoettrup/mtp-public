package ntnx

import (
	"bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/internal/pkg/validation"
	"bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/pkg/clients/ntnxv3Client/ntnxv3ClientProjects"
	clustersApi "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/api"
	clusterClient "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/client"
	subnetsApi "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/api"
	networkingClient "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/client"
	imagesApi "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/api"
	imagesClient "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/client"
)

type ConnectionConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Insecure bool
}

type NutanixValidationManager struct {
	connectionConfig  ConnectionConfig
	validationManager *validation.ValidationManager
	networkingClient  *networkingClient.ApiClient
	imagesClient      *imagesClient.ApiClient
	projectsClient    *ntnxv3ClientProjects.ProjectsService
	clusterClient     *clusterClient.ApiClient
}

type ProjectValidator struct {
	projectName string
	api         ntnxv3ClientProjects.ProjectsService
}

type ClusterValidator struct {
	clusterName string
	api         *clustersApi.ClustersApi
}

type ImageValidator struct {
	imageName string
	api       *imagesApi.ImagesApi
}

type SubnetValidator struct {
	subnetName string
	api        *subnetsApi.SubnetsApi
}
