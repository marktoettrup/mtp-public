package airgap

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"maps"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/pkg/api/iksctl/v1alpha1"
	"bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/pkg/cmd/flags"
	ctxhelpers "bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/pkg/context"
	"bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/pkg/pwdb"

	"github.com/go-logr/logr"
	"github.com/regclient/regclient"
	"github.com/regclient/regclient/config"
	"github.com/regclient/regclient/types/ref"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

const (
	flagAirgapConfigFile = "config"
	flagAirgapCharts     = "charts"
	flagAirgapImages     = "images"
	flagAirgapDryRun     = "dry-run"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "airgap",
		RunE: runAirgap,
	}

	cmd.Flags().StringP(flagAirgapConfigFile, "c", "", "config file")
	if err := cmd.MarkFlagRequired(flagAirgapConfigFile); err != nil {
		panic(err)
	}
	cmd.Flags().Bool(flagAirgapCharts, true, "wether to airgap charts")
	cmd.Flags().Bool(flagAirgapImages, true, "wether to airgap images")
	cmd.Flags().Bool(flagAirgapDryRun, false, "enable dry-run")

	return cmd
}

func runAirgap(cmd *cobra.Command, args []string) error {
	log, err := ctxhelpers.GetLogger(cmd.Context())
	if err != nil {
		return err
	}

	configFile, err := cmd.Flags().GetString(flagAirgapConfigFile)
	if err != nil {
		return fmt.Errorf("could not read config flag: %w", err)
	}
	log.Info("Using config file", "file", configFile)

	bs, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("could not read config file")
	}

	scheme := runtime.NewScheme()
	if err := v1alpha1.AddToScheme(scheme); err != nil {
		return fmt.Errorf("could not add iksctl api to scheme")
	}
	decoder := serializer.NewCodecFactory(scheme).UniversalDeserializer()

	obj, _, err := decoder.Decode(bs, nil, &v1alpha1.AirgapConfig{})
	if err != nil {
		return fmt.Errorf("could not decode config: %w", err)
	}
	cfg := obj.(*v1alpha1.AirgapConfig)

	ctx, err := ctxhelpers.WithTmpDir(cmd.Context())
	if err != nil {
		return err
	}
	defer ctxhelpers.RmTmpDir(ctx)

	airgapCharts, err := cmd.Flags().GetBool(flagAirgapCharts)
	if err != nil {
		return fmt.Errorf("could not get airgap charts flag: %w", err)
	}

	dryRun, err := cmd.Flags().GetBool(flagAirgapDryRun)
	if err != nil {
		return fmt.Errorf("could not get airgap dry run flag: %w", err)
	}
	if dryRun {
		log.Info("Dry run is enabled. Printing actions without doing them.")
	}

	log.Info("Fetching oci credentials from PWDB")
	pwdbAPIKey, err := flags.GetPWDBAPIKey(cmd)
	if err != nil {
		return err
	}
	pwdbEndpoint, err := flags.GetPWDBEndpoint(cmd)
	if err != nil {
		return err
	}
	pwdbClient := pwdb.New(
		pwdb.WithEndpoint(pwdbEndpoint),
		pwdb.WithAPIKey(pwdbAPIKey),
	)

	ociRepoCreds, err := pwdbClient.GetPassword(cmd.Context(), cfg.OCIRepo.PWDBCredentialID)
	if err != nil {
		return err
	}

	for _, chart := range cfg.Charts {
		log := log.WithValues("chart", chart.Name, "version", chart.Version, "repo", chart.Repo)
		if err := pullChart(ctx, log, chart); err != nil {
			return err
		}
		if airgapCharts {
			err := pushChart(ctx, log, chart, cfg.OCIRepo, ociRepoCreds, dryRun)
			if err != nil {
				return fmt.Errorf("could not push chart %s to OCI repo: %w", chart.Name, err)
			}
		}

		rendered, err := templateChart(ctx, log, chart)
		if err != nil {
			return err
		}

		log.Info("extracting images from chart")
		images := extractImages(rendered)
		if err := airgapImages(ctx, log, images, chart, cfg.OCIRepo, ociRepoCreds, dryRun); err != nil {
			return err
		}
	}

	return nil
}

