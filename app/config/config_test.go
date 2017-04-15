package config_test

import (
	"github.com/dcwangmit01/grpc-gw-poc/app/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
)

var _ = Describe("Config", func() {

	Describe("AppConfig", func() {
		Context("When created via NewAppConfig", func() {
			ac := config.NewAppConfig()

			It("Should have default settings", func() {
				Expect(ac.Settings.LogLevel).Should(Equal("DEBUG"))
			})

			It("Should have a default admin user", func() {
				Expect(len(ac.Users)).Should(Equal(1))
				Expect(ac.Users[0].Name).Should(Equal(config.DefaultAdminUser))
				Expect(len(ac.Users[0].Id)).Should(BeNumerically(">", 0))
				Expect(ac.Users[0].Role).Should(Equal("ADMIN"))
			})

			It("Should have a valid admin password hash", func() {
				err := ac.Users[0].ValidatePassword(config.DefaultAdminPass)
				Expect(err).Should(BeNil())

				err = ac.Users[0].ValidatePassword("not the right password")
				Expect(err).ShouldNot(BeNil())
			})
		})

		Context("AddUser", func() {

			It("Should increase the user count", func() {
				ac := config.NewAppConfig()

				Expect(len(ac.Users)).Should(Equal(1))

				u := config.NewUser()
				ac.AddUser(u)

				Expect(len(ac.Users)).Should(Equal(2))
			})
		})

		Context("Dump and Parse", func() {

			var err error
			var ac1, ac2 *config.AppConfig
			var ac1Str, ac2Str string

			ac1 = config.NewAppConfig()
			ac1.Users = append(ac1.Users, config.NewUser())
			ac1.Users[0].Email = "user1@gmail.com"
			ac1.Users[1].Email = "asdf"

			It("Should both match", func() {
				ac1Str, err = ac1.Dump()
				Expect(err).Should(BeNil())

				ac2, err = config.ParseAppConfig(ac1Str)
				Expect(err).Should(BeNil())

				ac2Str, err = ac2.Dump()
				Expect(err).Should(BeNil())

				Expect(ac1Str).Should(Equal(ac2Str))
			})
		})
	})

	Describe("User", func() {

		Context("When created via NewUser", func() {
			u := config.NewUser()

			It("Should have an ID", func() {
				Expect(len(u.Id)).Should(BeNumerically(">", 0))
			})

			It("Should have a Role", func() {
				Expect(u.Role).Should(Equal("USER"))
			})

			It("Should have an Empty Password", func() {
				Expect(u.PasswordHash).Should(BeZero())
			})
		})

		Context("HashPassword and ValidatePassword", func() {
			u := config.NewUser()
			myPassword := "1234asdf!@#$"

			It("Should create a hash of the password and validate", func() {
				var err error

				err = u.HashPassword(myPassword)
				Expect(myPassword).ShouldNot(Equal(u.PasswordHash))
				Expect(err).Should(BeNil())

				err = u.ValidatePassword(myPassword)
				Expect(err).Should(BeNil())

				err = u.ValidatePassword("not the right password")
				Expect(err).ShouldNot(BeNil())
			})
		})

	})

	Describe("Static Functions", func() {

		Context("Validate", func() {

			It("Should pass on good values", func() {
				myPassword := "1234asdf!@#$"
				u := config.NewUser()
				u.Email = "test@test.com"
				u.Name = "First Last"
				u.HashPassword(myPassword)
				u.Role = "USER"
				u.Phone = "012345678901234"
				errs := config.ValidateStruct(u)

				if len(errs) > 0 { // output to help debug
					fmt.Printf("%+v", errs)
				}
				Expect(len(errs)).Should(Equal(0))
			})

			It("Should fail on bad values", func() {
				u := config.NewUser()
				errs := config.ValidateStruct(u)

				Expect(len(errs)).Should(Equal(4))
				Expect(errs).Should(HaveKey("User.Email"))
				Expect(errs["User.Email"]).Should(Equal("required"))
				Expect(errs).Should(HaveKey("User.Name"))
				Expect(errs["User.Name"]).Should(Equal("required"))
				Expect(errs).Should(HaveKey("User.PasswordHash"))
				Expect(errs["User.PasswordHash"]).Should(Equal("required"))
				Expect(errs).Should(HaveKey("User.Phone"))
				Expect(errs["User.Phone"]).Should(Equal("phone"))
			})
		})
	})
})
