package riak_cs_service

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf-experimental/cf-test-helpers/runner"
	"github.com/vito/cmdtest"

	"../helpers"
	"testing"
)

func TestServices(t *testing.T) {
	helpers.SetupEnvironment(helpers.NewContext(IntegrationConfig))

	RegisterFailHandler(Fail)
	RunSpecs(t, "Riak CS Services Suite")
}

func AppUri(appname string) string {
	return IntegrationConfig.RiakCsScheme + appname + "." + IntegrationConfig.AppsDomain
}

func Curling(args ...string) func() *cmdtest.Session {
	return func() *cmdtest.Session {
		return Curl(args...)
	}
}

func ServiceName() string {
	return IntegrationConfig.ServiceName
}

func PlanName() string {
	return IntegrationConfig.PlanName
}

var IntegrationConfig = helpers.LoadConfig()

var AppName = ""

var sinatraPath = "../assets/app_sinatra_service"
