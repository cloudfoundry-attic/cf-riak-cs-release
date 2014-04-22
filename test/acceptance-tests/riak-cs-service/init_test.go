package riak_cs_service

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vito/cmdtest"
	. "github.com/pivotal-cf-experimental/cf-test-helpers/runner"

	"testing"
	"../helpers"
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

var IntegrationConfig = helpers.LoadConfig()

var AppName = ""

var sinatraPath = "../assets/app_sinatra_service"
