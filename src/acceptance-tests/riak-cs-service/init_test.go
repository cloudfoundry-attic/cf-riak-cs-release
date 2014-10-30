package riak_cs_service

import (
	. "github.com/cloudfoundry-incubator/cf-test-helpers/runner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"../helpers"
	context_setup "github.com/cloudfoundry-incubator/cf-test-helpers/services/context_setup"
	"testing"
)

func TestServices(t *testing.T) {
	context_setup.TimeoutScale = RiakCSIntegrationConfig.TimeoutScale

	context_setup.SetupEnvironment(context_setup.NewContext(RiakCSIntegrationConfig.IntegrationConfig, "RiakCSATS"))

	RegisterFailHandler(Fail)
	RunSpecs(t, "Riak CS Services Suite")
}

func AppUri(appname string) string {
	return RiakCSIntegrationConfig.RiakCsScheme + appname + "." + RiakCSIntegrationConfig.AppsDomain
}

func Curling(args ...string) func() *gexec.Session {
	return func() *gexec.Session {
		return Curl(args...)
	}
}

func ServiceName() string {
	return RiakCSIntegrationConfig.ServiceName
}

func PlanName() string {
	return RiakCSIntegrationConfig.PlanName
}

var RiakCSIntegrationConfig = helpers.LoadConfig()

var AppName = ""

var sinatraPath = "../assets/app_sinatra_service"
