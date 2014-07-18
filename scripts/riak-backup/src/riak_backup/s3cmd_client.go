package riak_backup

type S3CmdClientInterface interface {
	FetchBucket(bucket_name string, destination_dir string)
}

type S3CmdClient struct {
}

func(s3cmd *S3CmdClient) FetchBucket(bucket_name string, destination_dir string) {

}
