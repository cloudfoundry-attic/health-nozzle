package counter_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCounter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Counter Suite")
}
