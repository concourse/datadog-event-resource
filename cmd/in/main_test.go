package main_test

import (
	"encoding/json"
	"time"

	"strconv"

	"io/ioutil"
	"path/filepath"

	"bytes"

	"github.com/concourse/datadog-event-resource/cmd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
	"github.com/zorkian/go-datadog-api"
)

var _ = Describe("In", func() {
	var (
		session *gexec.Session

		inResponse cmd.InResponse

		event     datadog.Event
		t         time.Time
		timestamp string

		layout = "2006-01-02 15:04:05 -0700"
	)

	BeforeEach(func() {
		timestamp = "2016-12-12 14:33:04 -0800"
		t, err = time.Parse(layout, timestamp)
		Expect(err).NotTo(HaveOccurred())

		event = datadog.Event{
			Id:        id,
			Title:     "some-datadog-event",
			Text:      "some-datadog-event-text",
			Time:      int(t.Unix()),
			Priority:  "normal",
			AlertType: "info",
		}

		RespondWithEvent(event)
	})

	It("grabs the event corresponding to the id, and emits metadata", func() {
		session = RunIn()
		Expect(fakeDataDogServer.ReceivedRequests()).To(HaveLen(1))

		err = json.NewDecoder(bytes.NewBuffer(session.Out.Contents())).Decode(&inResponse)
		Expect(err).NotTo(HaveOccurred())

		Expect(inResponse.Version).To(Equal(cmd.Version{
			Id: strconv.Itoa(id),
		}))

		Expect(inResponse.Metadata).To(ConsistOf(
			cmd.Metadata{Name: "id", Value: strconv.Itoa(id)},
			cmd.Metadata{Name: "title", Value: "some-datadog-event"},
			cmd.Metadata{Name: "text", Value: "some-datadog-event-text"},
			cmd.Metadata{Name: "date_happened", Value: t.Local().Format(layout)},
			cmd.Metadata{Name: "priority", Value: "normal"},
			cmd.Metadata{Name: "alert_type", Value: "info"},
		))
	})

	It("writes an `event.json` file", func() {
		session = RunIn()
		Expect(fakeDataDogServer.ReceivedRequests()).To(HaveLen(1))

		b, err := ioutil.ReadFile(filepath.Join(tmpDir, "event.json"))
		Expect(err).NotTo(HaveOccurred())

		var (
			persistedEvent datadog.Event
		)

		err = json.Unmarshal(b, &persistedEvent)
		Expect(err).NotTo(HaveOccurred())

		Expect(persistedEvent).To(Equal(event))
	})

	It("writes a `version` file", func() {
		session = RunIn()
		Expect(fakeDataDogServer.ReceivedRequests()).To(HaveLen(1))

		b, err := ioutil.ReadFile(filepath.Join(tmpDir, "version"))
		Expect(err).NotTo(HaveOccurred())

		Expect(string(b)).To(Equal(strconv.Itoa(id)))
	})
})

type Response struct {
	Event datadog.Event `json:"event"`
}

func RespondWithEvent(e datadog.Event) {
	fakeDataDogServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/api/v1/events/"+strconv.Itoa(id)),
			ghttp.RespondWithJSONEncoded(200, Response{
				Event: e,
			}),
		),
	)
}
