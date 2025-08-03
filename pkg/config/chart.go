package config

import (
	"fmt"
	"path/filepath"
)

type Chart struct {
	Name           string   `json:"name,omitempty"`
	Version        string   `json:"version,omitempty"`
	Path           string   `json:"path,omitempty"`
	Source         Source   `json:"source,omitempty"`
	Patches        []string `json:"patches,omitempty"`
	FileOperations `json:",inline"`
}

type FileOperations struct {
	Remove []string `json:"remove,omitempty"`
	Keep   []string `json:"keep,omitempty"`
	Move   []Move   `json:"move,omitempty"`
}

func (f *FileOperations) HasOperations() bool {
	return len(f.Remove) != 0 || len(f.Move) != 0 || len(f.Keep) != 0
}

type Source struct {
	URL       string `json:"url,omitempty"`
	Version   string `json:"version,omitempty"`
	ChartName string `json:"chartName,omitempty"`
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
