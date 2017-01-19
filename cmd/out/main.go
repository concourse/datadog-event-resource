package main

import (
	"encoding/json"
	"os"

	"strconv"

	"strings"
	"time"

	"errors"

	"io/ioutil"
	"path/filepath"

	"github.com/concourse/datadog-event-resource/cmd"
	"github.com/zorkian/go-datadog-api"
)

func main() {
	artifactDirectory := os.Args[1]
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

	e, err := ConvertParamsToEvent(payload.Params, artifactDirectory)
	if err != nil {
		panic(err)
	}

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

func ConvertParamsToEvent(p cmd.OutParams, artifactDirectory string) (datadog.Event, error) {
	var (
		e datadog.Event
	)
	e.Title = p.Title
	if p.Text != "" {
		e.Text = p.Text
	} else if p.TextFile != "" {
		b, err := ioutil.ReadFile(filepath.Join(artifactDirectory, p.TextFile))
		if err != nil {
			return e, err
		}
		e.Text = string(b)
	} else {
		return e, errors.New("No Text or TextFile found in params")
	}
	e.Priority = p.Priority
	e.AlertType = p.AlertType
	e.Host = p.Host
	e.Aggregation = p.Aggregation
	e.SourceType = p.SourceType
	e.Tags = p.Tags

	return e, nil
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
		metadata = append(metadata, cmd.Metadata{Name: "date_happened", Value: date.Local().Format(layout)})
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
