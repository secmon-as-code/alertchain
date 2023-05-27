package lambda_test

import (
	_ "embed"

	"context"
	"encoding/json"
	"testing"

	"github.com/m-mizutani/alertchain/pkg/controller/lambda"
	"github.com/m-mizutani/gt"
)

//go:embed testdata/functional_url.json
var functionalURLData []byte

//go:embed testdata/sns_over_sqs.json
var snsOverSQSData []byte

func TestFunctionalURL(t *testing.T) {
	f := lambda.New(lambda.NewFunctionalURLHandler())

	var data map[string]any
	gt.NoError(t, json.Unmarshal(functionalURLData, &data))
	d := gt.R1(f(context.Background(), data)).NoError(t)
	gt.Value(t, d).Equal("OK")
}

func TestSNSoverSQS(t *testing.T) {
	f := lambda.New(lambda.NewSQSHandler("guardduty", lambda.WithDecodeSNS()))

	var data map[string]any
	gt.NoError(t, json.Unmarshal(snsOverSQSData, &data))
	d := gt.R1(f(context.Background(), data)).NoError(t)
	gt.Value(t, d).Equal("OK")
}
