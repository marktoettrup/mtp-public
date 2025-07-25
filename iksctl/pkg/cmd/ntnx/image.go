package ntnx

import (
	"context"
	"fmt"
	"strings"

	utils "bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/internal/pkg/utils"

	imagesApi "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/api"
	imagesClient "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/client"
	content "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/content"
)

func NewImageValidator(imageName string, client *imagesClient.ApiClient) *ImageValidator {
	return &ImageValidator{
		imageName: imageName,
		api:       imagesApi.NewImagesApi(client),
	}
}

func (v *ImageValidator) Name() string {
	return fmt.Sprintf("image/%s", v.imageName)
}

func (v *ImageValidator) Validate(ctx context.Context) error {
	cleanImageName := strings.TrimSpace(strings.ToLower(v.imageName))
	resource, err := v.GetResource(ctx)
	if err != nil {
		return err
	}
	images := resource.([]content.Image)

	if len(images) == 0 {
		return fmt.Errorf("image %s not found", v.imageName)
	}

	for _, image := range images {
		if strings.ToLower(strings.TrimSpace(*image.Name)) == cleanImageName {
			return nil
		}
	}

	return fmt.Errorf("image %s not found", v.imageName)
}

func (v *ImageValidator) GetResource(ctx context.Context) (interface{}, error) {
	cleanImageName := strings.TrimSpace(strings.ToLower(v.imageName))
	page, limit, filter := utils.GetDefaultFilter(cleanImageName)
	resp, err := v.api.ListImages(&page, &limit, &filter, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	images, ok := resp.Data.GetValue().([]content.Image)
	if !ok {
		return nil, fmt.Errorf("failed to parse image data")
	}

	return images, nil
}
