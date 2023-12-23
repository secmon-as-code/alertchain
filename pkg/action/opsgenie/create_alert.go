package opsgenie

import (
	"fmt"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/goerr"
	og_alert "github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
)

type Responder struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	UserName string `json:"username"`
	Type     string `json:"type"`
}

func CreateAlert(ctx *model.Context, alert model.Alert, args model.ActionArgs) (any, error) {
	var (
		apiKey     string
		responders []Responder
	)

	if err := args.Parse(
		model.ArgDef("secret_api_key", &apiKey),
		model.ArgDef("responder_teams", &responders, model.ArgOptional()),
	); err != nil {
		return nil, err
	}

	c, err := og_alert.NewClient(&client.Config{
		ApiKey: apiKey,
	})
	if err != nil {
		return nil, goerr.Wrap(err, "Failed to create OpsGenie client")
	}

	req := &og_alert.CreateAlertRequest{
		Message:     alert.Title,
		Description: alert.Description,
		Alias:       alert.ID.String(),
		Source:      "alertchain",
		Details:     map[string]string{},
	}
	for _, attr := range alert.Attrs {
		req.Details[attr.Key.String()] = fmt.Sprintf("%+v", attr.Value)
	}

	for _, r := range responders {
		req.Responders = append(req.Responders, og_alert.Responder{
			Type:     og_alert.ResponderType(r.Type),
			Id:       r.ID,
			Name:     r.Name,
			Username: r.UserName,
		})
	}

	resp, err := c.Create(ctx, req)
	if err != nil {
		return nil, goerr.Wrap(err, "Failed to create OpsGenie alert")
	}

	return resp, nil
}
