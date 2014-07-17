package test_support

import(
	"io/ioutil"
	"fmt"
	"path/filepath"
//	"os"
	"runtime"
)


type FakeCfClient struct {
}

func(cf *FakeCfClient) GetSpaces() string {
	bytes, err := ioutil.ReadFile(getFixturePath("successful_get_spaces_response.json"))
	if err != nil {
		fmt.Println(err.Error())
	}

	return string(bytes)
}

func(cf *FakeCfClient) GetServiceInstancesForSpace(space_guid string) string {
	var filename string
	switch space_guid {
		case "space-0": filename = "successful_get_instances_for_space_0_response.json"
		case "space-1": filename = "successful_get_instances_for_space_1_response.json"
	}

	bytes, err := ioutil.ReadFile(getFixturePath(filename))
	if err != nil {
		fmt.Println(err.Error())
	}
	return string(bytes)
}

func getFixturePath(filename string) string {
	_, current_filename, _, _ := runtime.Caller(0)
	current_dir := filepath.Dir(current_filename)
	return fmt.Sprintf("%s/cf_responses/%s", current_dir, filename)
}
