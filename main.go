package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func getAwsConfig() (aws.Config, error) {
	awsAccessKey, awsAccessKeyExists := os.LookupEnv("PB2S3_AWS_ACCESS_KEY")
	awsSecretKey, awsSecretKeyExists := os.LookupEnv("PB2S3_AWS_SECRET_KEY")
	if awsAccessKeyExists && awsSecretKeyExists {
		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsAccessKey, awsSecretKey, "")),
		)
		if err != nil {
			return aws.Config{}, errors.New("config not set")
		}
		return cfg, nil
	} else {
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return *aws.NewConfig(), errors.New("config not set")
		}
		return cfg, nil
	}
}

func main() {
	bucketName := os.Getenv("PB2S3_BUCKET_NAME")

	cfg, err := getAwsConfig()
	if err != nil {
		log.Println(err)
		return
	}

	tempDir, err := ioutil.TempDir("", "piholebackup")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	os.Chdir(tempDir)
	teleporter := exec.Command("pihole", "-a", "-t")
	if err := teleporter.Run(); err != nil {
		log.Fatal(err)
	}

	files, err := filepath.Glob("./pi-hole*")
	if err != nil || len(files) == 0 {
		log.Println(err)
		return
	}
	piholeBackup := files[0]

	file, err := os.Open(piholeBackup)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	input := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    &piholeBackup,
		Body:   file,
	}

	svc := s3.NewFromConfig(cfg)

	_, err = svc.PutObject(context.TODO(), input)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("uploaded " + piholeBackup + " to bucket " + bucketName)
}
