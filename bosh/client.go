package bosh

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/engineerbetter/compiled-release-server/util"
)

// Client is a struct representing a BOSH client
type Client struct {
	username string
	password string
	ip       string
	caCert   string
}

// New returns a new bosh client
func New(username, password, ip, caCert string) Client {
	return Client{
		username: username,
		password: password,
		ip:       ip,
		caCert:   caCert,
	}
}

func (client *Client) Compile(release *CompiledRelease) (io.ReadCloser, error) {
	manifestBytes, err := GenerateManifest(release)
	if err != nil {
		return nil, err
	}

	manifestFile, err := ioutil.TempFile("", "manifest")
	if err != nil {
		return nil, err
	}
	defer manifestFile.Close()

	if _, err := manifestFile.Write(manifestBytes); err != nil {
		return nil, err
	}

	if err := client.deploy(release.DeploymentName, manifestFile.Name()); err != nil {
		return nil, err
	}

	dir, err := ioutil.TempDir("", "compiled-release-server")
	if err != nil {
		return nil, err
	}

	if err := client.exportRelease(release, dir); err != nil {
		return nil, err
	}

	if err := client.deleteDeployment(release.DeploymentName); err != nil {
		return nil, err
	}

	filesInDir, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var releasePath string

	for _, file := range filesInDir {
		if strings.HasSuffix(file.Name(), ".tgz") {
			releasePath = filepath.Join(dir, file.Name())
		}
	}

	if releasePath == "" {
		return nil, fmt.Errorf("No release .tgz found in %s", dir)
	}

	return os.Open(releasePath)
}

func (client *Client) exportRelease(release *CompiledRelease, dir string) error {
	cmd := exec.Command(
		"bosh-cli",
		"--non-interactive",
		"--environment",
		client.ip,
		"--ca-cert",
		client.caCert,
		"--client",
		client.username,
		"--client-secret",
		client.password,
		"--deployment",
		release.DeploymentName,
		"export-release",
		"--dir",
		dir,
		fmt.Sprintf("%s/%s", release.ReleaseName, release.ReleaseVersion),
		fmt.Sprintf("%s/%s", release.StemcellName, release.StemcellVersion),
	)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func (client *Client) deploy(deploymentName, manifestPath string) error {
	cmd := exec.Command(
		"bosh-cli",
		"--non-interactive",
		"--environment",
		client.ip,
		"--ca-cert",
		client.caCert,
		"--client",
		client.username,
		"--client-secret",
		client.password,
		"--deployment",
		deploymentName,
		"deploy",
		manifestPath,
	)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func (client *Client) deleteDeployment(deploymentName string) error {

	cmd := exec.Command(
		"bosh-cli",
		"--non-interactive",
		"--environment",
		client.ip,
		"--ca-cert",
		client.caCert,
		"--client",
		client.username,
		"--client-secret",
		client.password,
		"--deployment",
		deploymentName,
		"delete-deployment",
	)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

const manifestTemplate = `
---
name: <% .DeploymentName %>

releases:
- name: <% .ReleaseName %>
  url: <% .BoshURL %>
  version: "<% .ReleaseVersion %>"

stemcells:
- alias: default
  os: <% .StemcellName %>
  version: "<% .StemcellVersion %>"

jobs: []

update:
  canaries: 1
  max_in_flight: 1
  canary_watch_time: 1000-90000
  update_watch_time: 1000-90000
`

func GenerateManifest(compiledRelease *CompiledRelease) ([]byte, error) {
	bytes, err := util.RenderTemplate(manifestTemplate, compiledRelease)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
