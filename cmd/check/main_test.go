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

	BeforeEach(func() {
		session = nil
	})

	Context("when called with source configuration but no version", func() {
		AfterEach(func() {
			Expect(fakeDataDogServer.ReceivedRequests()).To(HaveLen(1))
		})

		Context("when there are no events", func() {
			It("outputs an empty JSON", func() {
				RespondWithEvents(nil)

				session = RunCheckSuccessfully(nil, "")

				Expect(session.Out).To(gbytes.Say("\\[\\]"))
			})
		})

		Context("when there is one event", func() {
			It("outputs a single element array with that version as id", func() {
				RespondWithEvents([]datadog.Event{
					{Id: 100, Time: 10},
				})

				session = RunCheckSuccessfully(nil, "")

				Expect(session.Out).To(gbytes.Say(`\[{"id":"100"}\]`))
			})
		})

		Context("when there are multiple events", func() {
			It("outputs a single element array with the most recent time version as id", func() {
				RespondWithEvents([]datadog.Event{
					{Id: 100, Time: 80},
					{Id: 99, Time: 90},
					{Id: 98, Time: 100},
					{Id: 97, Time: 70},
				})

				session = RunCheckSuccessfully(nil, "")

				Expect(session.Out).To(gbytes.Say(`\[{"id":"98"}\]`))
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

				session = RunCheckSuccessfully(nil, "")

				Expect(session.Out).To(gbytes.Say("\\[\\]"))
			})
		})

		Context("when there is one event", func() {
			It("outputs a single element array with that version as id", func() {
				RespondWithEvents([]datadog.Event{
					{Id: 100, Time: 100},
				})

				session = RunCheckSuccessfully(nil, "")

				Expect(session.Out).To(gbytes.Say(`\[{"id":"100"}\]`))
			})
		})

		Context("when there are multiple events", func() {
			It("outputs a single element array with the most recent time version as id", func() {
				RespondWithEvents([]datadog.Event{
					{Id: 100, Time: 80},
					{Id: 99, Time: 90},
					{Id: 98, Time: 100},
					{Id: 97, Time: 70},
				})

				session = RunCheckSuccessfully(nil, "")

				Expect(session.Out).To(gbytes.Say(`\[{"id":"98"}\]`))
			})
		})
	})

	Context("when called with source configuration and version", func() {
		var (
			id string
		)

		BeforeEach(func() {
			id = "100"
		})

		AfterEach(func() {
			Expect(fakeDataDogServer.ReceivedRequests()).To(HaveLen(1))
		})

		Context("when there are no events", func() {
			It("outputs an empty JSON", func() {
				RespondWithEvents(nil)

				session = RunCheckSuccessfully(&id, "")

				Expect(session.Out).To(gbytes.Say("\\[\\]"))
			})
		})

		Context("when there is one event", func() {
			It("outputs a single element array with that version as id", func() {
				RespondWithEvents([]datadog.Event{
					{Id: 100},
				})

				session = RunCheckSuccessfully(&id, "")

				Expect(session.Out).To(gbytes.Say(`\[{"id":"100"}\]`))
			})
		})

		Context("when there are multiple events", func() {
			It("outputs a array with the all events more recent (self-inclusive), reversed, with the needle at the beginning", func() {
				RespondWithEvents([]datadog.Event{
					{Id: 100, Time: 80},
					{Id: 99, Time: 90},
					{Id: 98, Time: 100},
					{Id: 97, Time: 70},
					{Id: 96, Time: 90},
					{Id: 95, Time: 80},
				})

				session = RunCheckSuccessfully(&id, "")

				Expect(session.Out).To(gbytes.Say(`\[{"id":"100"},{"id":"95"},{"id":"96"},{"id":"99"},{"id":"98"}\]`))
			})

			Context("when the list of events does not include the one requested", func() {
				It("outputs the most recent version", func() {
					RespondWithEvents([]datadog.Event{
						{Id: 101, Time: 80},
						{Id: 99, Time: 90},
						{Id: 98, Time: 100},
						{Id: 97, Time: 70},
						{Id: 96, Time: 90},
						{Id: 95, Time: 80},
					})

					session = RunCheckSuccessfully(&id, "")

					Expect(session.Out).To(gbytes.Say(`\[{"id":"98"}\]`))
				})
			})
		})
	})

	Context("when called with source configuration with filter", func() {
		Context("when there are many events, but none match the given filter", func() {
			BeforeEach(func() {
				RespondWithEvents([]datadog.Event{
					{Title: "Hello-world", Id: 100, Time: 80},
					{Title: "Hello-world", Id: 99, Time: 90},
					{Title: "Hello-world", Id: 98, Time: 100},
					{Title: "Hello-world", Id: 97, Time: 70},
					{Title: "Hello-world", Id: 96, Time: 90},
					{Title: "Hello-world", Id: 95, Time: 80},
				})
			})

			It("doesn't print any event", func() {
				session = RunCheckSuccessfully(nil, "good-bye")
				Expect(session.Out).To(gbytes.Say(`\[\]`))
			})
		})

		Context("when there are many events, with some matches", func() {
			BeforeEach(func() {
				RespondWithEvents([]datadog.Event{
					{Title: "Hello-world", Id: 100, Time: 80},
					{Title: "good-bye", Id: 99, Time: 90},
					{Title: "g00d-bye", Id: 101, Time: 100},
					{Title: "good-bye", Id: 97, Time: 70},
					{Title: "Hello-world", Id: 96, Time: 90},
					{Title: "Hello-world", Id: 95, Time: 80},
				})
			})

			It("prints the latest match", func() {
				session = RunCheckSuccessfully(nil, "good-bye")
				Expect(session.Out).To(gbytes.Say(`\[{"id":"99"}\]`))
			})

			Context("when the filter is a regexp", func() {
				It("prints the latest match", func() {
					session = RunCheckSuccessfully(nil, "g[a-z]*-bye")
					Expect(session.Out).To(gbytes.Say(`\[{"id":"99"}\]`))
				})
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
