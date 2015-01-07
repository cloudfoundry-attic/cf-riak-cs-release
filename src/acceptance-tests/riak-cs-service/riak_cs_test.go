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
	. "github.com/cloudfoundry-incubator/cf-test-helpers/services/context_setup"
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

		Eventually(Cf("create-service", ServiceName, PlanName, ServiceInstanceName), ScaledTimeout(60*time.Second)).Should(Exit(0))
		Eventually(Cf("bind-service", AppName, ServiceInstanceName), ScaledTimeout(60*time.Second)).Should(Exit(0))
		Eventually(Cf("start", AppName), ScaledTimeout(5*time.Minute)).Should(Exit(0))

		uri := AppUri(AppName) + "/service/blobstore/" + ServiceInstanceName + "/mykey"
		delete_uri := AppUri(AppName) + "/service/blobstore/" + ServiceInstanceName

		fmt.Println("Posting to url: ", uri)
		Eventually(Curl("-k", "-d", "myvalue", uri), ScaledTimeout(10*time.Second), 1.0).Should(Say("myvalue"))
		fmt.Println("\n")

		fmt.Println("Curling url: ", uri)
		Eventually(Curl("-k", uri), ScaledTimeout(10*time.Second), 1.0).Should(Say("myvalue"))
		fmt.Println("\n")

		fmt.Println("Sending delete to: ", delete_uri)
		Eventually(Curl("-X", "DELETE", "-k", delete_uri), ScaledTimeout(10*time.Second), 1.0).Should(Say("successfully_deleted"))
		fmt.Println("\n")

		Eventually(Cf("unbind-service", AppName, ServiceInstanceName), ScaledTimeout(60*time.Second)).Should(Exit(0))
		Eventually(Cf("delete-service", "-f", ServiceInstanceName), ScaledTimeout(60*time.Second)).Should(Exit(0))
	})
})
