package http_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/storyscript/scheduler/http"
)

var _ = Describe("The healthcheck handler", func() {

	It("returns 200 and a body saying OK", func() {
		server := Server{}
		handler := http.HandlerFunc(server.HandleHealthcheck)

		request, err := http.NewRequest("GET", "/healthcheck", nil)
		Expect(err).NotTo(HaveOccurred())

		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)

		Expect(recorder.Code).To(Equal(http.StatusOK))
		Expect(recorder.Body.String()).To(Equal("OK"))
	})
})
