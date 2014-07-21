package main

import (
	"fmt"
	"riak_backup"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) != 5 {
		printUsage()
	}

	ensureS3cmdIsInstalled()

	backup_dir := os.Args[1]
	s3cfg := os.Args[2]
	cf_user := os.Args[3]
	cf_password := os.Args[4]

	if _, err := os.Stat(s3cfg); os.IsNotExist(err) {
		fmt.Printf("no such file or directory: %s", s3cfg)
		printUsage()
	}

	cf_client := riak_backup.CfClient{}
	cf_client.Login(cf_user, cf_password)

	s3cmd_client := *riak_backup.NewS3CmdClient(s3cfg)

	riak_backup.Backup(&cf_client, &s3cmd_client, backup_dir)
}

func printUsage() {
	fmt.Println("Usage: riak-backup BACKUP_DESTINATION_DIR PATH_TO_S3CFG_FILE CF_ADMIN_USER CF_ADMIN_PASSWORD")
	os.Exit(1)
}

func ensureS3cmdIsInstalled() {
	cmd := exec.Command("which", "s3cmd")
	_, err := cmd.Output()
	if err != nil {
		fmt.Println("Please install s3cmd and make sure it is in the $PATH before running this script.")
		os.Exit(1)
	}
}

