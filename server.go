package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/engineerbetter/compiled-release-server/bosh"
	"github.com/engineerbetter/compiled-release-server/s3"
	"github.com/gorilla/mux"
)

const s3Bucket = "summit-hackathon-compiled-releases"
const s3Region = "eu-west-1"
const requestTimeoutSeconds = 3600
const requestTimeout = time.Second * requestTimeoutSeconds

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", releaseHandler).
		Queries("release", "{release}").
		Queries("release_v", "{release_v}").
		Queries("stemcell", "{stemcell}").
		Queries("stemcell_v", "{stemcell_v}")

	srv := &http.Server{
		Handler:           router,
		Addr:              ":8080",
		ReadTimeout:       requestTimeout,
		WriteTimeout:      requestTimeout,
		ReadHeaderTimeout: requestTimeout,
	}

	log.Fatal(srv.ListenAndServe())
}
func streamFromBoshIO(w http.ResponseWriter, r bosh.CompiledRelease) error {
	boshResp, err := http.Get(r.BoshURL())
	if err != nil {
		return err
	}
	_, err = io.Copy(w, boshResp.Body)
	if err != nil {
		return err
	}
	return nil
}
func releaseHandler(w http.ResponseWriter, r *http.Request) {
	release := ReleaseFromRequestVars(mux.Vars(r))
	path := release.ToS3Path()
	log.Println("Writing headers")

	w.Header().Set("Content-Type", "application/gzip")
	w.Header().Set(
		"Keep-Alive",
		fmt.Sprintf("timeout=%d, max=%d", requestTimeoutSeconds, 100),
	)

	fileReader, _, err := s3.GetFile(s3Bucket, path, s3Region)

	if strings.HasPrefix(err.Error(), s3.AWSErrCodeNoSuchKey) {
		// compile release
		go compile(&release)

		err := streamFromBoshIO(w, release)
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	} else {
		log.Println("goto s3")
		// found in bucket
		_, err = io.Copy(w, fileReader)
		if err != nil {
			panic(err)
		}
	}

	return
}

func compile(release *bosh.CompiledRelease) {
	log.Printf("Compiling release: %s\n", release.ReleaseName)
	client := bosh.New(os.Getenv("BOSH_USER"), os.Getenv("BOSH_PASSWORD"), os.Getenv("BOSH_HOST"), os.Getenv("BOSH_CA_CERT"))
	output, err := client.Compile(release)
	if err != nil {
		panic(err)
	}

	if err := s3.PutFile(s3Bucket, release.ToS3Path(), s3Region, output); err != nil {
		panic(err)
	}
}

func ReleaseFromRequestVars(requestVars map[string]string) bosh.CompiledRelease {
	return bosh.CompiledRelease{
		ReleaseName:     requestVars["release"],
		ReleaseVersion:  requestVars["release_v"],
		StemcellName:    requestVars["stemcell"],
		StemcellVersion: requestVars["stemcell_v"],
	}
}
