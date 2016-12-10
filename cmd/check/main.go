package main

import (
	"bytes"
	"encoding/json"
	"os"

	"math"

	"strconv"

	"github.com/concourse/datadog-resource/cmd"
	"github.com/zorkian/go-datadog-api"
)

func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}

	var payload cmd.CheckPayload
	json.NewDecoder(bytes.NewBufferString(os.Args[1])).Decode(&payload)

	c := datadog.NewClient(payload.Source.ApiKey, payload.Source.ApplicationKey)

	events, err := c.GetEvents(0, math.MaxInt8, "", "", "")
	if err != nil {
		panic(err)
	}

	output := make(Output, 0)

	if payload.Version == nil {
		if len(events) > 0 {
			e := events[0]

			output = append(
				output,
				cmd.Version{
					Id: strconv.Itoa(e.Id),
				})
		}
	}

	json.NewEncoder(os.Stdout).Encode(&output)
}

type Output []cmd.Version
