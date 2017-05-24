package main

import (
	"flag"

	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awss3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	s3 "github.com/pvaass/s3/pkg/downloader"
)

func main() {
	options := parseOpts()
	validateOpts(options)

	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("eu-central-1"),
		Credentials: credentials.NewStaticCredentials(options.KeyID, options.KeySecret, options.token),
	}))

	s3 := s3.Downloader{
		s3manager.NewDownloader(sess),
		options.Bucket,
	}
	svc := awss3.New(sess)

	params := &awss3.ListObjectsInput{
		Bucket: aws.String(options.Bucket),
		Prefix: aws.String(options.BucketPath),
	}

	resp, _ := svc.ListObjects(params)
	for _, key := range resp.Contents {
		remotePath := *key.Key
		localPath := strings.Replace(remotePath, options.BucketPath, options.LocalPath, 1)

		s3.Get(remotePath, localPath)
	}

}

type Options struct {
	KeyID      string
	KeySecret  string
	Bucket     string
	BucketPath string
	LocalPath  string
	Region     string
	token      string
}

func parseOpts() Options {
	idPtr := flag.String("key_id", "", "Your AWS Access Key ID")
	secretPtr := flag.String("key_secret", "", "Your AWS Access Key Secret")
	bucketPtr := flag.String("bucket", "", "The bucket to download from")
	regionPtr := flag.String("region", "", "The region your bucket is in")
	flag.Parse()

	return Options{
		KeyID:      *idPtr,
		KeySecret:  *secretPtr,
		Bucket:     *bucketPtr,
		BucketPath: flag.Arg(0),
		LocalPath:  flag.Arg(1),
		Region:     *regionPtr,
	}
}

func validateOpts(opts Options) {
	isValidOpts := func(opts Options) bool {
		return opts.KeyID != "" &&
			opts.KeySecret != "" &&
			opts.Bucket != "" &&
			opts.BucketPath != "" &&
			opts.LocalPath != "" &&
			opts.Region != ""
	}

	if !isValidOpts(opts) {
		panic("Invalid arguments. Try --help for information on how to use this tool")
	}
}
