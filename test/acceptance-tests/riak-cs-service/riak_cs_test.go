package riak_cs_service

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"

	"fmt"
	. "github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	. "github.com/cloudfoundry-incubator/cf-test-helpers/generator"
	. "github.com/cloudfoundry-incubator/cf-test-helpers/runner"
	"time"
)

var _ = Describe("Riak CS Service Lifecycle", func() {
	BeforeEach(func() {
		AppName = RandomName()

		Eventually(Cf("push", AppName, "-m", "256M", "-p", sinatraPath, "-no-start"), 60*time.Second).Should(Exit(0))
	})

	AfterEach(func() {
		Eventually(Cf("delete", AppName, "-f"), 60*time.Second).Should(Exit(0))
	})

	It("Allows users to create, bind, write to, read from, unbind, and destroy the service instance", func() {
		ServiceName := ServiceName()
		PlanName := PlanName()
		ServiceInstanceName := RandomName()

		Eventually(Cf("create-service", ServiceName, PlanName, ServiceInstanceName), 60*time.Second).Should(Exit(0))
		Eventually(Cf("bind-service", AppName, ServiceInstanceName), 60*time.Second).Should(Exit(0))
		Eventually(Cf("start", AppName), 5*60*time.Second).Should(Exit(0))

		uri := AppUri(AppName) + "/service/blobstore/" + ServiceInstanceName + "/mykey"
		delete_uri := AppUri(AppName) + "/service/blobstore/" + ServiceInstanceName

		fmt.Println("Posting to url: ", uri)
		Eventually(Curl("-k", "-d", "myvalue", uri), 10.0, 1.0).Should(Say("myvalue"))
		fmt.Println("\n")

		fmt.Println("Curling url: ", uri)
		Eventually(Curl("-k", uri), 10.0, 1.0).Should(Say("myvalue"))
		fmt.Println("\n")

		fmt.Println("Sending delete to: ", delete_uri)
		Eventually(Curl("-X", "DELETE", "-k", delete_uri), 10.0, 1.0).Should(Say("successfully_deleted"))
		fmt.Println("\n")

		Eventually(Cf("unbind-service", AppName, ServiceInstanceName), 60*time.Second).Should(Exit(0))
		Eventually(Cf("delete-service", "-f", ServiceInstanceName), 60*time.Second).Should(Exit(0))
	})
})
