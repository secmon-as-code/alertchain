package firestore_test

import (
	"testing"
	"time"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra/firestore"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/gots/ptr"
	"github.com/m-mizutani/gt"
)

func TestWorkflow(t *testing.T) {
	var (
		projectID  string
		databaseID string
	)

	if err := utils.LoadEnv(
		utils.EnvDef("TEST_FIRESTORE_PROJECT_ID", &projectID),
		utils.EnvDef("TEST_FIRESTORE_DATABASE_ID", &databaseID),
	); err != nil {
		t.Skipf("Skip test due to missing env: %v", err)
	}

	ctx := model.NewContext()
	now := time.Now()
	client := gt.R1(firestore.New(ctx, projectID, databaseID)).NoError(t)

	workflow0 := model.WorkflowRecord{
		ID:        types.NewWorkflowID(),
		CreatedAt: now.Add(-time.Second),
	}
	workflow1 := model.WorkflowRecord{
		ID:        types.NewWorkflowID(),
		CreatedAt: now,
		Alert: &model.AlertRecord{
			ID:        types.NewAlertID(),
			CreatedAt: now,
			Source:    "test",
			Title:     "testing",
			InitAttrs: []*model.AttributeRecord{
				{Key: "key1", Value: "value1"},
				{Key: "key2", Value: "value2"},
			},
			Refs: []*model.ReferenceRecord{
				{
					Title: ptr.To("ref1"),
					URL:   ptr.To("https://example.com"),
				},
			},
		},
	}
	workflow2 := model.WorkflowRecord{
		ID:        types.NewWorkflowID(),
		CreatedAt: now.Add(time.Second),
	}

	// Test PutWorkflow method
	gt.NoError(t, client.PutWorkflow(ctx, workflow0))
	gt.NoError(t, client.PutWorkflow(ctx, workflow1))
	gt.NoError(t, client.PutWorkflow(ctx, workflow2))

	// Test GetWorkflows method with offset and limit
	workflows := gt.R1(client.GetWorkflows(ctx, 1, 1)).NoError(t)
	gt.A(t, workflows).Length(1).At(0, func(t testing.TB, v model.WorkflowRecord) {
		gt.Equal(t, v.ID, workflow1.ID)
		gt.V(t, v.Alert).Must().NotNil()
		gt.Equal(t, v.Alert.Title, "testing")
		gt.Equal(t, v.Alert.InitAttrs[0].Key, "key1")
		gt.Equal(t, v.Alert.InitAttrs[0].Value, "value1")
		gt.Equal(t, *v.Alert.Refs[0].Title, "ref1")
	})
}
