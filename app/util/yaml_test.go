package util_test

import (
	"github.com/dcwangmit01/goapi/app/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Util", func() {

	Context("Yaml", func() {

		Context("StructToYamlStr and StructFromYamlStr", func() {

			It("Should both match", func() {
				var err error
				var p1, p2 *Parent
				var s1, s2 string

				p1 = NewGoodParent()
				s1, err = util.StructToYamlStr(p1)
				Expect(err).Should(BeNil())

				p2 = &Parent{}
				err = util.StructFromYamlStr(p2, s1)
				Expect(err).Should(BeNil())

				s2, err = util.StructToYamlStr(p2)
				Expect(err).Should(BeNil())
				Expect(s1).Should(Equal(s2))
			})
		})
	})
})
