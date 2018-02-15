package main_test

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"

	. "github.com/hyperboloide/pdfgen"
)

var _ = Describe("Config", func() {

	It("IsValidTemplateDir should work in the 'templates' dir", func() {
		Expect(IsValidTemplateDir("./templates")).To(BeTrue())
	})

	It("IsValidTemplateDir should return false on a path that don't exists", func() {
		Expect(IsValidTemplateDir("./do_not_exists")).To(BeFalse())
	})

	It("IsValidTemplateDir should return false on a path that is not a directory", func() {
		Expect(IsValidTemplateDir("./main.go")).To(BeFalse())
	})

	It("SelectDir should return a valid path", func() {
		dirs := []string{
			"./do_not_exists",
			"./templates",
			"./main.go",
		}
		Expect(*SelectDir(dirs)).To(Equal("./templates"))
	})

	It("SelectDir should return nil if no valid path is found", func() {
		dirs := []string{
			"./do_not_exists",
			"./main.go",
		}
		Expect(SelectDir(dirs)).To(BeNil())
	})

	Context("bad path", func() {
		path := os.Getenv("PATH")

		_ = BeforeEach(func() {
			os.Setenv("PATH", "")
		})

		_ = AfterEach(func() {
			os.Setenv("PATH", path)
		})

		It("read the config and set the error", func() {
			err := ConfigRead()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("executable wkhtmltopdf could not be found in PATH"))
		})

	})

	It("set invalid template dir", func() {
		viper.Set("templates", "invalid")
		err := ConfigRead()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("invalid template directory"))
	})

	It("set invalid template dir", func() {
		viper.Set("templates", "")
		err := ConfigRead()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("template directory not found"))
	})

	It("no template found", func() {
		dir, err := ioutil.TempDir("", "")
		Expect(err).ToNot(HaveOccurred())
		defer os.RemoveAll(dir)
		viper.Set("templates", dir)
		err = ConfigRead()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("No template found"))
	})

})
