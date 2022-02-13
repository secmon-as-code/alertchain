package model_test

import (
	"testing"
	"time"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/stretchr/testify/assert"
)

func TestAlertValidate(t *testing.T) {
	makeTestAlert := func() *model.Alert {
		alert := model.NewAlert(&model.Alert{
			Title:     "t",
			Detector:  "d",
			CreatedAt: time.Now(),
		})
		return alert
	}

	t.Run("pass if required fields are filled", func(t *testing.T) {
		alert := makeTestAlert()
		assert.NoError(t, alert.Validate())
	})

	t.Run("fail if Title is not set", func(t *testing.T) {
		alert := makeTestAlert()
		alert.Title = ""
		assert.ErrorIs(t, alert.Validate(), types.ErrInvalidInput)
	})

	t.Run("fail if Detector is not set", func(t *testing.T) {
		alert := makeTestAlert()
		alert.Detector = ""
		assert.ErrorIs(t, alert.Validate(), types.ErrInvalidInput)
	})

	t.Run("fail if Severity is not set", func(t *testing.T) {
		alert := makeTestAlert()
		alert.Severity = ""
		assert.ErrorIs(t, alert.Validate(), types.ErrInvalidInput)
	})

	t.Run("fail if id is not set", func(t *testing.T) {
		alert := &model.Alert{
			Title:     "t",
			Detector:  "d",
			Severity:  types.SevUnclassified,
			CreatedAt: time.Now(),
		}
		assert.ErrorIs(t, alert.Validate(), types.ErrInvalidInput)
	})

	t.Run("fail if CreatedAt is not set", func(t *testing.T) {
		alert := makeTestAlert()
		empty := model.NewAlert(nil)
		alert.CreatedAt = empty.CreatedAt
		assert.ErrorIs(t, alert.Validate(), types.ErrInvalidInput)
	})
}
