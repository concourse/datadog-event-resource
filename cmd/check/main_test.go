package main_test

import (
	"os/exec"

	"github.com/concourse/datadog-resource/cmd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
	"github.com/zorkian/go-datadog-api"
)

var _ = Describe("Check", func() {
	var (
		session *gexec.Session
	)

	Context("when called with no arguments", func() {
		AfterEach(func() {
			Expect(fakeDataDogServer.ReceivedRequests()).To(HaveLen(0))
		})

		It("fails when called with no arguments", func() {
			session, err = gexec.Start(exec.Command(binPath), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			<-session.Exited

			Expect(session).To(gexec.Exit(1))
		})
	})

	Context("when called with source configuration", func() {
		AfterEach(func() {
			Expect(fakeDataDogServer.ReceivedRequests()).To(HaveLen(1))
		})

		Context("when there are no events", func() {
			It("outputs an empty JSON", func() {
				RespondWithEvents(nil)

				session = RunCheck(cmd.CheckPayload{
					Source: cmd.Source{
						ApplicationKey: "foobar",
						ApiKey:         "barbaz",
					},
				})

				Expect(session).To(gbytes.Say("\\[\\]"))
			})
		})

		Context("when there is one event", func() {
			It("outputs a single element array with that version as id", func() {
				RespondWithEvents([]datadog.Event{
					{Id: 100},
				})

				session = RunCheck(cmd.CheckPayload{
					Source: cmd.Source{
						ApplicationKey: "foobar",
						ApiKey:         "barbaz",
					},
				})

				Expect(session).To(gbytes.Say(`\[{"id":"100"}\]`))
			})
		})
	})
})

type Response struct {
	Events []datadog.Event `json:"events"`
}

func RespondWithEvents(e []datadog.Event) {
	fakeDataDogServer.AppendHandlers(
		ghttp.RespondWithJSONEncoded(200, Response{
			Events: e,
		}),
	)
}
