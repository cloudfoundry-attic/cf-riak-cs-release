package helpers

import (
	"fmt"
	"time"

	ginkgoconfig "github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/cf-test-helpers/cf"
	"github.com/vito/cmdtest"
	. "github.com/vito/cmdtest/matchers"
)

type ConfiguredContext struct {
	config IntegrationConfig

	organizationName string
	spaceName        string

	regularUserUsername string
	regularUserPassword string

	isPersistent bool
}

func NewContext(config IntegrationConfig) *ConfiguredContext {
	node := ginkgoconfig.GinkgoConfig.ParallelNode
	timeTag := time.Now().Format("2006_01_02-15h04m05.999s")

	return &ConfiguredContext{
		config: config,

		organizationName: fmt.Sprintf("CATS-ORG-%d-%s", node, timeTag),
		spaceName:        fmt.Sprintf("CATS-SPACE-%d-%s", node, timeTag),

		regularUserUsername: fmt.Sprintf("CATS-USER-%d-%s", node, timeTag),
		regularUserPassword: "meow",

		isPersistent: false,
	}
}

func (context *ConfiguredContext) Setup() {
	cf.AsUser(context.AdminUserContext(), func() {
		Expect(cf.Cf("create-user", context.regularUserUsername, context.regularUserPassword)).To(SayBranches(
			cmdtest.ExpectBranch{"OK", func() {}},
			cmdtest.ExpectBranch{"scim_resource_already_exists", func() {}},
		))

		Expect(cf.Cf("create-org", context.organizationName)).To(ExitWith(0))
	})
}

func (context *ConfiguredContext) Teardown() {
	cf.AsUser(context.AdminUserContext(), func() {
		Expect(cf.Cf("delete-user", "-f", context.regularUserUsername)).To(Say("OK"))

		if !context.isPersistent {
			Expect(cf.Cf("delete-org", "-f", context.organizationName)).To(Say("OK"))
		}
	})
}

func (context *ConfiguredContext) AdminUserContext() cf.UserContext {
	return cf.NewUserContext(
		context.config.ApiEndpoint,
		context.config.AdminUser,
		context.config.AdminPassword,
		"",
		"",
		context.config.SkipSSLValidation,
	)
}

func (context *ConfiguredContext) RegularUserContext() cf.UserContext {
	return cf.NewUserContext(
		context.config.ApiEndpoint,
		context.regularUserUsername,
		context.regularUserPassword,
		context.organizationName,
		context.spaceName,
		context.config.SkipSSLValidation,
	)
}
