package db

import "github.com/m-mizutani/alertchain/pkg/infra/ent"

func (x *Client) InjectClient(client *ent.Client) {
	x.client = client
}
