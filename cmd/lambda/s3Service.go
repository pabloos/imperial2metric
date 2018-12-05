package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//S3Service manages all the operations of the s2 service
type S3Service struct {
	downloader *s3manager.Downloader
	uploader   *s3manager.Uploader
}

// News3Service returns a new s3Service object
func News3Service(sess *session.Session) *S3Service {
	return &S3Service{
		downloader: s3manager.NewDownloader(sess),
		uploader:   s3manager.NewUploader(sess),
	}
}

// DownloadFile downloads a file as a byte array
func (s3s *S3Service) DownloadFile(bucket, filename string) io.Reader {
	buf := aws.NewWriteAtBuffer([]byte{}) // https://stackoverflow.com/questions/46019484/buffer-implementing-io-writerat-in-go

	filename, err := url.QueryUnescape(filename)
	if err != nil {
		fmt.Printf("Cannot convert the file to %v", err)
	}

	_, err = s3s.downloader.Download(buf,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(filename),
		},
	)

	if err != nil {
		log.Fatalf("Unable to download item %q, %v", filename, err)
	}

	return bytes.NewReader(buf.Bytes())
}

// UploadFile uploads a file to the bucket and dir passed
func (s3s *S3Service) UploadFile(bucket, key string, file io.Reader) error {
	// s3 API does not allow to use an io.Reader to pass to the body field directly. In fact it does but doesn't work --> 0 bytes file in the destiny
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(file)

	key, err := url.QueryUnescape(key)
	if err != nil {
		return err
	}

	_, err = s3s.uploader.Upload(
		&s3manager.UploadInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			// ContentDisposition: aws.String("attachment"),
			// ContentType:        aws.String(http.DetectContentType(buf.Bytes())),
			Body: buffer,
		},
	)

	return err
}
