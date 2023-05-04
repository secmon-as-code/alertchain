package types

import "github.com/google/uuid"

type (
	ParamID       string
	ParamName     string
	ParamValue    any
	ParameterType string
)

func NewParamID() ParamID {
	return ParamID(uuid.NewString())
}

const (
	IPAddr     ParameterType = "ipaddr"
	DomainName ParameterType = "domain"
	FileSha256 ParameterType = "file.sha256"
	FileSha512 ParameterType = "file.sha512"
	MarkDown   ParameterType = "markdown"
)
