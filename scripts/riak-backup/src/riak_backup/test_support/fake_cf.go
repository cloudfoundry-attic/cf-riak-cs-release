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
	return readFixtureFile("successful_get_spaces_response.json")
}

func(cf *FakeCfClient) GetServiceInstancesForSpace(space_guid string) string {
	var filename string
	switch space_guid {
		case "space-0": filename = "successful_get_instances_for_space_0_response.json"
		case "space-1": filename = "successful_get_instances_for_space_1_response.json"
		default: panic("fixture file not found")
	}

	return readFixtureFile(filename)
}

func(cf *FakeCfClient) GetBindings(service_instance_guid string) string {
	var filename string
	switch service_instance_guid {
		case "service-instance-0": filename = "successful_get_bindings_for_service_instance_0_response.json"
		case "service-instance-1", "service-instance-2", "service-instance-3": return "{}"
		default: panic(fmt.Sprintf("fixture file not found for %s", service_instance_guid))
	}

	return readFixtureFile(filename)
}

func readFixtureFile(filename string) string {
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
