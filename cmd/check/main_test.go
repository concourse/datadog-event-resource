package main_test

import (
	"os/exec"

	"github.com/concourse/datadog-resource/cmd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
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
				fakeDataDogServer.AppendHandlers(
					ghttp.RespondWith(200, nil, nil),
				)

				session = RunCheck(cmd.CheckPayload{
					Source: cmd.Source{
						ApplicationKey: "foobar",
						ApiKey:         "barbaz",
					},
				})

				Expect(session).To(gbytes.Say("\\[\\]"))
			})
		})
	})
})
