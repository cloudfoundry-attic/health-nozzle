package server_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/health-nozzle/counter"
	"github.com/cloudfoundry-incubator/health-nozzle/server"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	Context("Given a set of counters", func() {
		var request *http.Request
		var response *httptest.ResponseRecorder

		BeforeEach(func() {
			var err error
			counters := []*counter.HealthCounter{
				counter.NewMetronHealthCounter(),
				counter.NewDopplerHealthCounter(),
			}
			server := server.NewServer(counters)

			request, err = http.NewRequest("GET", "/", nil)
			Expect(err).ToNot(HaveOccurred())

			response = httptest.NewRecorder()
			server.ServeHTTP(response, request)
		})

		It("returns 200 OK", func() {
			Expect(response.Code).To(Equal(http.StatusOK))
		})

		It("sets the JSON content type", func() {
			Expect(response.HeaderMap["Content-Type"]).To(ContainElement("application/json"))
		})

		It("returns a nice JSON response", func() {
			Expect(response.Body).To(MatchJSON(`
				{
					"metron": {
						"ingress": 0,
						"egress": 0,
						"dropped": 0
					},
					"doppler": {
						"ingress": 0,
						"egress": 0,
						"dropped": 0
					}
				}
			`))
		})
	})
})
