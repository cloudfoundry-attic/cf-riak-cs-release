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
	It("Makes a directory for each space", func() {
		Backup(&test_support.FakeCfClient{})

		directories, _ := ioutil.ReadDir("/tmp/backup/spaces")
		Expect(directories).To(HaveLen(3))
		Expect(directories[0].IsDir()).To(BeTrue())
		Expect(directories[1].IsDir()).To(BeTrue())
		Expect(directories[2].IsDir()).To(BeTrue())

		guids := []string{ directories[0].Name(), directories[1].Name(), directories[2].Name() }
		Expect(guids).To(ContainElement("space-0"))
		Expect(guids).To(ContainElement("space-1"))
		Expect(guids).To(ContainElement("space-2"))
	})

	It("Makes a sub-directory for each riak-cs service instance in each space", func() {
		Backup(&test_support.FakeCfClient{})

		directories, _ := ioutil.ReadDir("/tmp/backup/spaces/space-0/service_instances")
		Expect(directories).To(HaveLen(2))
		Expect(directories[0].IsDir()).To(BeTrue())
		Expect(directories[1].IsDir()).To(BeTrue())

		guids := []string{ directories[0].Name(), directories[1].Name() }
		Expect(guids).To(ContainElement("service-instance-0"))
		Expect(guids).To(ContainElement("service-instance-1"))

		directories, _ = ioutil.ReadDir("/tmp/backup/spaces/space-1/service_instances")
		Expect(directories).To(HaveLen(2))
		Expect(directories[0].IsDir()).To(BeTrue())
		Expect(directories[1].IsDir()).To(BeTrue())

		guids = []string{ directories[0].Name(), directories[1].Name() }
		Expect(guids).To(ContainElement("service-instance-2"))
		Expect(guids).To(ContainElement("service-instance-3"))
		Expect(guids).NotTo(ContainElement("non-riak-service-instance"))
	})

	It("saves the instance name and list of bound apps in a metadata file for each instance", func() {
		Backup(&test_support.FakeCfClient{})

		entries, _ := ioutil.ReadDir("/tmp/backup/spaces/space-0/service_instances/service-instance-0")
		Expect(entries).To(HaveLen(1))
		Expect(entries[0].IsDir()).To(BeFalse())
		Expect(entries[0].Name()).To(Equal("metadata.yml"))

		// instance with a bound apps
		file_path := "/tmp/backup/spaces/space-0/service_instances/service-instance-0/metadata.yml"
		metadata := NewFromFilename(file_path)
		Expect(metadata.ServiceInstanceGuid).To(Equal("service-instance-0"))
		Expect(metadata.BoundApps).To(HaveLen(2))
		Expect(metadata.BoundApps[0].Guid).To(Equal("app-guid-0"))
		Expect(metadata.BoundApps[0].Name).To(Equal("app-name-0"))
		Expect(metadata.BoundApps[1].Guid).To(Equal("app-guid-1"))
		Expect(metadata.BoundApps[1].Name).To(Equal("app-name-1"))

		// instance with no bound apps
		file_path = "/tmp/backup/spaces/space-0/service_instances/service-instance-1/metadata.yml"
		metadata = NewFromFilename(file_path)
		Expect(metadata.ServiceInstanceGuid).To(Equal("service-instance-1"))
		Expect(metadata.BoundApps).To(HaveLen(0))
	})

	AfterEach(func(){
		err := os.RemoveAll("/tmp/backup")
		if err != nil {
			fmt.Println(err.Error())
		}
	})
})
