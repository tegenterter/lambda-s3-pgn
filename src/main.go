package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Response struct {
	Objects []string `json:"Objects"`
}

func handler(ctx context.Context, s3Event events.S3Event) (Response, error) {
	var response Response

	session, err := session.NewSession()
	if err != nil {
		return response, err
	}

	downloader := s3manager.NewDownloader(session)
	uploader := s3manager.NewUploader(session)

	for _, record := range s3Event.Records {
		buffer := &aws.WriteAtBuffer{}

		downloader.Download(buffer, &s3.GetObjectInput{
			Bucket: aws.String(record.S3.Bucket.Name),
			Key:    aws.String(record.S3.Object.Key),
		})

		originalPath := "/tmp/original-" + record.S3.Object.Key
		processedPath := "/tmp/processed-" + record.S3.Object.Key

		ioutil.WriteFile(originalPath, buffer.Bytes(), 0644)

		output, err := exec.Command("/bin/pgn-extract", originalPath, "--fencomments", "--output", processedPath).CombinedOutput()
		fmt.Println(string(output))
		if err != nil {
			return response, err
		}

		file, err := os.Open(processedPath)
		if err != nil {
			return response, err
		}
		defer file.Close()

		result, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(os.Getenv("AWS_S3_BUCKET")),
			Key:    aws.String(record.S3.Object.Key),
			Body:   file,
		})
		if err != nil {
			return response, err
		}

		response.Objects = append(response.Objects, result.Location)
	}

	return response, nil
}

func main() {
	lambda.Start(handler)
}
