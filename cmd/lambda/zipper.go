package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"

	"github.com/mholt/archiver"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	i2m "imperial2metric/pkg"
)

const (
	dirToDecompress = "."
)

/* func getFilesFromZip(source string) {
	err := archiver.Unarchive(source, dirToDecompress) //https://github.com/mholt/archiver
	if err != nil {
		fmt.Println(err)
	}
} */

// func main() {
// 	getFilesFromZip("es-ES.zip")
// }

var sess *session.Session

func init() {
	// sess = session.Must(session.NewSession(&aws.Config{
	// 	Region: aws.String("eu-west-1"),
	// }))
}

func handler(ctx context.Context, s3Event events.S3Event) {
	// data, _ := json.Marshal(s3Event)

	sess, err := session.NewSession()
	if err != nil {
		fmt.Println(err)
	}

	for _, record := range s3Event.Records {
		fmt.Printf(
			"[%s - %s] Bucket = %s, Key = %s \n",
			record.EventSource,
			record.EventTime,
			record.S3.Bucket.Name,
			record.S3.Object.Key,
		)

		obj, err := s3.New(sess).GetObject(&s3.GetObjectInput{
			Bucket: aws.String(record.S3.Bucket.Name),
			Key:    aws.String(record.S3.Object.Key),
		})

		fmt.Println(obj.Body)

		defer obj.Body.Close()
		if err != nil {
			fmt.Println("Problema al establecer la primera conexion con el s3")
			fmt.Println(err)
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(obj.Body)

		z := archiver.NewZip()

		b, err := ioutil.ReadAll(obj.Body) // The readCloser is the one from the zip-package
		if err != nil {
			panic(err)
		}

		// bytes.Reader implements io.Reader, io.ReaderAt, etc. All you need!
		readerAt := bytes.NewReader(b)

		err = z.Open(readerAt, int64(buf.Len()))
		if err != nil {
			fmt.Println("Problema al abrir el zip")
			fmt.Println(err)
		}

		file, err := z.Read()
		if err != nil {
			fmt.Println("Problema al leer un archivo dentro del zip")
			fmt.Println(err)
		}

		content := i2m.TransformFile(file, record.S3.Object.Key)

		out := new(bytes.Buffer)

		fmt.Fprintf(out, content)

		streamFile := out.Bytes()

		reader := bytes.NewReader(streamFile)

		svc := s3manager.NewUploader(sess)

		result, err := svc.Upload(&s3manager.UploadInput{
			Bucket: aws.String(record.S3.Bucket.Name),
			Key:    aws.String(record.S3.Object.Key),
			Body:   reader,
		})

		if err != nil {
			fmt.Println("Problema al establecer la segunda conexion con el s3")
			fmt.Println(err)
		} else {
			fmt.Println("result was: " + result.Location)
		}
	}
}

func main() {
	lambda.Start(handler)
}
