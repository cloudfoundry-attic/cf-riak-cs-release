package main

import (
	"fmt"
	"riak_backup"
	"os"
)

func main() {
	if len(os.Args) != 5 {
		printUsage()
	}

	operation := os.Args[1]
	s3cfg := os.Args[2]
	cf_user := os.Args[3]
	cf_password := os.Args[4]

	if _, err := os.Stat(s3cfg); os.IsNotExist(err) {
		fmt.Printf("no such file or directory: %s", s3cfg)
		printUsage()
	}

	cf_client := riak_backup.CfClient{}
	cf_client.Login(cf_user, cf_password)

	s3cmd_client := riak_backup.S3CmdClient{}

	switch operation {
		case "backup": riak_backup.Backup(&cf_client, &s3cmd_client)
		case "restore": fmt.Println("not implemented")
		default: printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage: riak-backup [backup|restore] PATH_TO_S3CFG_FILE CF_ADMIN_USER CF_ADMIN_PASSWORD")
	os.Exit(1)
}

