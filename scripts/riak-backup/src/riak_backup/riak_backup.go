package riak_backup

import (
	"fmt"
	"encoding/json"
	"os"
	"gopkg.in/v1/yaml"
	"io/ioutil"
)

type Organization struct {
	Entity struct {
		Name string
	}
}

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
		OrganizationGuid string `json:"organization_guid"`
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
	Name string
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

func Backup(cf CfClientInterface, s3cmd S3CmdClientInterface) {
	spaces := fetchSpaces(cf)

	for _, space := range spaces {
		space_guid := space.Metadata.Guid

		service_instances_json := cf.GetServiceInstancesForSpace(space_guid)
		service_instances := &ServiceInstances{}
		json.Unmarshal([]byte(service_instances_json), service_instances)

		if len(service_instances.Services) > 0 {
			organization := fetchOrganization(cf, space.Entity.OrganizationGuid)
			space_name := space.Entity.Name
			organization_name := organization.Entity.Name
			space_dir := spaceDirectory(organization_name, space_name)
			os.MkdirAll(space_dir, 0777)

			for _, service_instance := range service_instances.Services {
				if service_instance.ServicePlan.Service.Label == "p-riakcs" {
					service_instance_guid := service_instance.Guid
					service_instance_name := service_instance.Name
					instance_dir := space_dir + "/service_instances/" + service_instance_name
					os.MkdirAll(instance_dir, 0777)
					writeMetadataFile(cf, organization_name, space_name, service_instance_name, service_instance_guid)

					data_dir := instance_dir + "/data"
					os.MkdirAll(data_dir, 0777)

					bucket_name := bucketNameFromServiceInstanceGuid(service_instance_guid)
					s3cmd.FetchBucket(bucket_name, data_dir)
				}
			}
		}
	}
}

func bucketNameFromServiceInstanceGuid(service_instance_guid string) string {
	return "service-instance-" + service_instance_guid
}

func fetchOrganization(cf CfClientInterface, organization_guid string) Organization {
	organization_json := cf.GetOrganization(organization_guid)
	organization := &Organization{}
	json.Unmarshal([]byte(organization_json), organization)
	return *organization
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

func writeMetadataFile(cf CfClientInterface, organization_name string, space_name string, service_instance_name string, service_instance_guid string) {
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

	space_dir := spaceDirectory(organization_name, space_name)
	path := fmt.Sprintf("%s/service_instances/%s/metadata.yml", space_dir, service_instance_name)
	ioutil.WriteFile(path, bytes, 0777)
}

func spaceDirectory(organization_name, space_name string) string {
	return fmt.Sprintf("/tmp/backup/orgs/%s/spaces/%s", organization_name, space_name)
}
