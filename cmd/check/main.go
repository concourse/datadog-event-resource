package main

import (
	"encoding/json"
	"os"

	"strconv"

	"time"

	"github.com/concourse/datadog-event-resource/cmd"
	"github.com/zorkian/go-datadog-api"
)

func main() {
	var payload cmd.CheckPayload
	err := json.NewDecoder(os.Stdin).Decode(&payload)
	if err != nil {
		panic(err)
	}

	c := datadog.NewClient(payload.Source.ApiKey, payload.Source.ApplicationKey)

	end := int(time.Now().Add(24 * time.Hour).Unix())
	beginning := end - 2764800
	events, err := c.GetEvents(beginning, end, "", "", "")
	if err != nil {
		panic(err)
	}

	output := make(cmd.CheckResponse, 0)

	var e datadog.Event
	if payload.Version == nil {
		if len(events) > 0 {
			e = events[0]

			output = AddEvent(output, e)
		}
	} else {
		switch len(events) {
		case 0:
			break
		case 1:
			e = events[0]
			output = AddEvent(output, e)
			break
		default:
			needle, err := strconv.Atoi(payload.Version.Id)
			if err != nil {
				panic(err)
			}

			for _, e = range events {
				if e.Id >= needle {
					output = AddEvent(output, e)
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

func AddEvent(o cmd.CheckResponse, e datadog.Event) cmd.CheckResponse {
	return append(
		o,
		cmd.Version{
			Id: strconv.Itoa(e.Id),
		})
}
