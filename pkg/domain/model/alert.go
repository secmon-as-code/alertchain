package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

type Parameter struct {
	Key   string
	Value string
}

type AlertMetaData struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Source      string            `json:"source"`
	Params      []types.Parameter `json:"params"`
}

type Alert struct {
	AlertMetaData
	ID         types.AlertID `json:"id"`
	Schema     types.Schema  `json:"schema"`
	Data       any           `json:"-"`
	CreatedAt  time.Time     `json:"created_at"`
	References []Reference   `json:"reference"`

	Raw string `json:"-"`
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
	return Alert{
		AlertMetaData: meta,
		ID:            types.NewAlertID(),
		Schema:        schema,
		Data:          data,
		CreatedAt:     time.Now(),
		Raw:           encodeAlertData(data),
	}
}

type Reference struct {
	types.Parameter
	Source types.EnricherID `json:"source"`
	Data   any              `json:"data"`
}
