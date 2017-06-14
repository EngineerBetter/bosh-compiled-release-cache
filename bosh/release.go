package bosh

import (
	"fmt"
	"strings"
)

type CompiledRelease struct {
	DeploymentName  string
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

func (compiledRelease *CompiledRelease) ReleaseName() string {
	releasePathParts := strings.Split(compiledRelease.ReleasePath, "/")
	releaseName := releasePathParts[len(releasePathParts)-1]
	releaseName = strings.TrimSuffix(releaseName, "-release")
	releaseName = strings.TrimSuffix(releaseName, "-boshrelease")
	return releaseName
}

func (compiledRelease *CompiledRelease) StemcellOS() string {
	if strings.Contains(compiledRelease.StemcellName, "ubuntu-trusty") {
		return "ubuntu-trusty"
	}

	panic(fmt.Errorf("Stemcell not supported: %s", compiledRelease.StemcellName))
}

func (compiledRelease *CompiledRelease) StemcellURL() string {
	// https://bosh.io/d/stemcells/bosh-aws-xen-hvm-ubuntu-trusty-go_agent?v=3421.6
	return fmt.Sprintf("https://bosh.io/d/stemcells/%s?v=%s", compiledRelease.StemcellName, compiledRelease.StemcellVersion)
}

func (compiledRelease *CompiledRelease) BoshURL() string {
	return "https://bosh.io/d/" + compiledRelease.ReleasePath + "?v=" + compiledRelease.ReleaseVersion
}

func sanitizeS3Path(path string) string {
	return strings.Replace(path, "/", "-", -1)
}
