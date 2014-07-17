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

func Backup(cf CfClientInterface) {
	spaces_json := cf.GetSpaces()
	spaces := &Spaces{}
	json.Unmarshal([]byte(spaces_json), spaces)

	for _, space := range spaces.Resources {
		os.MkdirAll(fmt.Sprintf("/tmp/backup/spaces/%s", space.Metadata.Guid), 0777)
	}
}
