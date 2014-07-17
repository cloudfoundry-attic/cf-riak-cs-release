package riak_backup

import (
	"fmt"
	"encoding/json"
	"os"
	"gopkg.in/v1/yaml"
	"io/ioutil"
)

type Spaces struct{
	Resources []Space
}

type Space struct {
	Metadata struct {
		Guid string
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
	spaces_json := cf.GetSpaces()
	spaces := &Spaces{}
	json.Unmarshal([]byte(spaces_json), spaces)

	for _, space := range spaces.Resources {
		space_guid := space.Metadata.Guid
		os.MkdirAll(fmt.Sprintf("/tmp/backup/spaces/%s", space_guid), 0777)

		service_instances_json := cf.GetServiceInstancesForSpace(space_guid)
		service_instances := &ServiceInstances{}
		json.Unmarshal([]byte(service_instances_json), service_instances)

		for _, service_instance := range service_instances.Services {
			if service_instance.ServicePlan.Service.Label == "p-riakcs" {
				service_instance_guid := service_instance.Guid
				os.MkdirAll(fmt.Sprintf("/tmp/backup/spaces/%s/service_instances/%s", space_guid, service_instance_guid), 0777)
				writeMetadataFile(cf, space_guid, service_instance_guid)
			}
		}
	}
}

func writeMetadataFile(cf CfClientInterface, space_guid string, service_instance_guid string) {
	bindings_json := cf.GetBindings(service_instance_guid)
	bindings := &Bindings{}

	json.Unmarshal([]byte(bindings_json), bindings)

	metadata := InstanceMetadata{
		ServiceInstanceGuid: service_instance_guid,
	}

	app_metadatas := []AppMetadata{}
	for _, binding := range bindings.Resources {
		bound_app := binding.Entity.App
		app_metadatas = append(app_metadatas, AppMetadata{ Name: bound_app.Entity.Name, Guid: bound_app.Metadata.Guid })
	}
	metadata.BoundApps = app_metadatas

	bytes, err := yaml.Marshal(metadata)
	if err != nil {
		fmt.Println(err.Error())
	}

	path := fmt.Sprintf("/tmp/backup/spaces/%s/service_instances/%s/metadata.yml", space_guid, service_instance_guid)
	ioutil.WriteFile(path, bytes, 0777)
}
