package riak_backup

type CfClientInterface interface {
	GetSpaces() string
	GetServiceInstancesForSpace(space_guid string) string
}

type CfClient struct {
}

func(cf *CfClient) GetSpaces() string {
	return ""
}

func (cf *CfClient) GetServiceInstancesForSpace(space_guid string) string {
	return ""
}
