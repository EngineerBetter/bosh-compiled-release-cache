package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/engineerbetter/compiled-release-server/bosh"
	"github.com/satori/go.uuid"
)

func main() {
	client := bosh.New(os.Getenv("BOSH_USER"), os.Getenv("BOSH_PASSWORD"), os.Getenv("BOSH_HOST"), os.Getenv("BOSH_CA_CERT"))

	file, err := client.Compile(&bosh.CompiledRelease{
		// deployment name must start with a letter
		DeploymentName:  "a" + uuid.NewV4().String(),
		ReleaseName:     "bosh",
		ReleasePath:     "github.com/cloudfoundry/bosh",
		ReleaseVersion:  "262.1",
		StemcellName:    "ubuntu-trusty",
		StemcellVersion: "3363.20",
	})
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	fmt.Println("SUCCESSFULLY READ FILE")
}
