package otx

import (
	"context"
	"fmt"
	"time"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/gotx"
)

type otxClient interface {
	GetIPv4General(ctx context.Context, req *gotx.GetIPv4Request) (*gotx.GetIPv4GeneralResponse, error)
}

type Inquiry struct {
	client otxClient
}

const InquiryID = "otx-inquiry"

func NewInquiry(config model.ActionConfig) (model.Action, error) {
	newErr := func(msg string) error {
		return goerr.Wrap(types.ErrInvalidActionConfig, fmt.Sprintf("%s for %s", msg, InquiryID))
	}

	var apiKey string
	if v, ok := config["api_key"].(string); !ok {
		return nil, newErr("api_key is not set or invalid type")
	} else {
		apiKey = v
	}

	client, err := gotx.New(apiKey)
	if err != nil {
		return nil, err
	}

	return &Inquiry{
		client: client,
	}, nil
}

func (x *Inquiry) Run(ctx *types.Context, alert *model.Alert, args ...*model.Attribute) (*model.ChangeRequest, error) {
	cr := &model.ChangeRequest{}
	const pulseLimit = 5

	for _, attr := range alert.Attributes.FindByType(types.AttrIPAddr) {
		resp, err := x.client.GetIPv4General(ctx, &gotx.GetIPv4Request{
			IPAddr: attr.Value,
		})
		if err != nil {
			return nil, err
		}

		if len(resp.PulseInfo.Pulses) == 0 {
			continue
		}

		for i := 0; i < pulseLimit && i < len(resp.PulseInfo.Pulses); i++ {
			pulse := resp.PulseInfo.Pulses[i]
			ann := &model.Annotation{
				Source: "OTX",
				Name:   "Pulse",
				Value:  pulse.Name,
				URI:    "https://otx.alienvault.com/pulse/" + pulse.ID,
			}
			if dt, err := time.Parse("2006-01-02T15:04:05.000000", pulse.Created); err == nil {
				ann.Timestamp = &dt
			}
			cr.AddAnnotation(attr, ann)
		}
	}

	return cr, nil
}
