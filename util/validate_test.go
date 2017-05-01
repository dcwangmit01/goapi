package util_test

import (
	"github.com/dcwangmit01/goapi/util"
	validator "gopkg.in/go-playground/validator.v9"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
)

type Parent struct {
	Settings *Settings `validate:"dive"`
	Children []*Child  `validate:"dive"`
}

type Settings struct {
	Enum string `validate:"eq=ABC|eq=DEF"`
}

type Child struct {
	Username string `validate:"required,email|eq=admin"`
	Phone    string `validate:"phone,min=7"`
}

func NewGoodParent() *Parent {
	return &Parent{
		&Settings{
			Enum: "ABC",
		},
		[]*Child{
			&Child{ // default admin user
				Username: "admin",
				Phone:    `+155w5*5#5555ww55`,
			},
			&Child{ // a regular user
				Username: "user@domain.com",
				Phone:    `+1234567890ww2134w#0`,
			},
		},
	}
}

func NewBadParent() *Parent {
	return &Parent{
		&Settings{
			Enum: "JKL",
		},
		[]*Child{
			&Child{ // default admin user
				Username: "admin2",
				Phone:    "+12345z67890w*#w",
			},
			&Child{ // a regular user
				Username: "user@domaincom",
				Phone:    "wc38242384343",
			},
		},
	}
}

var _ = Describe("Util", func() {

	Context("Validate", func() {

		It("Should pass on good values", func() {
			parent := NewGoodParent()
			err := util.ValidateStruct(parent)
			if err != nil { // output to help debug
				errs := err.(validator.ValidationErrors)
				fmt.Printf("%+v\n", errs)
			}
			Expect(err).Should(BeNil())
		})

		It("Should fail on bad values", func() {
			parent := NewBadParent()
			err := util.ValidateStruct(parent)
			//errs := err.(validator.ValidationErrors)
			//fmt.Printf("%+v\n", errs)

			merrs := util.ValidationErrorToMap(err)

			Expect(len(merrs)).Should(Equal(5))
		})
	})
})
