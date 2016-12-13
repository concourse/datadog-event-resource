package main_test

import (
	"time"

	"bytes"
	"encoding/json"

	"github.com/concourse/datadog-resource/cmd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
	"github.com/zorkian/go-datadog-api"
)

var _ = Describe("Out", func() {
	var (
		session *gexec.Session

		outResponse cmd.OutResponse

		t         time.Time
		timestamp string

		layout = "2006-01-02 15:04:05 -0700"
	)

	BeforeEach(func() {
		timestamp = "2016-12-12 14:33:04 -0800"
		t, err = time.Parse(layout, timestamp)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("when given params containing a static event to emit", func() {
		var (
			event  datadog.Event
			params cmd.OutParams
		)

		BeforeEach(func() {
			params = cmd.OutParams{
				Title:       "some-datadog-event",
				Text:        "some-datadog-event-text",
				Priority:    "normal",
				AlertType:   "info",
				Host:        "localhost",
				Aggregation: "some-aggregation-key",
				SourceType:  "some-source-type",
				Tags: []string{
					"some-tag",
					"some-other-tag",
				},
			}

			event = datadog.Event{
				Id:          1234,
				Title:       "some-datadog-event",
				Text:        "some-datadog-event-text",
				Time:        int(t.Unix()),
				Priority:    "normal",
				AlertType:   "info",
				Host:        "localhost",
				Aggregation: "some-aggregation-key",
				SourceType:  "some-source-type",
				Tags: []string{
					"some-tag",
					"some-other-tag",
				},
			}

			fakeDataDogServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/api/v1/events"),
					ghttp.RespondWithJSONEncoded(200, Response{
						Event: event,
					}),
				),
			)
		})

		It("creates the event via the API", func() {
			session = RunOut(params)
			Expect(fakeDataDogServer.ReceivedRequests()).To(HaveLen(1))
		})

		It("emits metadata about the event", func() {
			session = RunOut(params)

			err = json.NewDecoder(bytes.NewBuffer(session.Out.Contents())).Decode(&outResponse)
			Expect(err).NotTo(HaveOccurred())

			Expect(outResponse.Version).To(Equal(cmd.Version{
				Id: "1234",
			}))

			Expect(outResponse.Metadata).To(ConsistOf(
				cmd.Metadata{Name: "id", Value: "1234"},
				cmd.Metadata{Name: "title", Value: "some-datadog-event"},
				cmd.Metadata{Name: "text", Value: "some-datadog-event-text"},
				cmd.Metadata{Name: "date_happened", Value: t.Local().Format(layout)},
				cmd.Metadata{Name: "priority", Value: "normal"},
				cmd.Metadata{Name: "alert_type", Value: "info"},
				cmd.Metadata{Name: "host", Value: "localhost"},
				cmd.Metadata{Name: "aggregation_key", Value: "some-aggregation-key"},
				cmd.Metadata{Name: "source_type_name", Value: "some-source-type"},
				cmd.Metadata{Name: "tags", Value: "some-tag, some-other-tag"},
			))
		})
	})
})

type Response struct {
	Event datadog.Event `json:"event"`
}
