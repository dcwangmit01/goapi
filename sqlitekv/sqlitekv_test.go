package sqlitekv_test

import (
	"io/ioutil"
	"log"

	"github.com/dcwangmit01/goapi/sqlitekv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"
)

var _ = Describe("Sqlitekv", func() {

	Context("New", func() {

		// create a tmp directory for the db file
		dir, err := ioutil.TempDir("", "tmp_dir")
		if err != nil {
			log.Fatal(err)
		}
		defer os.RemoveAll(dir) // clean up

		// create the sqlitekv database
		skv := sqlitekv.New(dir + "sqlitekv_test")

		It("Create Read Update Delete", func() {
			k1 := "key1"
			v1 := "val1"
			k2 := "key2"
			v2 := "val2"
			v3 := "val3"

			// create, read
			skv.SetString(k1, v1)
			Expect(skv.HasKey(k1)).Should(BeTrue())
			Expect(skv.String(k1)).Should(Equal(v1))

			// create, read
			Expect(skv.HasKey(k2)).Should(BeFalse())
			skv.SetString(k2, v2)
			Expect(skv.HasKey(k2)).Should(BeTrue())
			Expect(skv.String(k2)).Should(Equal(v2))

			// update
			skv.SetString(k2, v3)
			Expect(skv.String(k2)).Should(Equal(v3))

			// delete
			Expect(skv.HasKey(k1)).Should(BeTrue())
			Expect(skv.String(k1)).Should(Equal(v1))
			skv.Del(k1)
			Expect(skv.HasKey(k1)).Should(BeFalse())
		})
	})
})
