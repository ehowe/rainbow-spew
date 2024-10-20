package spew_test

import (
	"testing"

	spew "github.com/ehowe/rainbow-spew"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSpew(t *testing.T) {
	RegisterFailHandler(Fail)
	spew.Config = *spew.NewTestConfig()
	RunSpecs(t, "Spew Suite")
}
