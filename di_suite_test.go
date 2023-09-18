package di_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2" //nolint:depguard
	. "github.com/onsi/gomega"    //nolint:depguard
)

func TestDi(t *testing.T) { //nolint:paralleltest
	RegisterFailHandler(Fail)
	RunSpecs(t, "Di Suite")
}
