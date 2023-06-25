package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

type AlertMetaData struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Source      string     `json:"source"`
	Attrs       Attributes `json:"attrs"`
}

func (x AlertMetaData) Copy() AlertMetaData {
	newMeta := AlertMetaData{
		Title:       x.Title,
		Description: x.Description,
		Source:      x.Source,
		Attrs:       x.Attrs.Copy(),
	}
	return newMeta
}

type Alert struct {
	AlertMetaData
	ID        types.AlertID `json:"id"`
	Schema    types.Schema  `json:"schema"`
	Data      any           `json:"data"`
	CreatedAt time.Time     `json:"created_at"`

	Raw string `json:"-"`
}

func (x Alert) Copy() (Alert, error) {
	newAlert := Alert{
		AlertMetaData: x.AlertMetaData.Copy(),
	}

	return newAlert, nil
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
		Raw:           encodeAlertData(data),
	}
	alert.AlertMetaData.Attrs = TidyAttributes(alert.AlertMetaData.Attrs)

	return alert
}
