package main

import (
	"context"
	"fmt"
	"imperial2metric/pkg"
	"path/filepath"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	outputDir = "output"
)

func downloadHandler(ctx context.Context, s3Event events.S3Event) {
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String("eu-west-2"),
		},
	)

	if err != nil {
		fmt.Printf("error getting a session with aws: %v", err)
	}

	s3 := News3Service(sess)

	for _, record := range s3Event.Records {
		imperialfile := s3.DownloadFile(record.S3.Bucket.Name, record.S3.Object.Key)

		if isAZip(record.S3.Object.Key) {
			fmt.Println("Detectado archivo zip")

			filename := filepath.Base(record.S3.Object.Key)

			finalZip := pkg.ZipProducer(imperialfile, filename)

			err = s3.UploadFile(
				record.S3.Bucket.Name,
				fmt.Sprintf("%s/%s", outputDir, filename),
				finalZip,
			)

			if err != nil {
				fmt.Printf("Cannot upload the file: %s", err.Error())
			}

			return
		}

		metricfile, err := pkg.TransformFile(imperialfile)
		if err != nil {
			fmt.Printf("Error transforming the file: %v", err.Error())
		}

		err = s3.UploadFile(
			record.S3.Bucket.Name,
			fmt.Sprintf("%s/%s", outputDir, filepath.Base(record.S3.Object.Key)),
			metricfile,
		)

		if err != nil {
			fmt.Printf("Cannot upload the file: %s", err.Error())
		}
	}
}
