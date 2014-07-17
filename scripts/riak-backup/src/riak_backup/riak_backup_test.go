package riak_backup_test

import (
	. "riak_backup"
	"riak_backup/test_support"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
)

var _ = Describe("RiakBackup", func() {
	It("Makes a directory for each space", func() {
		Backup(&test_support.FakeCfClient{})

		directories, _ := ioutil.ReadDir("/tmp/backup/spaces")
		Expect(directories).To(HaveLen(2))
		Expect(directories[0].IsDir()).To(BeTrue())
		Expect(directories[1].IsDir()).To(BeTrue())

		names := []string{ directories[0].Name(), directories[1].Name() }
		Expect(names).To(ContainElement("413c4df3-66b6-4a7e-a681-bd7f89cffcd9"))
		Expect(names).To(ContainElement("0d8ff79a-b9b1-4d84-9de1-015a6c884269"))
	})

	AfterEach(func(){
		os.Remove("/tmp/backup")
	})
})
