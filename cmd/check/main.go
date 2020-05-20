package main

import (
	"encoding/json"
	"os"
	"regexp"

	"strconv"

	"time"

	"sort"

	"github.com/concourse/datadog-event-resource/cmd"
	"github.com/zorkian/go-datadog-api"
)

func main() {
	var payload cmd.CheckPayload
	err := json.NewDecoder(os.Stdin).Decode(&payload)
	if err != nil {
		panic(err)
	}

	if payload.Source.DatadogHost != "" {
		os.Setenv("DATADOG_HOST", payload.Source.DatadogHost)
	}

	c := datadog.NewClient(payload.Source.ApiKey, payload.Source.ApplicationKey)

	end := int(time.Now().Add(24 * time.Hour).Unix())
	beginning := end - 2764800
	events, err := c.GetEvents(beginning, end, "", "", "")
	if err != nil {
		panic(err)
	}

	if payload.Source.Filter != "" {
		events, err = FilterEventsByTitle(events, payload.Source.Filter)
		if err != nil {
			panic(err)
		}
	}
	sort.Sort(ByLaterDateAndEarlierId(events))

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
			i, err := strconv.Atoi(payload.Version.Id)
			if err != nil {
				panic(err)
			}

			var (
				needle datadog.Event
				found  bool
			)
			for _, e = range events {
				if e.Id == i {
					needle = e
					found = true
					break
				}
			}

			if !found {
				e = events[0]
				output = AddEvent(output, e)
			} else {
				for _, e = range events {
					if e.Time >= needle.Time && e.Id != needle.Id {
						output = AddEvent(output, e)
					}
				}

				if found {
					output = AddEvent(output, needle)
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

type ByLaterDateAndEarlierId []datadog.Event

func (a ByLaterDateAndEarlierId) Len() int      { return len(a) }
func (a ByLaterDateAndEarlierId) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByLaterDateAndEarlierId) Less(i, j int) bool {
	return a[i].Time > a[j].Time && a[i].Id < a[j].Id
}

func FilterEventsByTitle(events []datadog.Event, filter string) ([]datadog.Event, error) {
	filterRegexp, err := regexp.Compile(filter)
	if err != nil {
		return nil, err
	}

	filteredEvents := []datadog.Event{}
	for _, event := range events {
		if filterRegexp.MatchString(event.Title) {
			filteredEvents = append(filteredEvents, event)
		}
	}

	return filteredEvents, nil
}
