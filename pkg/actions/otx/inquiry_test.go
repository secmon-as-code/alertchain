package otx

import (
	"context"
	"testing"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/gotx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type dummyClient struct {
	getIPv4General func(ctx context.Context, req *gotx.GetIPv4Request) (*gotx.GetIPv4GeneralResponse, error)
}

func (x *dummyClient) GetIPv4General(ctx context.Context, req *gotx.GetIPv4Request) (*gotx.GetIPv4GeneralResponse, error) {
	return x.getIPv4General(ctx, req)
}

func TestInquiry(t *testing.T) {
	action, err := NewInquiry(model.ActionConfig{
		"api_key": "xxx",
	})
	require.NoError(t, err)

	client, ok := action.(*Inquiry)
	require.True(t, ok)
	client.client = &dummyClient{
		getIPv4General: func(ctx context.Context, req *gotx.GetIPv4Request) (*gotx.GetIPv4GeneralResponse, error) {
			assert.Equal(t, "10.1.2.3", req.IPAddr)
			return &gotx.GetIPv4GeneralResponse{
				PulseInfo: gotx.PulseInfo{
					Pulses: []gotx.Pulse{
						{
							ID:      "abc123",
							Name:    "blue",
							Created: "2022-01-29T12:38:17.828000",
						},
						{
							ID:      "xyz098",
							Name:    "orange",
							Created: "2022-01-29T12:38:17.828000",
						},
					},
				},
			}, nil
		},
	}

	cr, err := action.Run(types.NewContext(), &model.Alert{
		Attributes: model.Attributes{
			{
				Type:  types.AttrNoType,
				Value: "xxx.xxx.xxx.xxx",
			},
			{
				Type:  types.AttrIPAddr,
				Value: "10.1.2.3",
			},
		},
	})

	require.NoError(t, err)
	require.NotNil(t, cr)
	assert.Len(t, cr.AddingAnnotations(), 2)
}
