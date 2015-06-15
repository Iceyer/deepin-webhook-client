package client

type Event struct {
	ID        string                 `json:"id"`
	Publisher string                 `json:"publisher"`
	Name      string                 `json:"name"`
	Secret    string                 `json:"secret"`
	Schema    map[string]interface{} `json:"schema"`
}
