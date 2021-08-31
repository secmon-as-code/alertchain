package main

import (
	"fmt"

	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/types"
)

type BanIPAddr struct {
}

func (x *BanIPAddr) Name() string { return "BAN IP address" }
func (x *BanIPAddr) Executable(attr *alertchain.Attribute) bool {
	return attr.Type == types.AttrIPAddr
}

func (x *BanIPAddr) Execute(ctx *types.Context, attr *alertchain.Attribute) error {
	// Send request to firewall, WAF or etc.
	w := ctx.Writer()
	fmt.Fprintf(w, "Done\n")
	return nil
}
