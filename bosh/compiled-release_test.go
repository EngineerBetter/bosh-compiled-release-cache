package bosh_test

import (
	. "github.com/engineerbetter/compiled-release-server/bosh"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CompiledRelease", func() {
	Describe("ToS3Path", func() {
		It("Returns the correct path", func() {
			release := CompiledRelease{
				ReleasePath:     "github.com/cloudfoundry/garden-runc-release",
				ReleaseVersion:  "1.7.0",
				StemcellName:    "bosh-aws-xen-hvm-ubuntu-trusty-go_agent",
				StemcellVersion: "3421.6",
			}

			expectedPath := "github.com-cloudfoundry-community-docker-registry-boshrelease-11-bosh-aws-xen-hvm-ubuntu-trusty-go_agent-3421.6.tgz"

			Expect(release.ToS3Path()).To(Equal(expectedPath))
		})
	})
})
