package config_test

import (
	"github.com/dcwangmit01/goapi/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {

	Describe("AppConfig", func() {
		Context("When created via NewAppConfig", func() {
			ac := config.NewAppConfig()

			It("Should have default settings", func() {
				Expect(ac.Settings.LogLevel).Should(Equal("debug"))
			})

			It("Should have a default admin user", func() {
				Expect(len(ac.Users)).Should(Equal(1))
				Expect(ac.Users[0].Name).Should(Equal(config.DefaultAdminUsername))
				Expect(len(ac.Users[0].Id)).Should(BeNumerically(">", 0))
				Expect(ac.Users[0].Role).Should(Equal("admin"))
			})

			It("Should have a valid admin password hash", func() {
				err := ac.Users[0].ValidatePassword(config.DefaultAdminPassword)
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

		Context("GetUserByUsername", func() {
			ac := config.NewAppConfig()
			ac.Users = append(ac.Users, config.NewUser())
			ac.Users[1].Username = "user@domain.com"

			It("Should return an admin user", func() {
				u, err := ac.GetUserByUsername("admin")
				Expect(u).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
			})
			It("Should return a user, user", func() {
				u, err := ac.GetUserByUsername("user@domain.com")
				Expect(u).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
			})
			It("Should return nil when a user does not exist", func() {
				u, err := ac.GetUserByUsername("nonexistant@domain.com")
				Expect(u).Should(BeNil())
				Expect(err).ShouldNot(BeNil())
			})

		})

		Context("ToYaml and FromYaml", func() {

			var err error
			var ac1, ac2 *config.AppConfig
			var ac1Str, ac2Str string

			ac1 = config.NewAppConfig()
			ac1.Users = append(ac1.Users, config.NewUser())
			ac1.Users[0].Username = "user1@domain.com"
			ac1.Users[1].Username = "asdf"

			It("Should both match", func() {
				ac1Str, err = ac1.ToYaml()
				Expect(err).Should(BeNil())

				ac2, err = config.AppConfigFromYaml(ac1Str)
				Expect(err).Should(BeNil())

				ac2Str, err = ac2.ToYaml()
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
				Expect(u.Role).Should(Equal("user"))
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
})