func pullChart(ctx context.Context, log logr.Logger, chart v1alpha1.Chart) error {
	tmpDir, err := ctxhelpers.GetTmpDir(ctx)
	if err != nil {
		return err
	}

	log.Info("Pulling chart")
	cmd := exec.CommandContext(
		ctx,
		"helm",
		"pull",
		chart.Name,
		"--repo",
		chart.Repo,
		"--version",
		chart.Version,
		"--destination",
		tmpDir,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not pull chart %s from %s: %w",
			chart.Name, chart.Repo, err)
	}

	return nil
}

func pushChart(
	ctx context.Context,
	log logr.Logger,
	chart v1alpha1.Chart,
	repo v1alpha1.OCIRepo,
	creds *pwdb.PasswordResponse,
	dryRun bool,
) error {
	tmpDir, err := ctxhelpers.GetTmpDir(ctx)
	if err != nil {
		return err
	}

	chartPath := filepath.Join(tmpDir, fmt.Sprintf("%s-%s.tgz", chart.Name, chart.Version))
	log = log.WithValues("tarball", chartPath, "targetRepo", repo.URL())
	if dryRun {
		log.Info("helm push")
		return nil
	}

	log.Info("Pushing helm chart")
	cmd := exec.CommandContext(
		ctx,
		"helm",
		"push",
		chartPath,
		"--username",
		creds.Username,
		"--password",
		creds.Password,
		repo.URL(),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not push chart %s to %s: %w", chart, repo.URL(), err)
	}
	return nil
}

func templateChart(ctx context.Context, log logr.Logger, chart v1alpha1.Chart) ([]byte, error) {
	log.Info("templating chart in order to look for images")

	tmpDir, err := ctxhelpers.GetTmpDir(ctx)
	if err != nil {
		return nil, err
	}

	args := []string{
		"template",
		filepath.Join(tmpDir, chart.TarBallName()),
	}
	if chart.Values != nil {
		args = append(args, "-f", "-")
	}

	var b bytes.Buffer
	c := exec.CommandContext(ctx, "helm", args...)
	c.Stdout = &b

	var w *io.PipeWriter
	if chart.Values != nil {
		var r *io.PipeReader
		r, w = io.Pipe()
		defer r.Close()
		defer w.Close()
		c.Stdin = r
	}

	if err := c.Start(); err != nil {
		return nil, fmt.Errorf("could not start helm template command for chart %s: %w", chart.Name, err)
	}

	if chart.Values != nil {
		if _, err := w.Write(chart.Values.Raw); err != nil {
			return nil, fmt.Errorf("could not write to command pipe: %w", err)
		}
		w.Close()
	}

	if err := c.Wait(); err != nil {
		return nil, fmt.Errorf("could not run helm template command for chart %s: %w", chart.Name, err)
	}

	return b.Bytes(), err
}

var re *regexp.Regexp = regexp.MustCompile(`\s*image:\s*(.*)`)

func extractImages(data []byte) []string {
	matches := re.FindAllSubmatch(data, -1)
	result := map[string]struct{}{}
	for _, match := range matches {
		result[string(match[1])] = struct{}{}
	}
	keys := slices.Collect(maps.Keys(result))
	trimmed := make([]string, len(keys))
	for i, key := range keys {
		trimmed[i] = strings.Trim(key, `"`)
	}
	return trimmed
}

func airgapImages(
	ctx context.Context,
	log logr.Logger,
	images []string,
	chart v1alpha1.Chart,
	repo v1alpha1.OCIRepo,
	creds *pwdb.PasswordResponse,
	dryRun bool,
) error {
	h := config.HostNewName(repo.Host)
	h.User = creds.Username
	h.Pass = creds.Password
	rc := regclient.New(regclient.WithConfigHost(*h))

	images = append(images, chart.ExtraImages...)

	for _, img := range images {
		refSrc, err := ref.New(img)
		if err != nil {
			return fmt.Errorf("could not parse ref from image %s: %w", img, err)
		}

		refTgt := refSrc
		refTgt.Registry = repo.Host
		refTgt.Repository = path.Join(repo.Repository, chart.Name, path.Base(refSrc.Repository))
		refTgt.Digest = ""

		log.Info("copying image", "srcRef", refSrc.CommonName(), "tgtRef", refTgt.CommonName())
		if dryRun {
			continue
		}

		if err := rc.ImageCopy(ctx, refSrc, refTgt); err != nil {
			return fmt.Errorf("could not copy image: %w", err)
		}
	}

	return nil
}
