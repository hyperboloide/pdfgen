package main_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"

	. "github.com/hyperboloide/pdfgen"
)

var _ = Describe("Main", func() {

	var srv *httptest.Server
	var docSize int
	var js []byte

	_ = BeforeEach(func() {
		viper.Set("templates", "./templates")
		err := ConfigRead()
		Expect(err).ToNot(HaveOccurred())

		srv = httptest.NewServer(Router())
		tmp, err := ioutil.ReadFile("./templates/demo.json")
		Expect(err).To(BeNil())
		js = tmp
	})

	_ = AfterEach(func() {
		srv.Close()
	})

	It("should respond StatusMethodNotAllowed to GET requests", func() {
		resp, err := http.Get(srv.URL + "/invoice")
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusMethodNotAllowed))
	})

	It("should respond StatusMethodNotAllowed is content type is not application/json", func() {
		resp, err := http.Post(srv.URL+"/invoice", "text/plain", bytes.NewReader(js))
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
	})

	It("should respond StatusNotFound if the template do not exists", func() {
		resp, err := http.Post(srv.URL+"/not_found", "application/json", bytes.NewReader(js))
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
	})

	It("should respond StatusBadRequest json is invalid", func() {
		resp, err := http.Post(srv.URL+"/invoice", "application/json", strings.NewReader("invalid!!!"))
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
	})

	It("should make a valid post", func() {
		resp, err := http.Post(srv.URL+"/invoice", "application/json", bytes.NewReader(js))
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		b, err := ioutil.ReadAll(resp.Body)
		Expect(err).ToNot(HaveOccurred())
		docSize = len(b)
	})

	It("should work with appended slash", func() {
		resp, err := http.Post(srv.URL+"/invoice/", "application/json", bytes.NewReader(js))
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		b, err := ioutil.ReadAll(resp.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(len(b)).To(Equal(docSize))
	})

	// It("should be able to run in parallel", func() {
	//
	// 	var test = func() pool.WorkFunc {
	// 		return func(wu pool.WorkUnit) (interface{}, error) {
	// 			resp, err := http.Post(srv.URL+"/invoice", "application/json", bytes.NewReader(js))
	// 			if err != nil || resp.StatusCode != http.StatusOK {
	// 				return false, err
	// 			} else if b, err := ioutil.ReadAll(resp.Body); err != nil {
	// 				return false, err
	// 			} else if len(b) != docSize {
	// 				return false, errors.New("invalid size")
	// 			}
	// 			return true, nil
	// 		}
	// 	}
	//
	// 	p := pool.NewLimited(20)
	// 	defer p.Close()
	// 	batch := p.Batch()
	//
	// 	go func() {
	// 		for i := 0; i < 100; i++ {
	// 			batch.Queue(test())
	// 		}
	// 		batch.QueueComplete()
	// 	}()
	//
	// 	for res := range batch.Results() {
	// 		Expect(res.Error()).ToNot(HaveOccurred())
	// 	}
	//
	// })

})
