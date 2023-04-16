package types

type (
	ParamKey      string
	ParamValue    any
	ParameterType string
)

const (
	IPAddr     ParameterType = "ipaddr"
	DomainName ParameterType = "domain"
	FileSha256 ParameterType = "file.sha256"
	FileSha512 ParameterType = "file.sha512"
)
