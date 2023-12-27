package infra_test

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra/firestore"
	"github.com/m-mizutani/alertchain/pkg/infra/memory"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/gt"
)

func TestMemory(t *testing.T) {
	testClient(t, memory.New())
}

func TestFirestore(t *testing.T) {
	var (
		projectID  string
		collection string
	)

	if err := utils.LoadEnv(
		utils.Env("TEST_FIRESTORE_PROJECT_ID", &projectID),
		utils.Env("TEST_FIRESTORE_COLLECTION_PREFIX", &collection),
	); err != nil {
		t.Skipf("Skip test due to missing env: %v", err)
	}

	ctx := model.NewContext()
	client := gt.R1(firestore.New(ctx, projectID, collection)).NoError(t)

	testClient(t, client)
}

func testClient(t *testing.T, client interfaces.Database) {
	t.Run("PutGet", func(t *testing.T) {
		testPutGet(t, client)
	})
	t.Run("Lock", func(t *testing.T) {
		testLock(t, client)
	})
	t.Run("LockExpire", func(t *testing.T) {
		testLockExpires(t, client)
	})
	t.Run("Workflow", func(t *testing.T) {
		testWorkflow(t, client)
	})
}

func testPutGet(t *testing.T, client interfaces.Database) {
	ctx := model.NewContext()

	attrs1 := model.Attributes{
		{
			ID:    types.NewAttrID(),
			Key:   "key1",
			Value: "value1",
		},
		{
			ID:    types.NewAttrID(),
			Key:   "key2",
			Value: "value2",
		},
	}
	attrs2 := model.Attributes{
		{
			ID:    types.NewAttrID(),
			Key:   "key3",
			Value: "value3",
		},
	}

	ns1 := types.Namespace(uuid.New().String())
	ns2 := types.Namespace(uuid.New().String())
	ns3 := types.Namespace(uuid.New().String())
	gt.NoError(t, client.PutAttrs(ctx, ns1, attrs1))
	gt.NoError(t, client.PutAttrs(ctx, ns2, attrs2))

	t.Run("GetAttrs from ns1", func(t *testing.T) {
		resp := gt.R1(client.GetAttrs(ctx, ns1)).NoError(t)
		gt.A(t, resp).Length(2).At(0, func(t testing.TB, v model.Attribute) {
			gt.V(t, v.Key).In("key1", "key2")
		})

		t.Run("Update the attribute", func(t *testing.T) {
			gt.NoError(t, client.PutAttrs(ctx, ns1, model.Attributes{
				{
					ID:    attrs1[0].ID,
					Key:   "keyA", // should not be updated
					Value: "valueA",
				},
			}))

			check := func(t testing.TB, v model.Attribute) {
				gt.V(t, v.Key).In("key1", "key2").NotEqual("keyA")
				gt.V(t, v.Value).In("valueA", "value2").NotEqual("value1")
			}

			resp := gt.R1(client.GetAttrs(ctx, ns1)).NoError(t)
			gt.A(t, resp).Length(2).
				At(0, check).
				At(1, check)
		})
	})

	t.Run("GetAttrs from ns2", func(t *testing.T) {
		resp := gt.R1(client.GetAttrs(ctx, ns2)).NoError(t)
		gt.A(t, resp).Length(1).At(0, func(t testing.TB, v model.Attribute) {
			gt.V(t, v.Key).In("key3").NotEqual("key1").NotEqual("key2")
		})
	})

	t.Run("GetAttrs from ns3", func(t *testing.T) {
		resp := gt.R1(client.GetAttrs(ctx, ns3)).NoError(t)
		gt.A(t, resp).Length(0)
	})
}

func testLock(t *testing.T, client interfaces.Database) {
	ns := types.Namespace(uuid.New().String())
	taskNum := 10
	type record struct {
		ts    time.Time
		label string
	}
	records := make(chan *record, taskNum*2)
	wg := sync.WaitGroup{}

	ctx := model.NewContext(
		model.WithAlert(
			model.NewAlert(model.AlertMetaData{
				Title: "test",
			}, types.Schema("test"), "test",
			),
		),
	)
	now := time.Now()
	for i := 0; i < taskNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			gt.NoError(t, client.Lock(ctx, ns, now.Add(10*time.Second)))

			t.Log("lock")
			records <- &record{time.Now(), "lock"}
			time.Sleep(time.Duration(rand.Intn(300)) * time.Millisecond)
			t.Log("unlock")
			records <- &record{time.Now(), "unlock"}
			if err := client.Unlock(ctx, ns); err != nil {
				t.Error(err)
			}
		}()
	}

	wg.Wait()
	close(records)
	var results []*record
	for r := range records {
		if len(results)%2 == 0 {
			gt.V(t, r.label).Equal("lock")
		} else {
			gt.V(t, r.label).Equal("unlock")
		}

		if len(results) > 0 {
			gt.B(t, results[len(results)-1].ts.Before(r.ts)).True()
		}
		results = append(results, r)
	}
	gt.Array(t, results).Length(taskNum * 2)
}

func testLockExpires(t *testing.T, client interfaces.Database) {
	ns := types.Namespace(uuid.New().String())

	ctx := model.NewContext(
		model.WithAlert(
			model.NewAlert(model.AlertMetaData{
				Title: "test",
			}, types.Schema("test"), "test",
			),
		),
	)

	// Lock with 1 second
	gt.NoError(t, client.Lock(ctx, ns, time.Now().Add(2*time.Second)))
	// Next lock can be done after 1 second without unlock
	gt.NoError(t, client.Lock(ctx, ns, time.Now().Add(100*time.Millisecond)))
}

func testWorkflow(t *testing.T, client interfaces.Database) {
	now := time.Now()
	workflows := []model.WorkflowRecord{
		{
			ID:        types.NewWorkflowID(),
			CreatedAt: now,
		},
		{
			ID:        types.NewWorkflowID(),
			CreatedAt: now.Add(1 * time.Second),
		},
		{
			ID:        types.NewWorkflowID(),
			CreatedAt: now.Add(2 * time.Second),
			Alert: &model.AlertRecord{
				ID:        types.NewAlertID(),
				CreatedAt: now,
				Source:    "test",
				Title:     "testing",
			},
		},
		{
			ID:        types.NewWorkflowID(),
			CreatedAt: now.Add(3 * time.Second),
		},
		{
			ID:        types.NewWorkflowID(),
			CreatedAt: now.Add(4 * time.Second),
		},
	}

	ctx := model.NewContext()
	for _, wf := range workflows {
		gt.NoError(t, client.PutWorkflow(ctx, wf))
	}

	t.Run("GetWorkflows", func(t *testing.T) {
		resp := gt.R1(client.GetWorkflows(ctx, 1, 2)).NoError(t)
		gt.A(t, resp).Length(2).At(0, func(t testing.TB, v model.WorkflowRecord) {
			gt.V(t, v.ID).Equal(workflows[3].ID)
		}).At(1, func(t testing.TB, v model.WorkflowRecord) {
			gt.V(t, v.ID).Equal(workflows[2].ID)
		})
	})

	t.Run("GetWorkflow by ID", func(t *testing.T) {
		resp := gt.R1(client.GetWorkflow(ctx, types.WorkflowID(workflows[2].ID))).NoError(t)
		gt.V(t, resp.Alert.ID).Equal(workflows[2].Alert.ID)
	})
}
