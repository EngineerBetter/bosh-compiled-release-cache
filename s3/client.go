package s3

import (
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	// New versions of github.com/aws/aws-sdk-go/aws have these consts
	// but the version currently pinned by bosh-cli v2 does not

	// ErrCodeNoSuchBucket for service response error code
	// "NoSuchBucket".
	//
	// The specified bucket does not exist.
	awsErrCodeNoSuchBucket = "NoSuchBucket"

	// AWSErrCodeNoSuchKey for service response error code
	// "NoSuchKey".
	//
	// The specified key does not exist.
	AWSErrCodeNoSuchKey = "NoSuchKey"

	// Returned when calling HEAD on non-existant bucket or object
	awsErrCodeNotFound = "NotFound"
)

// HasFile returns true if the specified S3 object exists
func HasFile(bucket, path, region string) (bool, error) {
	sess, err := session.NewSession(aws.NewConfig().WithCredentialsChainVerboseErrors(true))
	if err != nil {
		return false, err
	}
	client := s3.New(sess, &aws.Config{Region: &region})

	_, err = client.HeadObject(&s3.HeadObjectInput{Bucket: &bucket, Key: &path})
	if err != nil {
		awsErrCode := err.(awserr.Error).Code()
		if awsErrCode == awsErrCodeNotFound || awsErrCode == AWSErrCodeNoSuchKey {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// GetFile returns a file on S3
func GetFile(bucket, path, region string) (io.ReadCloser, string, error) {
	sess, err := session.NewSession(aws.NewConfig().WithCredentialsChainVerboseErrors(true))
	if err != nil {
		return nil, "", err
	}
	client := s3.New(sess, &aws.Config{Region: &region})
	response, err := client.GetObject(&s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &path,
	})

	if err != nil {
		return nil, "", err
	}

	status := response.Metadata["Status"]
	fmt.Printf("%#v\n", response.Metadata)
	if status == nil {
		return nil, "", errors.New("Metadata not set: Status")
	}

	return response.Body, *status, nil
}

// PutFile writes a file to S3
func PutFile(bucket, path, region string, file io.ReadSeeker) error {
	sess, err := session.NewSession(aws.NewConfig().WithCredentialsChainVerboseErrors(true))
	if err != nil {
		return err
	}
	client := s3.New(sess, &aws.Config{Region: &region})

	acl := "public-read"

	_, err = client.PutObject(&s3.PutObjectInput{
		ACL:    &acl,
		Body:   file,
		Bucket: &bucket,
		Key:    &path,
	})

	return err
}
