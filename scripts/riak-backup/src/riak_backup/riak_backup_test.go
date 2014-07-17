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
		Expect(directories).To(HaveLen(2))
		Expect(directories[0].IsDir()).To(BeTrue())
		Expect(directories[1].IsDir()).To(BeTrue())

		guids := []string{ directories[0].Name(), directories[1].Name() }
		Expect(guids).To(ContainElement("space-0"))
		Expect(guids).To(ContainElement("space-1"))
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
	})

	AfterEach(func(){
		err := os.RemoveAll("/tmp/backup")
		if err != nil {
			fmt.Println(err.Error())
		}
	})
})
