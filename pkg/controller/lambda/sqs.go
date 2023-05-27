package lambda

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/goerr"
)

type sqsConfig struct {
	decodeSNS bool
}

type SQSOption func(*sqsConfig)

func WithDecodeSNS() SQSOption {
	return func(cfg *sqsConfig) {
		cfg.decodeSNS = true
	}
}

// NewSQSHandler is a handler for SQS event.
func NewSQSHandler(schema types.Schema, options ...SQSOption) func(context.Context, any, Callback) error {
	var cfg sqsConfig
	for _, opt := range options {
		opt(&cfg)
	}

	return func(ctx context.Context, data any, cb Callback) error {
		var event events.SQSEvent
		if err := remapEvent(data, &event); err != nil {
			return goerr.Wrap(err, "fail to remap event")
		}
		if len(event.Records) == 0 {
			return goerr.Wrap(types.ErrInvalidLambdaRequest, "Event is not SQSEvent")
		}

		for _, record := range event.Records {
			body := record.Body

			if cfg.decodeSNS {
				var snsEvent events.SNSEntity
				if err := json.Unmarshal([]byte(record.Body), &snsEvent); err != nil {
					return goerr.Wrap(err, "fail to decode SNS event")
				}
				body = snsEvent.Message
			}

			var data any
			if err := json.Unmarshal([]byte(body), &data); err != nil {
				return goerr.Wrap(err, "fail to encode body of SQS event").With("body", body)
			}

			if err := cb(ctx, schema, data); err != nil {
				return err
			}
		}

		return nil
	}
}
