package config

import (
	"fmt"
	"path/filepath"
)

type Chart struct {
	// Name of the Chart
	Name string `json:"name,omitempty"`
	// Version to set on the chart after applying changes, leave empty to keep the original version
	Version string `json:"version,omitempty"`
	// Path where to store the Chart
	Path string `json:"path,omitempty"`
	// Source holds information about where to pull the chart from and what version to pull
	Source Source `json:"source,omitempty"`
	// Patches is the set of patch files to apply to this chart
	Patches []string `json:"patches,omitempty"`
	// Repack indicates wether or not the chart should be stored as a tarball
	Repack         bool `json:"repack,omitempty"`
	FileOperations `json:",inline"`
}

type FileOperations struct {
	// Remove is a list of files or directories to delete from the chart
	Remove []string `json:"remove,omitempty"`
	// Keep defines what files or directories to keep fromm the original chart
	// when set, any path not included in this list will be removed
	Keep []string `json:"keep,omitempty"`
	// Move is a set of instructions for renaming chart content
	Move []Move `json:"move,omitempty"`
}

func (f *FileOperations) HasOperations() bool {
	return len(f.Remove) != 0 || len(f.Move) != 0 || len(f.Keep) != 0
}

type Source struct {
	// URL from where to get the chart from, this can be a git repo, OCI registry or a helm repo
	URL string `json:"url,omitempty"`
	// Version to pull from the repo or registry
	Version string `json:"version,omitempty"`
	// ChartName is the name of the chart in the repo or registry
	ChartName string `json:"chartName,omitempty"`
	// ChartPath is a sub-path where to find the chart in the repo or registry
	ChartPath string `json:"chartPath,omitempty"`
}

type Move struct {
	Source string `json:"source,omitempty"`
	Dest   string `json:"dest,omitempty"`
}

func (c *Chart) GetName() string {
	return filepath.Join(c.Path, c.Source.ChartName)
}

func (c *Chart) GetFullName(workDir string) string {
	return filepath.Join(workDir, c.GetName())
}

func (c *Chart) GetTarName() string {
	return fmt.Sprintf("%s-%s.tgz", c.GetName(), c.Version)
}
