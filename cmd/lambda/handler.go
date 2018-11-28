package main

import (
	"context"
	"fmt"
	"path/filepath"

	i2m "imperial2metric/pkg"

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
		fmt.Printf(
			"[%s - %s] Bucket = %s, Key = %s \n",
			record.EventSource,
			record.EventTime,
			record.S3.Bucket.Name,
			record.S3.Object.Key, //key contains both dir (if any) and the object name
		)

		imperialfile := s3.DownloadFile(record.S3.Bucket.Name, record.S3.Object.Key)

		metricfile := i2m.TransformFile(imperialfile)

		err = s3.UploadFile(
			record.S3.Bucket.Name,
			fmt.Sprintf("%s/%s", outputDir, filepath.Base(record.S3.Object.Key)),
			metricfile,
		)

		if err != nil {
			fmt.Printf("Cannot upload the file: %s", err)
		}
	}
}
