package model

type Config struct {
	Policy   PolicyConfig     `json:"policy"`
	Database DBConfig         `json:"database"`
	Actions  ActionDefinition `json:"actions"`
	Jobs     JobDefinition    `json:"jobs"`
}

type PolicyConfig struct {
	Type string `json:"type"`
	Path string `json:"path"`
	URL  string `json:"url"`
}

type DBConfig struct {
	Type   string `json:"type"`
	Config string `json:"config"`
}
