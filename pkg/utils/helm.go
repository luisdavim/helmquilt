package utils

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"helm.sh/helm/v4/pkg/action"
	helmchart "helm.sh/helm/v4/pkg/chart/v2"
	"helm.sh/helm/v4/pkg/cli"
	"helm.sh/helm/v4/pkg/registry"
	helmrepo "helm.sh/helm/v4/pkg/repo/v1"
	"sigs.k8s.io/yaml"
)

const (
	chartMetaFile = "Chart.yaml"
)

func DownloadChart(repo, chart, version, dst string) (string, error) {
	data, err := HTTPGet(repo + "/index.yaml")
	if err != nil {
		var httpErr *ErrHTTP
		if errors.As(err, &httpErr) {
			if httpErr.Code == http.StatusNotFound {
				data, err = HTTPGet(repo + "/index.json")
			}
		}
		if err != nil {
			return "", fmt.Errorf("failed to download charts index: %w", err)
		}
	}

	var idx helmrepo.IndexFile
	if err := yaml.Unmarshal(data, &idx); err != nil {
		return "", fmt.Errorf("faield to read index: %w", err)
	}

	// lookup the chart version in the index
	for _, e := range idx.Entries[chart] {
		if e.Version == version {
			// TODO: on error try the next URL
			chartTar := filepath.Join(dst, filepath.Base(e.URLs[0]))
			if err := DownloadFile(chartTar, e.URLs[0]); err != nil {
				return "", fmt.Errorf("failed to download %s: %w", e.URLs[0], err)
			}
			chartPath := filepath.Join(dst, chart)
			if err := Extract(chartTar, chartPath); err != nil {
				return "", fmt.Errorf("failed to extract %s: %w", chartTar, err)
			}
			_ = os.Remove(chartTar)

			return chartPath, nil
		}
	}

	return "", fmt.Errorf("no matching version %q for %q", version, chart)
}

func PullChart(reg, chart, version, dst string, debug bool) (string, error) {
	if dst == "" {
		dst = "."
	}
	// registry client
	registryClient, err := registry.NewClient(
		registry.ClientOptDebug(debug),
		registry.ClientOptEnableCache(false),
	)
	if err != nil {
		return "", err
	}

	// init helm action config
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(nil, "", "secret"); err != nil {
		return "", err
	}

	actionConfig.RegistryClient = registryClient

	// pull the chart
	pull := action.NewPull(action.WithConfig(actionConfig))
	pull.Settings = cli.New() // didn't want to do this but otherwise it goes nil pointer
	pull.Version = version
	pull.DestDir = dst
	pull.UntarDir = chart
	pull.Untar = true
	if _, err := pull.Run(fmt.Sprintf("%s/%s", reg, chart)); err != nil {
		return "", err
	}

	return filepath.Join(dst, chart), nil
}

func GetChartVersion(chartDir string) (string, error) {
	chartMeta, err := os.ReadFile(filepath.Join(chartDir, chartMetaFile))
	if err != nil {
		return "", err
	}

	var cm helmchart.Metadata
	if err := yaml.Unmarshal(chartMeta, &cm); err != nil {
		return "", err
	}

	return cm.Version, nil
}

func UpdateChartMetadata(version, name string, chartDir string) error {
	if version == "" && name == "" {
		return nil
	}

	chartMeta, err := os.ReadFile(filepath.Join(chartDir, chartMetaFile))
	if err != nil {
		return err
	}

	var cm helmchart.Metadata
	if err := yaml.Unmarshal(chartMeta, &cm); err != nil {
		return err
	}

	if (version == "" || cm.Version == version) && (name == "" || cm.Name == name) {
		return nil
	}

	if version != "" {
		cm.Version = version
	}

	if name != "" {
		cm.Name = name
	}

	chartMeta, err = yaml.Marshal(cm)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(chartDir, chartMetaFile), chartMeta, os.ModePerm); err != nil {
		return err
	}

	return nil
}
