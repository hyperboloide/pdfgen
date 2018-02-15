package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPdfgen(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pdfgen Suite")
}
