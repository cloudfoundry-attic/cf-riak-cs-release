package route_register_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "test_helpers"
	"testing"
	"fmt"
)

func TestRoute_register(t *testing.T) {
	RegisterFailHandler(Fail)

	fmt.Println("starting gnatsd...")
	natsCmd := StartNats(4222)

	RunSpecs(t, "Route_register Suite")

	StopCmd(natsCmd)
}
