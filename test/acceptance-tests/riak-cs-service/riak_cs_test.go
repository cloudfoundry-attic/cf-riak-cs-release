package riak_cs_service

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/vito/cmdtest/matchers"

	. "github.com/pivotal-cf-experimental/cf-test-helpers/cf"
	. "github.com/pivotal-cf-experimental/cf-test-helpers/generator"
	"fmt"
	"time"
)

var _ = Describe("Riak CS Service Lifecycle", func() {
	BeforeEach(func() {
		AppName = RandomName()

		Expect(Cf("push", AppName, "-m", "256M", "-p", sinatraPath, "-no-start")).To(ExitWithTimeout(0, 60*time.Second))
	})

	AfterEach(func() {
		Expect(Cf("delete", AppName, "-f")).To(ExitWithTimeout(0, 20*time.Second))
	})

	It("Allows users to create, bind, write to, read from, unbind, and destroy the service instance", func() {
		ServiceName := "riak-cs"
		PlanName := "bucket"
		ServiceInstanceName := RandomName()

		Expect(Cf("create-service", ServiceName, PlanName, ServiceInstanceName)).To(ExitWithTimeout(0, 60*time.Second))
		Expect(Cf("bind-service", AppName, ServiceInstanceName)).To(ExitWithTimeout(0, 60*time.Second))
		Expect(Cf("start", AppName)).To(ExitWithTimeout(0, 5*60*time.Second))

		uri := AppUri(AppName) + "/service/blobstore/" + ServiceInstanceName + "/mykey"
		delete_uri := AppUri(AppName) + "/service/blobstore/" + ServiceInstanceName

		fmt.Println("Posting to url: ", uri)
		Eventually(Curling("-k", "-d", "myvalue", uri), 10.0, 1.0).Should(Say("myvalue"))
		fmt.Println("\n")

		fmt.Println("Curling url: ", uri)
		Eventually(Curling("-k", uri), 10.0, 1.0).Should(Say("myvalue"))
		fmt.Println("\n")

		fmt.Println("Sending delete to: ", delete_uri)
		Eventually(Curling("-X", "DELETE", "-k", delete_uri), 10.0, 1.0).Should(Say(""))
		fmt.Println("\n")

		Expect(Cf("unbind-service", AppName, ServiceInstanceName)).To(ExitWithTimeout(0, 20*time.Second))
		Expect(Cf("delete-service", "-f", ServiceInstanceName)).To(ExitWithTimeout(0, 20*time.Second))
	})
})
