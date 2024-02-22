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
		databaseID string
	)

	if err := utils.LoadEnv(
		utils.EnvDef("TEST_FIRESTORE_PROJECT_ID", &projectID),
		utils.EnvDef("TEST_FIRESTORE_DATABASE_ID", &databaseID),
	); err != nil {
		t.Skipf("Skip test due to missing env: %v", err)
	}

	ctx := model.NewContext()
	client := gt.R1(firestore.New(ctx, projectID, databaseID)).NoError(t)

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
	t.Run("Action", func(t *testing.T) {
		testAction(t, client)
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

func testAction(t *testing.T, db interfaces.Database) {
	ctx := model.NewContext()
	workflow1 := model.WorkflowRecord{
		ID:        types.NewWorkflowID(),
		CreatedAt: time.Now(),
	}
	workflow2 := model.WorkflowRecord{
		ID:        types.NewWorkflowID(),
		CreatedAt: time.Now(),
	}

	gt.NoError(t, db.PutWorkflow(ctx, workflow1))

	actions := []model.ActionRecord{
		{
			ID:         types.NewActionID(),
			WorkflowID: workflow1.ID,
			Seq:        0,
			Uses:       "test1",
			Args: []*model.ArgumentRecord{
				{
					Key:   "key1",
					Value: "value1",
				},
			},
		},
		{
			ID:         types.NewActionID(),
			WorkflowID: workflow1.ID,
			Seq:        1,
			Uses:       "test2",
			Args: []*model.ArgumentRecord{
				{
					Key:   "key2",
					Value: "value2",
				},
			},
		},
		{
			ID:         types.NewActionID(),
			WorkflowID: workflow1.ID,
			Seq:        2,
			Uses:       "test3",
			Args:       []*model.ArgumentRecord{},
			Next:       []*model.NextRecord{},
			StartedAt:  time.Now(),
			FinishedAt: time.Now(),
		},
		{
			ID:         types.NewActionID(),
			WorkflowID: workflow2.ID,
			Seq:        0,
			Uses:       "test4",
			Args:       []*model.ArgumentRecord{},
			Next:       []*model.NextRecord{},
			StartedAt:  time.Now(),
			FinishedAt: time.Now(),
		},
	}

	for _, action := range actions {
		gt.NoError(t, db.PutAction(ctx, action))
	}

	t.Run("GetAction", func(t *testing.T) {
		resp := gt.R1(db.GetAction(ctx, actions[1].ID)).NoError(t)
		gt.V(t, resp).Must().NotNil()
		gt.V(t, resp.Uses).Equal("test2")
	})

	t.Run("GetActionByWorkflowID", func(t *testing.T) {
		resp := gt.R1(db.GetActionByWorkflowID(ctx, types.WorkflowID(workflow1.ID))).NoError(t)
		gt.A(t, resp).Length(3).
			MatchThen(func(v model.ActionRecord) bool {
				return v.ID == actions[0].ID
			}, func(t testing.TB, v model.ActionRecord) {
				gt.V(t, v.Uses).Equal("test1")
			}).
			MatchThen(func(v model.ActionRecord) bool {
				return v.ID == actions[1].ID
			}, func(t testing.TB, v model.ActionRecord) {
				gt.V(t, v.Uses).Equal("test2")
			}).
			MatchThen(func(v model.ActionRecord) bool {
				return v.ID == actions[2].ID
			}, func(t testing.TB, v model.ActionRecord) {
				gt.V(t, v.Uses).Equal("test3")
			})
	})

	t.Run("GetActions", func(t *testing.T) {
		resp := gt.R1(db.GetActions(ctx, []types.ActionID{actions[0].ID, actions[3].ID})).NoError(t)
		gt.A(t, resp).Length(2).
			MatchThen(func(v model.ActionRecord) bool {
				return v.ID == actions[0].ID
			}, func(t testing.TB, v model.ActionRecord) {
				gt.V(t, v.Uses).Equal("test1")
				gt.V(t, v.WorkflowID).Equal(workflow1.ID)
			}).
			MatchThen(func(v model.ActionRecord) bool {
				return v.ID == actions[3].ID
			}, func(t testing.TB, v model.ActionRecord) {
				gt.V(t, v.Uses).Equal("test4")
				gt.V(t, v.WorkflowID).Equal(workflow2.ID)
			})
	})
}
