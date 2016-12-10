package cmd

type CheckPayload struct {
	Source Source `json:"source"`
	Version *Version `json:"version,omitempty"`
}

type Source struct {
	ApplicationKey string `json:"application_key"`
	ApiKey         string `json:"api_key"`
}

type Version struct {
	Id string `json:"id"`
}
