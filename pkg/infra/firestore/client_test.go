package firestore_test

import (
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/infra/firestore"
	"github.com/m-mizutani/gots/ptr"
	"github.com/m-mizutani/gt"
)

func TestWorkflow(t *testing.T) {
	projectID, ok := os.LookupEnv("TEST_FIRESTORE_PROJECT_ID")
	if !ok {
		t.Skip("TEST_FIRESTORE_PROJECT_ID is not set")
	}
	collection, ok := os.LookupEnv("TEST_FIRESTORE_COLLECTION")
	if !ok {
		t.Skip("TEST_FIRESTORE_COLLECTION is not set")
	}

	ctx := model.NewContext()
	now := time.Now()
	client := gt.R1(firestore.New(ctx, projectID, collection)).NoError(t)

	workflow1 := model.WorkflowRecord{
		ID:        uuid.NewString(),
		CreatedAt: now,
		Alert: &model.AlertRecord{
			ID:        uuid.NewString(),
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
	workflow2 := model.WorkflowRecord{ID: uuid.NewString(), CreatedAt: now.Add(time.Second)}

	// Test PutWorkflow method
	gt.NoError(t, client.PutWorkflow(ctx, workflow1))
	gt.NoError(t, client.PutWorkflow(ctx, workflow2))

	// Test GetWorkflows method with offset and limit
	workflows := gt.R1(client.GetWorkflows(ctx, 1, 1)).NoError(t)
	gt.A(t, workflows).Length(1).At(0, func(t testing.TB, v model.WorkflowRecord) {
		gt.Equal(t, v.ID, workflow1.ID)
		gt.Equal(t, v.Alert.Title, "testing")
		gt.Equal(t, v.Alert.InitAttrs[0].Key, "key1")
		gt.Equal(t, v.Alert.InitAttrs[0].Value, "value1")
		gt.Equal(t, *v.Alert.Refs[0].Title, "ref1")
	})
}
