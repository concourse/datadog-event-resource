package cmd

import "errors"

type Source struct {
	ApplicationKey string `json:"application_key"`
	ApiKey         string `json:"api_key"`
	Filter         string `json:"filter,omitempty"`
	DatadogHost    string `json:"datadog_host,omitempty"`
}

type Version struct {
	Id string `json:"id"`
}

type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CheckPayload struct {
	Source  Source   `json:"source"`
	Version *Version `json:"version,omitempty"`
}

type CheckResponse []Version

type InPayload struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type InResponse struct {
	Version  Version    `json:"version"`
	Metadata []Metadata `json:"metadata"`
}

type OutPayload struct {
	Source Source    `json:"source"`
	Params OutParams `json:"params"`
}

type OutParams struct {
	Title       string   `json:"title,omitempty"`
	Text        string   `json:"text,omitempty"`
	TextFile    string   `json:"text_file,omitempty"`
	Priority    string   `json:"priority,omitempty"`
	AlertType   string   `json:"alert_type,omitempty"`
	Host        string   `json:"host,omitempty"`
	Aggregation string   `json:"aggregation_key,omitempty"`
	SourceType  string   `json:"source_type_name,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type OutResponse struct {
	Version  Version    `json:"version"`
	Metadata []Metadata `json:"metadata"`
}

var NoTextOrTextFileInOutParamsErr = errors.New("No Text or TextFile found in params")
