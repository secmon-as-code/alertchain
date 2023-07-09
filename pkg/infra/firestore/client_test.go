package firestore_test

import (
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra/firestore"
	"github.com/m-mizutani/gt"
)

func setupClient(t *testing.T, ctx *model.Context) *firestore.Client {
	projectID, ok := os.LookupEnv("TEST_FIRESTORE_PROJECT_ID")
	if !ok {
		t.Skip("TEST_FIRESTORE_PROJECT_ID not set")
	}

	collection, ok := os.LookupEnv("TEST_FIRESTORE_COLLECTION")
	if !ok {
		t.Skip("TEST_FIRESTORE_COLLECTION not set")
	}

	client := gt.R1(firestore.New(ctx, projectID, collection)).NoError(t)

	return client
}

func TestPutGet(t *testing.T) {
	ctx := model.NewContext()
	client := setupClient(t, ctx)
	defer client.Close()

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

func TestLock(t *testing.T) {
	client := setupClient(t, model.NewContext())

	ns := types.Namespace(uuid.New().String())
	taskNum := 4
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

			println("lock")
			records <- &record{time.Now(), "lock"}
			time.Sleep(time.Duration(rand.Intn(300)) * time.Millisecond)
			println("unlock")
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

func TestLockExpires(t *testing.T) {
	client := setupClient(t, model.NewContext())
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
