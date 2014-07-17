package riak_backup

type CfClientInterface interface {
	GetSpaces() string
}

type CfClient struct {
}

func(cf *CfClient) GetSpaces() string {
	return ""
}
