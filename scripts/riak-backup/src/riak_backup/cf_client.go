package riak_backup

import (
	"os"
	"os/exec"
	"fmt"
)

type CfClientInterface interface {
	GetSpaces(next_url string) string
	GetOrganization(organization_guid string) string
	GetServiceInstancesForSpace(space_guid string) string
	GetBindings(service_instance_guid string) string
	Login(user, password string)
}

type CfClient struct {
}

func(cf *CfClient) GetSpaces(next_url string) string {
	cmd := exec.Command("cf", "curl", next_url)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("failed to get spaces using URL: ", next_url)
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return string(output)
}

func(cf *CfClient) GetOrganization(organization_guid string) string {
	organization_url := "/v2/organizations/" + organization_guid
	cmd := exec.Command("cf", "curl", organization_url)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("failed to get organization using URL: ", organization_url)
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return string(output)
}

func (cf *CfClient) GetServiceInstancesForSpace(space_guid string) string {
	summary_url := "/v2/spaces/" + space_guid + "/summary"
	cmd := exec.Command("cf", "curl", summary_url)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("failed to get spaces using URL: ", summary_url)
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return string(output)
}

func (cf *CfClient) GetBindings(next_url string) string {
	cmd := exec.Command("cf", "curl", next_url)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("failed to get bindings using URL: ", next_url)
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return string(output)
}

func (cf *CfClient) Login(user, password string) {
	cmd := exec.Command("cf", "auth", user, password)
	output, err := cmd.CombinedOutput()
	fmt.Println("Command: cf auth " + user + " " + password)
	fmt.Println(string(output))
	if err != nil {
		os.Exit(1)
	}
}
