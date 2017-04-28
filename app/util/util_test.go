package util_test

import (
	"github.com/dcwangmit01/goapi/app/config"
	"github.com/dcwangmit01/goapi/app/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
)

var _ = Describe("Util", func() {

	Context("Validate", func() {

		It("Should pass on good values", func() {
			myPassword := "1234asdf!@#$"
			u := config.NewUser()
			u.Username = "user@domain.com"
			u.Name = "First Last"
			u.HashPassword(myPassword)
			u.Role = "user"
			u.Phone = "012345678901234"
			errs := util.ValidateStruct(u)

			if len(errs) > 0 { // output to help debug
				fmt.Printf("%+v", errs)
			}
			Expect(len(errs)).Should(Equal(0))
		})

		It("Should fail on bad values", func() {
			u := config.NewUser()
			errs := util.ValidateStruct(u)

			Expect(len(errs)).Should(Equal(4))
			Expect(errs).Should(HaveKey("User.Username"))
			Expect(errs["User.Username"]).Should(Equal("required"))
			Expect(errs).Should(HaveKey("User.Name"))
			Expect(errs["User.Name"]).Should(Equal("required"))
			Expect(errs).Should(HaveKey("User.PasswordHash"))
			Expect(errs["User.PasswordHash"]).Should(Equal("required"))
			Expect(errs).Should(HaveKey("User.Phone"))
			Expect(errs["User.Phone"]).Should(Equal("phone"))
		})

		It("User.username should pass on special username value of 'admin'", func() {
			myPassword := "1234asdf!@#$"
			u := config.NewUser()
			u.Username = "admin" // <- special value
			u.Name = "First Last"
			u.HashPassword(myPassword)
			u.Role = "user"
			u.Phone = "012345678901234"
			errs := util.ValidateStruct(u)

			if len(errs) > 0 { // output to help debug
				fmt.Printf("%+v", errs)
			}
			Expect(len(errs)).Should(Equal(0))
		})
	})
})
