package director_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	. "github.com/cloudfoundry/bosh-cli/director"
)

var _ = Describe("Director", func() {
	var (
		director   Director
		deployment Deployment
		server     *ghttp.Server
	)

	BeforeEach(func() {
		director, server = BuildServer()

		var err error

		deployment, err = director.FindDeployment("dep1")
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("Errands", func() {
		It("returns errands", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/deployments/dep1/errands"),
					ghttp.VerifyBasicAuth("username", "password"),
					ghttp.RespondWith(http.StatusOK, `[
	{"name": "errand1"},
	{"name": "errand2"}
]`),
				),
			)

			errands, err := deployment.Errands()
			Expect(err).ToNot(HaveOccurred())
			Expect(errands).To(Equal([]Errand{
				{Name: "errand1"},
				{Name: "errand2"},
			}))
		})

		It("returns error if response is non-200", func() {
			AppendBadRequest(ghttp.VerifyRequest("GET", "/deployments/dep1/errands"), server)

			_, err := deployment.Errands()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(
				"Finding errands: Director responded with non-successful status code"))
		})

		It("returns error if response cannot be unmarshalled", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/deployments/dep1/errands"),
					ghttp.RespondWith(http.StatusOK, ``),
				),
			)

			_, err := deployment.Errands()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(
				"Finding errands: Unmarshaling Director response"))
		})
	})

	Describe("RunErrand", func() {
		It("runs errand and returns result", func() {
			respBody := `{
	"exit_code":1,
	"stdout":"stdout",
	"stderr":"stderr",
	"logs": {
		"blobstore_id": "logs-blob-id",
		"sha1": "logs-sha1"
	}
}`
			ConfigureTaskResult(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/deployments/dep1/errands/errand1/runs"),
					ghttp.VerifyBasicAuth("username", "password"),
					ghttp.VerifyHeader(http.Header{
						"Content-Type": []string{"application/json"},
					}),
					ghttp.VerifyBody([]byte(`{"keep-alive":false}`)),
				),
				respBody,
				server,
			)

			result, err := deployment.RunErrand("errand1", false)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(ErrandResult{
				ExitCode: 1,

				Stdout: "stdout",
				Stderr: "stderr",

				LogsBlobstoreID: "logs-blob-id",
				LogsSHA1:        "logs-sha1",
			}))
		})

		It("runs errand, keeping it alive and returns result", func() {
			ConfigureTaskResult(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/deployments/dep1/errands/errand1/runs"),
					ghttp.VerifyBasicAuth("username", "password"),
					ghttp.VerifyHeader(http.Header{
						"Content-Type": []string{"application/json"},
					}),
					ghttp.VerifyBody([]byte(`{"keep-alive":true}`)),
				),
				`{"exit_code":1}`,
				server,
			)

			result, err := deployment.RunErrand("errand1", true)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(ErrandResult{ExitCode: 1}))
		})

		It("returns error if response is non-200", func() {
			AppendBadRequest(ghttp.VerifyRequest("POST", "/deployments/dep1/errands/errand1/runs"), server)

			_, err := deployment.RunErrand("errand1", false)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(
				"Running errand 'errand1': Director responded with non-successful status code"))
		})

		It("returns error if task result cannot be unmarshalled", func() {
			ConfigureTaskResult(ghttp.VerifyRequest("POST", "/deployments/dep1/errands/errand1/runs"), "", server)

			_, err := deployment.RunErrand("errand1", false)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Unmarshaling errand result"))
		})
	})
})
