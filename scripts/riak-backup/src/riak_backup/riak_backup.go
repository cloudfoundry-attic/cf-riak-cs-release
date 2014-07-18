package riak_backup

import (
	"fmt"
	"encoding/json"
	"os"
	"gopkg.in/v1/yaml"
	"io/ioutil"
)

type Spaces struct{
	NextUrl string `json:"next_url"`
	Resources []Space
}

type Space struct {
	Metadata struct {
		Guid string
	}
	Entity struct {
		Name string
	}
}

type ServiceInstances struct {
	Services []ServiceInstance
}

type ServicePlan struct {
	Service struct {
		Label string
	}
}

type ServiceInstance struct {
	Guid string
	ServicePlan ServicePlan `json:"service_plan"`
}

type Bindings struct {
	NextUrl string `json:"next_url"`
	Resources []Binding
}

type Binding struct {
	Entity struct {
		App App
	}
}

type App struct {
	Metadata struct {
		Guid string
	}
	Entity struct {
		Name string
	}
}

func Backup(cf CfClientInterface) {
	spaces := fetchSpaces(cf)

	for _, space := range spaces {
		space_guid := space.Metadata.Guid
		space_name := space.Entity.Name
		os.MkdirAll(fmt.Sprintf("/tmp/backup/spaces/%s", space_name), 0777)

		service_instances_json := cf.GetServiceInstancesForSpace(space_guid)
		service_instances := &ServiceInstances{}
		json.Unmarshal([]byte(service_instances_json), service_instances)

		for _, service_instance := range service_instances.Services {
			if service_instance.ServicePlan.Service.Label == "p-riakcs" {
				service_instance_guid := service_instance.Guid
				os.MkdirAll(fmt.Sprintf("/tmp/backup/spaces/%s/service_instances/%s", space_name, service_instance_guid), 0777)
				writeMetadataFile(cf, space_name, service_instance_guid)
			}
		}
	}
}

func fetchSpaces(cf CfClientInterface) []Space {
	spaces := []Space{}
	next_url := "/v2/spaces"
	for next_url != "" {
		spaces_json := cf.GetSpaces(next_url)
		page := &Spaces{}
		json.Unmarshal([]byte(spaces_json), page)

		spaces = append(spaces, page.Resources...)
		next_url = page.NextUrl
	}
	return spaces
}

func fetchBindings(cf CfClientInterface, service_instance_guid string) []Binding {
	next_url := "/v2/service_instances/" + service_instance_guid + "/service_bindings?inline-relations-depth=1"
	bindings := []Binding{}
	for next_url != "" {
		bindings_json := cf.GetBindings(next_url)
		page := &Bindings{}
		json.Unmarshal([]byte(bindings_json), page)

		bindings = append(bindings, page.Resources...)
		next_url = page.NextUrl
	}
	return bindings
}

func writeMetadataFile(cf CfClientInterface, space_name string, service_instance_guid string) {
	bindings := fetchBindings(cf, service_instance_guid)

	metadata := InstanceMetadata{
		ServiceInstanceGuid: service_instance_guid,
	}

	app_metadatas := []AppMetadata{}
	for _, binding := range bindings {
		bound_app := binding.Entity.App
		app_metadatas = append(app_metadatas, AppMetadata{ Name: bound_app.Entity.Name, Guid: bound_app.Metadata.Guid })
	}
	metadata.BoundApps = app_metadatas

	bytes, err := yaml.Marshal(metadata)
	if err != nil {
		fmt.Println(err.Error())
	}

	path := fmt.Sprintf("/tmp/backup/spaces/%s/service_instances/%s/metadata.yml", space_name, service_instance_guid)
	ioutil.WriteFile(path, bytes, 0777)
}
