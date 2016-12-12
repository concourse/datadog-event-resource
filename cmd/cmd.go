package cmd

type CheckPayload struct {
	Source  Source   `json:"source"`
	Version *Version `json:"version,omitempty"`
}

type InPayload struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type Source struct {
	ApplicationKey string `json:"application_key"`
	ApiKey         string `json:"api_key"`
}

type Version struct {
	Id string `json:"id"`
}

type InResponse struct {
	Version  Version    `json:"version"`
	Metadata []Metadata `json:"metadata"`
}

type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
