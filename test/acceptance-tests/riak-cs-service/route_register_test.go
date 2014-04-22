package riak_cs_service

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/vito/cmdtest/matchers"
	"github.com/vito/cmdtest"
	"github.com/pivotal-cf-experimental/cf-test-helpers/runner"
	"fmt"
)

var _ = Describe("Riak CS Nodes Register a Route", func() {

		It("Allows users to access the riak-cs service using external url instead of IP of single machine after register the route", func() {
			endpointURL := IntegrationConfig.RiakCsScheme + "riakcs." + IntegrationConfig.AppsDomain + "/riak-cs/ping"
			fmt.Println("Endpoint URL: " + endpointURL)

			var session *cmdtest.Session
			session = runner.Curl("-k", endpointURL)

			Expect(session).To(Say("OK"))
		})
	})
