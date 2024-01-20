package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

type AlertMetaData struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Source      string          `json:"source"`
	Namespace   types.Namespace `json:"namespace"`
	Attrs       Attributes      `json:"attrs"`
	Refs        References      `json:"refs"`
}

func (x AlertMetaData) Copy() AlertMetaData {
	newMeta := AlertMetaData{
		Title:       x.Title,
		Description: x.Description,
		Source:      x.Source,
		Namespace:   x.Namespace,
		Attrs:       x.Attrs.Copy(),
		Refs:        x.Refs.Copy(),
	}
	return newMeta
}

type Alert struct {
	AlertMetaData
	ID        types.AlertID `json:"id"`
	Schema    types.Schema  `json:"schema"`
	Data      any           `json:"data"`
	CreatedAt time.Time     `json:"created_at"`

	// Raw is a JSON string of Data. The field will be redacted by masq because of verbosity.
	Raw string `json:"raw" masq:"quiet"`
}

func (x Alert) Copy() Alert {
	newAlert := Alert{
		AlertMetaData: x.AlertMetaData.Copy(),

		ID:        x.ID,
		Schema:    x.Schema,
		Data:      x.Data,
		CreatedAt: x.CreatedAt,

		Raw: x.Raw,
	}

	return newAlert
}

func encodeAlertData(a any) string {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(a); err != nil {
		return fmt.Sprintf("%v", a)
	}

	return buf.String()
}

func NewAlert(meta AlertMetaData, schema types.Schema, data any) Alert {
	alert := Alert{
		AlertMetaData: meta,
		ID:            types.NewAlertID(),
		Schema:        schema,
		Data:          data,
		CreatedAt:     time.Now(),

		Raw: encodeAlertData(data),
	}
	alert.AlertMetaData.Attrs = alert.AlertMetaData.Attrs.Tidy()

	return alert
}
