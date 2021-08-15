package types

type AttrType string

const (
	AttrIPAddr   AttrType = "ipaddr"
	AttrDomain   AttrType = "domain"
	AttrPort     AttrType = "port"
	AttrSha256   AttrType = "sha256"
	AttrFilePath AttrType = "filepath"
	AttrURL      AttrType = "url"
)

type AttrContext string

const (
	CtxRemote AttrContext = "remote"
	CtxLocal  AttrContext = "local"
	CtxServer AttrContext = "server"
	CtxClient AttrContext = "client"
)
