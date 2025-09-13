package helmquilt

import (
	"context"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"golang.org/x/mod/semver"

	"github.com/luisdavim/helmquilt/pkg/config"
	"github.com/luisdavim/helmquilt/pkg/logger"
	"github.com/luisdavim/helmquilt/pkg/utils"
)

// fetchChart gets the helm chart from the given source, storing it in the given path
func fetchChart(ctx context.Context, source config.Source, chartDownloadDir string) (string, error) {
	logger := logger.FromContext(ctx)
	logger.Printf("Fetching chart %s from %s\n", filepath.Join(source.ChartPath, source.ChartName), source.URL)

	chartDownloadPath := filepath.Join(chartDownloadDir, source.ChartName)
	if _, err := os.Stat(chartDownloadPath); err == nil {
		if err := os.RemoveAll(chartDownloadPath); err != nil {
			return "", err
		}
	}

	// the source is a git repository
	if strings.HasPrefix(source.URL, "git@") || strings.HasSuffix(source.URL, ".git") {
		ref := plumbing.ReferenceName(source.Version)
		if ref != "" && ref.Validate() != nil {
			if semver.IsValid(source.Version) {
				ref = plumbing.NewTagReferenceName(source.Version)
			} else {
				ref = plumbing.NewBranchReferenceName(source.Version)
			}
		}

		_, err := git.PlainClone(chartDownloadPath, false, &git.CloneOptions{
			URL:           source.URL,
			Depth:         1,
			SingleBranch:  true,
			ReferenceName: ref,
		})
		if err != nil {
			return "", err
		}

		logger.Println("Cloned chart to", chartDownloadPath)
		return chartDownloadPath, nil
	}

	// the source is an HTTP helm repo
	if strings.HasPrefix(source.URL, "http") {
		chartURL, err := url.JoinPath(source.URL, source.ChartPath)
		if err != nil {
			return "", err
		}
		chartDownloadPath, err = utils.DownloadChart(chartURL, source.ChartName, source.Version, chartDownloadDir)
		if err != nil {
			return "", err
		}
		logger.Println("Downloaded chart to", chartDownloadPath)
		return chartDownloadPath, nil
	}

	// the source is an OCI registry
	chartDownloadPath, err := utils.PullChart(source.URL, source.ChartName, source.Version, chartDownloadDir, false)
	if err != nil {
		return "", err
	}
	logger.Println("Pulled chart to", chartDownloadPath)
	return chartDownloadPath, nil
}
