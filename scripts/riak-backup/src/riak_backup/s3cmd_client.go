package riak_backup

import (
	"fmt"
	"os"
	"os/exec"
)

type S3CmdClientInterface interface {
	FetchBucket(bucket_name string, destination_dir string)
}

type S3CmdClient struct {
	config_file string
}

func NewS3CmdClient(config_file string) *S3CmdClient {
	return &S3CmdClient{config_file: config_file}
}

func(s3cmd *S3CmdClient) FetchBucket(bucket_name string, destination_dir string) {
	fmt.Printf("Backing up bucket %s to %s\n", bucket_name, destination_dir)

	bucket_path := "s3://" + bucket_name
	cmd := exec.Command("s3cmd", "-c", s3cmd.config_file, "sync", bucket_path, destination_dir)

	output, err := cmd.Output()
	fmt.Println(string(output))
	if err != nil {
		fmt.Println("failed to back up contents of bucket ", bucket_name)
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
