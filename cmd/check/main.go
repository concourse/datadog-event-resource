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

	var e datadog.Event
	if payload.Version == nil {
		if len(events) > 0 {
			e = events[0]

			output = output.AddEvent(e)
		}
	} else {
		switch len(events) {
		case 0:
			break
		case 1:
			e = events[0]
			output = output.AddEvent(e)
			break
		default:
			needle, err := strconv.Atoi(payload.Version.Id)
			if err != nil {
				panic(err)
			}

			for _, e = range events {
				if e.Id >= needle {
					output = output.AddEvent(e)
				}
			}

			// Reverse
			for i := len(output)/2 - 1; i >= 0; i-- {
				opp := len(output) - 1 - i
				output[i], output[opp] = output[opp], output[i]
			}

			break
		}
	}

	json.NewEncoder(os.Stdout).Encode(&output)
}

type Output []cmd.Version

func (o Output) AddEvent(e datadog.Event) Output {
	return append(
		o,
		cmd.Version{
			Id: strconv.Itoa(e.Id),
		})
}
