package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestHealthNozzle(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HealthNozzle Suite")
}
