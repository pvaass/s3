package downloader

import (
	"os"
	"path/filepath"

	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Downloader struct {
	*s3manager.Downloader
	Bucket string
}

func (d *Downloader) Get(remotePath string, localPath string) {
	if err := os.MkdirAll(filepath.Dir(localPath), 0775); err != nil {
		panic(err)
	}

	file, err := os.Create(localPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fmt.Printf("\nDownloading %s to %s", remotePath, localPath)

	n, err := d.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(d.Bucket),
		Key:    aws.String(remotePath),
	})

	if err != nil {
		panic(fmt.Errorf(": failed to download file or directory, %v", err))
	}
	fmt.Printf(": file downloaded, %d bytes\n", n)
}
