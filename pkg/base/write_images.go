package base

import (
	"io"

	"github.com/pkg/errors"
	kotsv1beta1 "github.com/replicatedhq/kots/kotskinds/apis/kots/v1beta1"
	"github.com/replicatedhq/kots/pkg/docker/registry"
	"github.com/replicatedhq/kots/pkg/image"
	"github.com/replicatedhq/kots/pkg/logger"
	kustomizeimage "sigs.k8s.io/kustomize/api/types"
)

type WriteUpstreamImageOptions struct {
	BaseDir        string
	AppSlug        string
	SourceRegistry registry.RegistryOptions
	DestRegistry   registry.RegistryOptions
	DryRun         bool
	IsAirgap       bool
	Log            *logger.Logger
	ReportWriter   io.Writer
	Installation   *kotsv1beta1.Installation
}

type WriteUpstreamImageResult struct {
	Images        []kustomizeimage.Image          // images to be rewritten
	CheckedImages []kotsv1beta1.InstallationImage // all images found in the installation
}

func CopyUpstreamImages(options WriteUpstreamImageOptions) (*WriteUpstreamImageResult, error) {
	checkedImages := makeImageInfoMap(options.Installation.Spec.KnownImages)
	newImages, err := image.CopyImages(options.SourceRegistry, options.DestRegistry, options.AppSlug, options.Log, options.ReportWriter, options.BaseDir, options.DryRun, options.IsAirgap, checkedImages)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save images")
	}

	return &WriteUpstreamImageResult{
		Images:        newImages,
		CheckedImages: makeInstallationImages(checkedImages),
	}, nil
}
