package incident_io_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/m-mizutani/alertchain/pkg/action/incident_io"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/gt"
)

func TestWithSink(t *testing.T) {
	alert := model.Alert{
		ID: types.NewAlertID(),
		AlertMetaData: model.AlertMetaData{
			Title:       "Test alert",
			Description: "This is test alert",
			Source:      "test_source",
			Attrs: []model.Attribute{
				{
					Key:   "key1",
					Value: "value1",
				},
			},
		},
		CreatedAt: time.Now(),
	}

	testCases := map[string]struct {
		args  model.ActionArgs
		isErr bool
		check func(t *testing.T, req *http.Request)
	}{
		"Create alert with minimal arguments": {
			args: model.ActionArgs{
				"secret_api_token":       "test_token",
				"alert_source_config_id": "test_config_id",
			},
			isErr: false,
			check: func(t *testing.T, req *http.Request) {
				gt.Equal(t, req.Method, "POST")
				gt.Equal(t, req.URL.String(), "https://api.incident.io/v2/alert_events/http/test_config_id")
				gt.Equal(t, req.Header.Get("Authorization"), "Bearer test_token")

				var reqBody incident_io.CreateAlertRequest
				gt.NoError(t, json.NewDecoder(req.Body).Decode(&reqBody))
				gt.Equal(t, reqBody.Status, "firing")
				gt.Equal(t, reqBody.Title, "Test alert")
				gt.Equal(t, reqBody.DeduplicationKey, alert.ID.String())
				gt.Equal(t, reqBody.Description, "This is test alert")
				gt.Equal(t, reqBody.MetaData, map[string]interface{}{
					"key1": "value1",
				})
				gt.Equal(t, reqBody.SourceURL, "")
			},
		},

		"Create alert with all arguments": {
			args: model.ActionArgs{
				"secret_api_token":       "test_token",
				"alert_source_config_id": "test_config_id",
				"status":                 "resolved",

				"title":             "title form args",
				"description":       "description from args",
				"deduplication_key": "deduplication_key",
				"metadata": map[string]interface{}{
					"key2": "value2",
				},
				"source_url": "https://example.com",
			},
			isErr: false,
			check: func(t *testing.T, req *http.Request) {
				gt.Equal(t, req.Method, "POST")
				gt.Equal(t, req.URL.String(), "https://api.incident.io/v2/alert_events/http/test_config_id")
				gt.Equal(t, req.Header.Get("Authorization"), "Bearer test_token")

				var reqBody incident_io.CreateAlertRequest
				gt.NoError(t, json.NewDecoder(req.Body).Decode(&reqBody))
				gt.Equal(t, reqBody.Status, "resolved")
				gt.Equal(t, reqBody.Title, "title form args")
				gt.Equal(t, reqBody.DeduplicationKey, "deduplication_key")
				gt.Equal(t, reqBody.Description, "description from args")
				gt.Equal(t, reqBody.MetaData, map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				})
				gt.Equal(t, reqBody.SourceURL, "https://example.com")
			},
		},

		"Create alert without api token": {
			args: model.ActionArgs{
				"alert_source_config_id": "test_config_id",
			},
			isErr: true,
		},
		"Create alert without alert source config ID": {
			args: model.ActionArgs{
				"secret_api_token": "test_token",
			},
			isErr: true,
		},

		"Create alert with invalid metadata": {
			args: model.ActionArgs{
				"secret_api_token":       "test_token",
				"alert_source_config_id": "test_config_id",
				"metadata": map[string]interface{}{
					"key2": make(chan int),
				},
			},
			isErr: true,
		},

		"Create alert with merging metadata": {
			args: model.ActionArgs{
				"secret_api_token":       "test_token",
				"alert_source_config_id": "test_config_id",
				"metadata": map[string]interface{}{
					"key1": "valueX",
				},
			},
			isErr: false,
			check: func(t *testing.T, req *http.Request) {
				var reqBody incident_io.CreateAlertRequest
				gt.NoError(t, json.NewDecoder(req.Body).Decode(&reqBody))
				gt.Equal(t, reqBody.MetaData, map[string]interface{}{
					"key1": "valueX",
				})
			},
		},
	}

	ctx := model.NewContext(model.WithTest())
	for name, tc := range testCases {
		incident_io.Sink.Reset()

		t.Run(name, func(t *testing.T) {
			_, err := incident_io.CreateAlert(ctx, alert, tc.args)
			if tc.isErr {
				gt.True(t, err != nil)
				return
			}
			gt.V(t, err).Nil()
			if tc.check != nil {
				tc.check(t, incident_io.Sink.Requests[0])
			}
		})
	}
}

func TestCreateAlert(t *testing.T) {
	var (
		configID string
		apiToken string
	)
	if err := utils.LoadEnv(
		utils.EnvDef("TEST_INCIDENT_IO_ALERT_SOURCE_CONFIG_ID", &configID),
		utils.EnvDef("TEST_INCIDENT_IO_API_TOKEN", &apiToken),
	); err != nil {
		t.Skip("Skipping test because TEST_INCIDENT_IO_ALERT_SOURCE_CONFIG_ID and TEST_INCIDENT_IO_API_TOKEN are not set")
	}

	alert := model.Alert{
		ID: types.NewAlertID(),
		AlertMetaData: model.AlertMetaData{
			Title:       "Test alert",
			Description: "This is test alert",
			Source:      "test_source",
			Attrs: []model.Attribute{
				{
					Key:   "key1",
					Value: "value1",
				},
			},
		},
		CreatedAt: time.Now(),
	}

	args := model.ActionArgs{
		"secret_api_token":       apiToken,
		"alert_source_config_id": configID,
	}

	ctx := model.NewContext()
	resp := gt.R1(incident_io.CreateAlert(ctx, alert, args)).NoError(t)
	gt.V(t, resp).NotNil()
}
