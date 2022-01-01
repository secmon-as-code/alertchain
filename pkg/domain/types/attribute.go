package types

import "github.com/m-mizutani/goerr"

type AttrType string

const (
	AttrIPAddr   AttrType = "ipaddr"
	AttrDomain   AttrType = "domain"
	AttrPort     AttrType = "port"
	AttrUserID   AttrType = "user_id"
	AttrEmail    AttrType = "email"
	AttrSha256   AttrType = "sha256"
	AttrFilePath AttrType = "filepath"
	AttrURL      AttrType = "url"
	AttrNoType   AttrType = ""
)

func (x AttrType) IsValid() error {
	switch x {
	case AttrIPAddr,
		AttrDomain,
		AttrPort,
		AttrUserID,
		AttrEmail,
		AttrSha256,
		AttrFilePath,
		AttrURL,
		AttrNoType:
		return nil
	}
	return goerr.Wrap(ErrInvalidInput, "invalid attribute type")
}

type AttrContext string

const (
	CtxRemote AttrContext = "remote"
	CtxLocal  AttrContext = "local"
	CtxServer AttrContext = "server"
	CtxClient AttrContext = "client"
)

func (x AttrContext) IsValid() error {
	switch x {
	case CtxRemote, CtxLocal, CtxServer, CtxClient:
		return nil
	}
	return goerr.Wrap(ErrInvalidInput, "invalid attribute context")
}
