package riak_backup

type CfClientInterface interface {
	GetSpaces(next_url string) string
	GetServiceInstancesForSpace(space_guid string) string
	GetBindings(service_instance_guid string) string
}

type CfClient struct {
}

func(cf *CfClient) GetSpaces(next_url string) string {
	return ""
}

func (cf *CfClient) GetServiceInstancesForSpace(space_guid string) string {
	return ""
}

func (cf *CfClient) GetBindings(service_instance_guid string) string {
	return ""
}
