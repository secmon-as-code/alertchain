package types

type Parameter struct {
	Key   string        `json:"key"`
	Value string        `json:"value"`
	Type  ParameterType `json:"type"`
}

type ParameterType string

const (
	IPAddr     ParameterType = "ipaddr"
	DomainName ParameterType = "domain"
	FileSha256 ParameterType = "file.sha256"
	FileSha512 ParameterType = "file.sha512"
)
