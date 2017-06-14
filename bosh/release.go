package bosh

import (
	"fmt"
	"strings"
)

type CompiledRelease struct {
	DeploymentName  string
	ReleaseName     string
	ReleasePath     string
	ReleaseVersion  string
	StemcellName    string
	StemcellVersion string
}

func (compiledRelease *CompiledRelease) ToS3Path() string {
	return fmt.Sprintf("%s-%s-%s-%s.tgz",
		sanitizeS3Path(compiledRelease.ReleasePath),
		sanitizeS3Path(compiledRelease.ReleaseVersion),
		sanitizeS3Path(compiledRelease.StemcellName),
		sanitizeS3Path(compiledRelease.StemcellVersion),
	)
}

func (compiledRelease *CompiledRelease) BoshURL() string {
	return "https://bosh.io/d/" + compiledRelease.ReleaseName + "?v=" + compiledRelease.ReleaseVersion
}

func sanitizeS3Path(path string) string {
	return strings.Replace(path, "/", "-", -1)
}
