package main_test

import (
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

	Context("when called with source configuration but no version", func() {
		AfterEach(func() {
			Expect(fakeDataDogServer.ReceivedRequests()).To(HaveLen(1))
		})

		Context("when there are no events", func() {
			It("outputs an empty JSON", func() {
				RespondWithEvents(nil)

				session = RunCheck(nil)

				Expect(session.Out).To(gbytes.Say("\\[\\]"))
			})
		})

		Context("when there is one event", func() {
			It("outputs a single element array with that version as id", func() {
				RespondWithEvents([]datadog.Event{
					{Id: 100},
				})

				session = RunCheck(nil)

				Expect(session.Out).To(gbytes.Say(`\[{"id":"100"}\]`))
			})
		})

		Context("when there are multiple events", func() {
			It("outputs a single element array with the first version (most recent) as id", func() {
				RespondWithEvents([]datadog.Event{
					{Id: 100},
					{Id: 99},
					{Id: 98},
					{Id: 97},
				})

				session = RunCheck(nil)

				Expect(session.Out).To(gbytes.Say(`\[{"id":"100"}\]`))
			})
		})
	})

	Context("when called with source configuration but no version", func() {
		AfterEach(func() {
			Expect(fakeDataDogServer.ReceivedRequests()).To(HaveLen(1))
		})

		Context("when there are no events", func() {
			It("outputs an empty JSON", func() {
				RespondWithEvents(nil)

				session = RunCheck(nil)

				Expect(session.Out).To(gbytes.Say("\\[\\]"))
			})
		})

		Context("when there is one event", func() {
			It("outputs a single element array with that version as id", func() {
				RespondWithEvents([]datadog.Event{
					{Id: 100},
				})

				session = RunCheck(nil)

				Expect(session.Out).To(gbytes.Say(`\[{"id":"100"}\]`))
			})
		})

		Context("when there are multiple events", func() {
			It("outputs a single element array with the first version (most recent) as id", func() {
				RespondWithEvents([]datadog.Event{
					{Id: 100},
					{Id: 99},
					{Id: 98},
					{Id: 97},
				})

				session = RunCheck(nil)

				Expect(session.Out).To(gbytes.Say(`\[{"id":"100"}\]`))
			})
		})
	})

	Context("when called with source configuration and version", func() {
		var (
			id string
		)

		BeforeEach(func() {
			id = "90"
		})

		AfterEach(func() {
			Expect(fakeDataDogServer.ReceivedRequests()).To(HaveLen(1))
		})

		Context("when there are no events", func() {
			It("outputs an empty JSON", func() {
				RespondWithEvents(nil)

				session = RunCheck(&id)

				Expect(session.Out).To(gbytes.Say("\\[\\]"))
			})
		})

		Context("when there is one event", func() {
			It("outputs a single element array with that version as id", func() {
				RespondWithEvents([]datadog.Event{
					{Id: 100},
				})

				session = RunCheck(&id)

				Expect(session.Out).To(gbytes.Say(`\[{"id":"100"}\]`))
			})
		})

		Context("when there are multiple events", func() {
			It("outputs a array with the all events more recent (self-inclusive), reversed", func() {
				RespondWithEvents([]datadog.Event{
					{Id: 110},
					{Id: 100},
					{Id: 90},
					{Id: 80},
				})

				session = RunCheck(&id)

				Expect(session.Out).To(gbytes.Say(`\[{"id":"90"},{"id":"100"},{"id":"110"}\]`))
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
