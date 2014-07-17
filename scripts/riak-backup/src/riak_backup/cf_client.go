package riak_backup

type CfClientInterface interface {
	GetSpaces() string
	GetServiceInstancesForSpace(space_guid string) string
	GetBindings(service_instance_guid string) string
}

type CfClient struct {
}

func(cf *CfClient) GetSpaces() string {
	return ""
}

func (cf *CfClient) GetServiceInstancesForSpace(space_guid string) string {
	return ""
}

func (cf *CfClient) GetBindings(service_instance_guid string) string {
	return ""
}
