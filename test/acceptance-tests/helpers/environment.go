package helpers

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	"time"
)

var AdminUserContext cf.UserContext
var RegularUserContext cf.UserContext

type SuiteContext interface {
	Setup()
	Teardown()

	AdminUserContext() cf.UserContext
	RegularUserContext() cf.UserContext
}

func SetupEnvironment(context SuiteContext) {
	var originalCfHomeDir, currentCfHomeDir string

	BeforeEach(func() {
		AdminUserContext = context.AdminUserContext()
		RegularUserContext = context.RegularUserContext()

		context.Setup()

		cf.AsUser(AdminUserContext, func() {
			setUpSpaceWithUserAccess(RegularUserContext)
		})

		originalCfHomeDir, currentCfHomeDir = cf.InitiateUserContext(RegularUserContext)
		cf.TargetSpace(RegularUserContext)
	})

	AfterEach(func() {
		cf.RestoreUserContext(RegularUserContext, originalCfHomeDir, currentCfHomeDir)

		context.Teardown()
	})
}

func setUpSpaceWithUserAccess(uc cf.UserContext) {
	Eventually(cf.Cf("create-space", "-o", uc.Org, uc.Space), 30 * time.Second).Should(Exit(0))
	Eventually(cf.Cf("set-space-role", uc.Username, uc.Org, uc.Space, "SpaceManager"), 30 * time.Second).Should(Exit(0))
	Eventually(cf.Cf("set-space-role", uc.Username, uc.Org, uc.Space, "SpaceDeveloper"), 30 * time.Second).Should(Exit(0))
	Eventually(cf.Cf("set-space-role", uc.Username, uc.Org, uc.Space, "SpaceAuditor"), 30 * time.Second).Should(Exit(0))
}
