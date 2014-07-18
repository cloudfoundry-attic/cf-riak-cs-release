package test_support

import (
	"io/ioutil"
)

type FakeS3CmdClient struct {
}

func(s3cmd *FakeS3CmdClient) FetchBucket(bucket_name string, destination_dir string) {
	path := destination_dir + "/datafile.dat"
	bytes := []byte(`data from bucket ` + bucket_name)
	ioutil.WriteFile(path, bytes, 0644)
}
