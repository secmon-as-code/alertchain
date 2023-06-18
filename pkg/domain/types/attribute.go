package types

import "github.com/google/uuid"

type (
	AttrID        string
	AttrName      string
	AttrValue     any
	AttributeType string
)

func NewAttrID() AttrID {
	return AttrID(uuid.NewString())
}

const (
	IPAddr     AttributeType = "ipaddr"
	DomainName AttributeType = "domain"
	FileSha256 AttributeType = "file.sha256"
	FileSha512 AttributeType = "file.sha512"
	MarkDown   AttributeType = "markdown"
)
