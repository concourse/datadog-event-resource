package main

import (
	"encoding/json"
	"os"

	"strconv"

	"github.com/concourse/datadog-resource/cmd"
	"github.com/zorkian/go-datadog-api"
	"strings"
	"time"
)

func main() {
	// Parse payload
	var (
		err     error
		payload cmd.OutPayload
	)
	err = json.NewDecoder(os.Stdin).Decode(&payload)
	if err != nil {
		panic(err)
	}

	// Use datadog client to get event
	c := datadog.NewClient(payload.Source.ApiKey, payload.Source.ApplicationKey)

	e := ConvertParamsToEvent(payload.Params)
	event, err := c.PostEvent(&e)
	if err != nil {
		panic(err)
	}

	// Emit metadata
	metadata, err := GenerateMetadata(event)
	if err != nil {
		panic(err)
	}

	response := cmd.OutResponse{
		Version:  cmd.Version{Id: strconv.Itoa(event.Id)},
		Metadata: metadata,
	}

	err = json.NewEncoder(os.Stdout).Encode(response)
	if err != nil {
		panic(err)
	}
}

func ConvertParamsToEvent(p cmd.OutParams) (e datadog.Event) {
	e.Title = p.Title
	e.Text = p.Text
	e.Priority = p.Priority
	e.AlertType = p.AlertType
	e.Host = p.Host
	e.Aggregation = p.Aggregation
	e.SourceType = p.SourceType
	e.Tags = p.Tags

	return
}

var layout = "2006-01-02 15:04:05 -0700"

func GenerateMetadata(e *datadog.Event) ([]cmd.Metadata, error) {
	metadata := make([]cmd.Metadata, 0)

	if e.Id != 0 {
		metadata = append(metadata, cmd.Metadata{Name: "id", Value: strconv.Itoa(e.Id)})
	}

	if e.Title != "" {
		metadata = append(metadata, cmd.Metadata{Name: "title", Value: e.Title})
	}

	if e.Text != "" {
		metadata = append(metadata, cmd.Metadata{Name: "text", Value: e.Text})
	}

	if e.Time != 0 {
		date := time.Unix(int64(e.Time), 0)
		metadata = append(metadata, cmd.Metadata{Name: "date_happened", Value: date.Format(layout)})
	}

	if e.Priority != "" {
		metadata = append(metadata, cmd.Metadata{Name: "priority", Value: e.Priority})
	}

	if e.AlertType != "" {
		metadata = append(metadata, cmd.Metadata{Name: "alert_type", Value: e.AlertType})
	}

	if e.Host != "" {
		metadata = append(metadata, cmd.Metadata{Name: "host", Value: e.Host})
	}

	if e.Aggregation != "" {
		metadata = append(metadata, cmd.Metadata{Name: "aggregation_key", Value: e.Aggregation})
	}

	if e.SourceType != "" {
		metadata = append(metadata, cmd.Metadata{Name: "source_type_name", Value: e.SourceType})
	}

	if len(e.Tags) > 0 {
		t := strings.Join(e.Tags, ", ")

		metadata = append(metadata, cmd.Metadata{Name: "tags", Value: t})
	}

	return metadata, nil
}
