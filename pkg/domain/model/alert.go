package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/goerr"
)

type AlertMetaData struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Source      string     `json:"source"`
	Params      Parameters `json:"params"`
}

type Alert struct {
	AlertMetaData
	ID        types.AlertID `json:"id"`
	Schema    types.Schema  `json:"schema"`
	Data      any           `json:"-"`
	CreatedAt time.Time     `json:"created_at"`

	Raw string `json:"-"`
}

func (x Alert) Copy() (Alert, error) {
	raw, err := json.Marshal(x)
	if err != nil {
		return Alert{}, goerr.Wrap(err, "Failed to marshal alert")
	}

	var newAlert Alert
	if err := json.Unmarshal(raw, &newAlert); err != nil {
		return Alert{}, goerr.Wrap(err, "Failed to unmarshal alert")
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
	alert.AlertMetaData.Params = TidyParameters(alert.AlertMetaData.Params)

	return alert
}
