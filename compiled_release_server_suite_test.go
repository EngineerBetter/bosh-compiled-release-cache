package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCompiledReleaseServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CompiledReleaseServer Suite")
}
