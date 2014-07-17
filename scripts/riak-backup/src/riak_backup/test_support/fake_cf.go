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
	bytes, err := ioutil.ReadFile(getFixturePath("successful_get_instances_for_space_0_response.json"))
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
