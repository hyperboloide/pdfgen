package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

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

})
