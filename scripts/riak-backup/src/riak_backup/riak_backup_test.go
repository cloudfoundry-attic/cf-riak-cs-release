package riak_backup_test

import (
	. "riak_backup"
	"riak_backup/test_support"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"fmt"
)

var _ = Describe("RiakBackup", func() {
	It("Makes a directory for each space that has instances", func() {
		Backup(&test_support.FakeCfClient{}, &test_support.FakeS3CmdClient{}, "/tmp/backup")

		directories, _ := ioutil.ReadDir("/tmp/backup/orgs/organization-name-0/spaces")
		Expect(directories).To(HaveLen(2))
		Expect(directories[0].IsDir()).To(BeTrue())
		Expect(directories[1].IsDir()).To(BeTrue())

		space_dir_names := []string{ directories[0].Name(), directories[1].Name() }
		Expect(space_dir_names).To(ContainElement("space-name-0"))
		Expect(space_dir_names).To(ContainElement("space-name-2"))
	})

	It("Makes a sub-directory for each riak-cs service instance in each space", func() {
		Backup(&test_support.FakeCfClient{}, &test_support.FakeS3CmdClient{}, "/tmp/backup")

		directories, _ := ioutil.ReadDir("/tmp/backup/orgs/organization-name-0/spaces/space-name-0/service_instances")
		Expect(directories).To(HaveLen(2))
		Expect(directories[0].IsDir()).To(BeTrue())
		Expect(directories[1].IsDir()).To(BeTrue())

		instance_names := []string{ directories[0].Name(), directories[1].Name() }
		Expect(instance_names).To(ContainElement("service-instance-name-0"))
		Expect(instance_names).To(ContainElement("service-instance-name-1"))

		directories, _ = ioutil.ReadDir("/tmp/backup/orgs/organization-name-0/spaces/space-name-2/service_instances")
		Expect(directories).To(HaveLen(2))
		Expect(directories[0].IsDir()).To(BeTrue())
		Expect(directories[1].IsDir()).To(BeTrue())

		instance_names = []string{ directories[0].Name(), directories[1].Name() }
		Expect(instance_names).To(ContainElement("service-instance-name-2"))
		Expect(instance_names).To(ContainElement("service-instance-name-3"))
		Expect(instance_names).NotTo(ContainElement("non-riak-service-instance-name"))
	})

	It("saves the instance guid and list of bound apps in a metadata file for each instance", func() {
		Backup(&test_support.FakeCfClient{}, &test_support.FakeS3CmdClient{}, "/tmp/backup")

		entries, _ := ioutil.ReadDir("/tmp/backup/orgs/organization-name-0/spaces/space-name-0/service_instances/service-instance-name-0")

		Expect(entries).To(HaveLen(2))
		Expect(entries[1].IsDir()).To(BeFalse())
		Expect(entries[1].Name()).To(Equal("metadata.yml"))

		// instance with a bound apps
		file_path := "/tmp/backup/orgs/organization-name-0/spaces/space-name-0/service_instances/service-instance-name-0/metadata.yml"
		metadata := NewMetadataFromFilename(file_path)
		Expect(metadata.ServiceInstanceGuid).To(Equal("service-instance-guid-0"))
		Expect(metadata.BoundApps).To(HaveLen(2))
		Expect(metadata.BoundApps[0].Guid).To(Equal("app-guid-0"))
		Expect(metadata.BoundApps[0].Name).To(Equal("app-name-0"))
		Expect(metadata.BoundApps[1].Guid).To(Equal("app-guid-1"))
		Expect(metadata.BoundApps[1].Name).To(Equal("app-name-1"))

		// instance with no bound apps
		file_path = "/tmp/backup/orgs/organization-name-0/spaces/space-name-0/service_instances/service-instance-name-1/metadata.yml"
		metadata = NewMetadataFromFilename(file_path)
		Expect(metadata.ServiceInstanceGuid).To(Equal("service-instance-guid-1"))
		Expect(metadata.BoundApps).To(HaveLen(0))
	})

	It("saves the data from each instance into the data directory", func() {
		Backup(&test_support.FakeCfClient{}, &test_support.FakeS3CmdClient{}, "/tmp/backup")

		entries, _ := ioutil.ReadDir("/tmp/backup/orgs/organization-name-0/spaces/space-name-0/service_instances/service-instance-name-0")
		Expect(entries).To(HaveLen(2))
		Expect(entries[0].IsDir()).To(BeTrue())
		Expect(entries[0].Name()).To(Equal("data"))

		data_path := "/tmp/backup/orgs/organization-name-0/spaces/space-name-0/service_instances/service-instance-name-0/data"
		data_items, _ := ioutil.ReadDir(data_path)
		Expect(data_items).To(HaveLen(1))
		Expect(data_items[0].IsDir()).To(BeFalse())
		Expect(data_items[0].Name()).To(Equal("datafile.dat"))

		bytes, _ := ioutil.ReadFile(data_path + "/" + data_items[0].Name())
		Expect(string(bytes)).To(Equal("data from bucket service-instance-service-instance-guid-0"))
	})

	AfterEach(func(){
		err := os.RemoveAll("/tmp/backup")
		if err != nil {
			fmt.Println(err.Error())
		}
	})
})
