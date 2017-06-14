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
	uuid "github.com/satori/go.uuid"
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
	log.Printf("Streaming response from %s\n", r.BoshURL())
	boshResp, err := http.Get(r.BoshURL())
	if err != nil {
		return err
	}

	if boshResp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status code 200 from %s, recieved %d\n", r.BoshURL(), boshResp.StatusCode)
	}

	byteCount, err := io.Copy(w, boshResp.Body)
	if err != nil {
		return err
	}

	log.Printf("Copied %d bytes\n", byteCount)

	return nil
}
func releaseHandler(w http.ResponseWriter, r *http.Request) {
	release := ReleaseFromRequestVars(mux.Vars(r))
	log.Printf("fetching release %#v\n", release)

	path := release.ToS3Path()
	log.Println("writing headers")

	w.Header().Set("Content-Type", "application/gzip")
	w.Header().Set(
		"Keep-Alive",
		fmt.Sprintf("timeout=%d, max=%d", requestTimeoutSeconds, 100),
	)

	fileReader, err := s3.GetFile(s3Bucket, path, s3Region)

	if err != nil && strings.HasPrefix(err.Error(), s3.AWSErrCodeNoSuchKey) {
		go compile(release)

		if err := streamFromBoshIO(w, release); err != nil {
			log.Fatal(err)
		}

		return
	}

	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("found release %s in S3\n", release.ReleasePath)

	bytesCount, err := io.Copy(w, fileReader)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("copied %d bytes from S3\n", bytesCount)
}

func compile(release bosh.CompiledRelease) {
	log.Printf("compiling release: %s\n", release.ReleasePath)

	client := bosh.New(os.Getenv("BOSH_CLIENT"), os.Getenv("BOSH_CLIENT_SECRET"), os.Getenv("BOSH_HOST"), os.Getenv("BOSH_CA_CERT"))
	output, err := client.Compile(&release)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("uploading release to s3: %s\n", release.ToS3Path())
	if err := s3.PutFile(s3Bucket, release.ToS3Path(), s3Region, output); err != nil {
		log.Fatal(err)
		return
	}
}

func ReleaseFromRequestVars(requestVars map[string]string) bosh.CompiledRelease {
	return bosh.CompiledRelease{
		DeploymentName:  fmt.Sprintf("compilation-%s", uuid.NewV4().String()),
		ReleasePath:     requestVars["release"],
		ReleaseVersion:  requestVars["release_v"],
		StemcellName:    requestVars["stemcell"],
		StemcellVersion: requestVars["stemcell_v"],
	}
}
