package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/hyperboloide/pdfgen"
)

var _ = Describe("Template", func() {

	It("should not create a NewTemplate if there is no index.html in the directory", func() {
		_, err := NewTemplate("./", "templates")
		Expect(err).To(HaveOccurred())
	})

	It("should not create a NewTemplate if the directory doesn't exists", func() {
		_, err := NewTemplate("./", "not_found")
		Expect(err).To(HaveOccurred())
	})

})
