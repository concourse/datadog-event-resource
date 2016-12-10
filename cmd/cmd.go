package cmd

type CheckPayload struct {
	Source Source `json:"source"`
}

type Source struct {
	ApplicationKey string `json:"application_key"`
	ApiKey         string `json:"api_key"`
}
