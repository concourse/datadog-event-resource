package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"strconv"

	"time"

	"github.com/concourse/datadog-resource/cmd"
	"github.com/zorkian/go-datadog-api"
)

func main() {
	if len(os.Args) < 2 {
		panic("must be called with filepath as $1")
	}

	dir := os.Args[1]

	// Parse payload
	var (
		err     error
		payload cmd.InPayload
	)
	err = json.NewDecoder(os.Stdin).Decode(&payload)
	if err != nil {
		panic(err)
	}

	// Use datadog client to get event
	c := datadog.NewClient(payload.Source.ApiKey, payload.Source.ApplicationKey)

	i, err := strconv.Atoi(payload.Version.Id)
	if err != nil {
		panic(err)
	}

	event, err := c.GetEvent(i)
	if err != nil {
		panic(err)
	}

	// Write `version` and `event.json` files
	err = ioutil.WriteFile(filepath.Join(dir, "version"), []byte(strconv.Itoa(event.Id)), 0644)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(filepath.Join(dir, "event.json"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(event)
	if err != nil {
		panic(err)
	}

	// Emit event metadata
	metadata, err := GenerateMetadata(event)
	if err != nil {
		panic(err)
	}

	response := cmd.InResponse{
		Version:  payload.Version,
		Metadata: metadata,
	}

	err = json.NewEncoder(os.Stdout).Encode(&response)
	if err != nil {
		panic(err)
	}
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

	return metadata, nil
}
