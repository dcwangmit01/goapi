package sqlitekv_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSqlitekv(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sqlitekv Suite")
}
