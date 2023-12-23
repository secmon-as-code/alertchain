package types

import (
	"github.com/google/uuid"
)

type (
	AttrID    string
	AttrKey   string
	AttrValue any
	AttrType  string
	AttrTTL   int64
)

func NewAttrID() AttrID {
	return AttrID(uuid.NewString())
}

const (
	IPAddr     AttrType = "ipaddr"
	DomainName AttrType = "domain"
	FileSha256 AttrType = "file.sha256"
	FileSha512 AttrType = "file.sha512"
	MarkDown   AttrType = "markdown"
)

func (x AttrID) String() string  { return string(x) }
func (x AttrKey) String() string { return string(x) }
