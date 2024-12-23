package bigquery_test

import (
	"context"
	"testing"
	"time"

	"github.com/m-mizutani/gt"
	"github.com/secmon-lab/alertchain/action/bigquery"
	"github.com/secmon-lab/alertchain/pkg/domain/model"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
	"github.com/secmon-lab/alertchain/pkg/utils"
)

func TestInsertDataIntegration(t *testing.T) {
	var (
		projectID string
		datasetID string
		tableID   string
	)

	if err := utils.LoadEnv(
		utils.EnvDef("TEST_BIGQUERY_PROJECT_ID", &projectID),
		utils.EnvDef("TEST_BIGQUERY_DATASET_ID", &datasetID),
		utils.EnvDef("TEST_BIGQUERY_DATA_TABLE_ID", &tableID),
	); err != nil {
		t.Skipf("Skip test due to missing env: %v", err)
	}

	ctx := context.Background()
	args := model.ActionArgs{
		"project_id": projectID,
		"dataset_id": datasetID,
		"table_id":   tableID,
		"data": map[string]interface{}{
			"color":  "blue",
			"number": 5,
			"nested": map[string]interface{}{
				"foo": "bar",
				"bar": "baz",
			},
		},
	}

	ret := gt.R1(bigquery.InsertData(ctx, model.Alert{
		ID: types.NewAlertID(),
	}, args)).NoError(t)
	gt.V(t, ret).Nil()
}

func TestInsertAlertIntegration(t *testing.T) {
	var (
		projectID string
		datasetID string
		tableID   string
	)

	if err := utils.LoadEnv(
		utils.EnvDef("TEST_BIGQUERY_PROJECT_ID", &projectID),
		utils.EnvDef("TEST_BIGQUERY_DATASET_ID", &datasetID),
		utils.EnvDef("TEST_BIGQUERY_ALERT_TABLE_ID", &tableID),
	); err != nil {
		t.Skipf("Skip test due to missing env: %v", err)
	}

	ctx := context.Background()
	args := model.ActionArgs{
		"project_id": projectID,
		"dataset_id": datasetID,
		"table_id":   tableID,
	}
	alert := model.Alert{
		ID: types.NewAlertID(),
		AlertMetaData: model.AlertMetaData{
			Source:      "test_source",
			Namespace:   "test_namespace",
			Title:       "test alert",
			Description: "test description",
			Attrs: model.Attributes{
				{
					Key:   "color",
					Value: "blue",
				},
			},
			Refs: model.References{
				{
					Title: "test reference",
					URL:   "https://example.com",
				},
			},
		},
		Schema:    "test_schema",
		CreatedAt: time.Now(),
		Data: map[string]interface{}{
			"color":  "blue",
			"number": 5,
			"nested": map[string]interface{}{
				"foo": "bar",
			},
		},
		Raw: `{ "color": "blue", "number": 5, "nested": { "foo": "bar" } }`,
	}

	ret := gt.R1(bigquery.InsertAlert(ctx, alert, args)).NoError(t)
	gt.V(t, ret).Nil()
}
