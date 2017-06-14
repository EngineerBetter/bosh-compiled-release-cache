package main_test

import (
	. "github.com/engineerbetter/compiled-release-server"
	"github.com/engineerbetter/compiled-release-server/bosh"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	Describe("ReleaseFromRequestVars", func() {
		It("Returns the correct compiled release", func() {
			vars := map[string]string{
				"release":    "github.com/cloudfoundry-community/docker-registry-boshrelease",
				"release_v":  "11",
				"stemcell":   "bosh-aws-xen-hvm-ubuntu-trusty-go_agent",
				"stemcell_v": "3421.6",
			}

			Expect(ReleaseFromRequestVars(vars)).To(Equal(bosh.CompiledRelease{
				ReleaseName:     "github.com/cloudfoundry-community/docker-registry-boshrelease",
				ReleaseVersion:  "11",
				StemcellName:    "bosh-aws-xen-hvm-ubuntu-trusty-go_agent",
				StemcellVersion: "3421.6",
			}))
		})
	})
})
