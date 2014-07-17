package riak_backup

import (
	"fmt"
	"encoding/json"
	"os"
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
	Resources []ServiceInstance
}

type ServiceInstance struct {
	Metadata struct {
		Guid string
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

		for _, service_instance := range service_instances.Resources {
			service_instance_guid := service_instance.Metadata.Guid
			os.MkdirAll(fmt.Sprintf("/tmp/backup/spaces/%s/service_instances/%s", space_guid, service_instance_guid), 0777)
		}
	}
}
