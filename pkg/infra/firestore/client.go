package firestore

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"math"
	"math/rand"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/goerr"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Client struct {
	client             *firestore.Client
	projectID          string
	attrCollection     string
	workflowCollection string
}

const (
	attrKeyPrefix     = "attr:"
	lockKeyPrefix     = "lock:"
	workflowKeyPrefix = "workflow:"
)

func hashNamespace(input types.Namespace) string {
	hash := sha512.New()
	hash.Write([]byte(input))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}

// GetAttrs implements interfaces.Database.
func (x *Client) GetAttrs(ctx *model.Context, ns types.Namespace) (model.Attributes, error) {
	key := attrKeyPrefix + hashNamespace(ns)
	docs, err := x.client.Collection(x.attrCollection).Doc(key).Collection("attributes").Documents(ctx).GetAll()
	if err != nil {
		return nil, goerr.Wrap(err, "failed to get attributes from firestore")
	}

	now := time.Now().UTC()
	var attrs model.Attributes
	for _, doc := range docs {
		if !doc.Exists() {
			continue
		}

		var attr attribute
		if err := doc.DataTo(&attr); err != nil {
			return nil, goerr.Wrap(err, "failed to unmarshal attribute from firestore")
		}
		if attr.ExpiresAt.Before(now) {
			continue
		}
		attrs = append(attrs, attr.Attribute)
	}

	return attrs, nil
}

// PutAttrs implements interfaces.Database.
func (x *Client) PutAttrs(ctx *model.Context, ns types.Namespace, attrs model.Attributes) error {
	err := x.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		key := attrKeyPrefix + hashNamespace(ns)
		collection := x.client.Collection(x.attrCollection).Doc(key).Collection("attributes")

		attrRefMap := map[types.AttrID]*firestore.DocumentRef{}
		for _, attr := range attrs {
			doc, err := collection.Doc(string(attr.ID)).Get(ctx)
			if err != nil {
				if status.Code(err) != codes.NotFound {
					return goerr.Wrap(err, "failed to get attributes from firestore")
				}
				continue
			}
			attrRefMap[attr.ID] = doc.Ref
		}

		now := time.Now().UTC()

		for _, base := range attrs {
			ttl := base.TTL
			if ttl == 0 {
				ttl = types.DefaultAttributeTTL
			}
			attr := attribute{
				Attribute: base,
				ExpiresAt: now.Add(time.Duration(ttl) * time.Second),
			}

			if ref, ok := attrRefMap[attr.ID]; ok {
				if err := tx.Set(ref, map[string]any{
					"value":      attr.Value,
					"expires_at": attr.ExpiresAt,
				}, firestore.MergeAll); err != nil {
					return goerr.Wrap(err, "failed to unmarshal attribute from firebase")
				}
			} else {
				ref := collection.Doc(string(attr.ID))
				if err := tx.Create(ref, attr); err != nil {
					return goerr.Wrap(err, "failed to create attribute")
				}
			}
		}

		return nil
	})
	if err != nil {
		return goerr.Wrap(err, "failed firestore transaction")
	}

	return nil
}

func (x *Client) PutWorkflow(ctx *model.Context, workflow model.WorkflowRecord) error {
	key := workflowKeyPrefix + workflow.ID

	if _, err := x.client.Collection(x.workflowCollection).Doc(key).Set(ctx, workflow); err != nil {
		return goerr.Wrap(err, "failed to put workflow")
	}
	return nil
}

func (x *Client) GetWorkflows(ctx *model.Context, offset, limit int) ([]model.WorkflowRecord, error) {
	var workflows []model.WorkflowRecord
	iter := x.client.Collection(x.workflowCollection).
		OrderBy("CreatedAt", firestore.Desc).
		Offset(offset).
		Limit(limit).
		Documents(ctx)

	for {
		doc, err := iter.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				return workflows, nil
			}
			return nil, goerr.Wrap(err, "failed to get workflow")
		}

		var workflow model.WorkflowRecord
		if err := doc.DataTo(&workflow); err != nil {
			return nil, goerr.Wrap(err, "failed to unmarshal workflow")
		}
		workflows = append(workflows, workflow)
	}
}

func (x *Client) GetWorkflow(ctx *model.Context, id types.WorkflowID) (*model.WorkflowRecord, error) {
	key := workflowKeyPrefix + id.String()
	doc, err := x.client.Collection(x.workflowCollection).Doc(key).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, goerr.Wrap(err, "failed to get workflow")
	}

	var workflow model.WorkflowRecord
	if err := doc.DataTo(&workflow); err != nil {
		return nil, goerr.Wrap(err, "failed to unmarshal workflow")
	}

	return &workflow, nil
}

type attribute struct {
	model.Attribute
	ExpiresAt time.Time `firestore:"expires_at"`
}

type lock struct {
	AlertID   types.AlertID `firestore:"alert_id"`
	ExpiresAt time.Time     `firestore:"expires_at"`
}

const (
	expBackOffMaxDelay  float64 = 10000
	expBackOffBaseDelay float64 = 50
)

func exponentialBackoff(attempt int) time.Duration {
	delay := expBackOffBaseDelay * math.Pow(2, float64(attempt))
	if delay > expBackOffMaxDelay {
		delay = expBackOffMaxDelay
	}

	// #nosec
	jitter := rand.Float64() * delay / 2
	backoff := delay + jitter

	return time.Duration(backoff) * time.Millisecond
}

// Lock implements interfaces.Database.
func (x *Client) Lock(ctx *model.Context, ns types.Namespace, timeout time.Time) error {
	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			return goerr.Wrap(ctx.Err(), "context is done")
		default:
			if err := x.tryLock(ctx, ns, timeout); err != nil {
				if !errors.Is(err, errLockFailed) {
					return goerr.Wrap(err, "failed to lock")
				}
			} else {
				return nil
			}
		}

		wait := exponentialBackoff(i)

		select {
		case <-ctx.Done():
			return goerr.Wrap(ctx.Err(), "context is done")
		case <-time.After(wait):
			// wait
		}
	}
}

var (
	errLockFailed = goerr.New("failed to lock")
)

func (x *Client) tryLock(ctx *model.Context, ns types.Namespace, timeout time.Time) error {
	key := lockKeyPrefix + hashNamespace(ns)
	now := time.Now().UTC()
	alert := ctx.Alert()

	err := x.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		var doc *firestore.DocumentSnapshot
		resp, err := tx.Get(x.client.Collection(x.attrCollection).Doc(key))
		if err != nil {
			if status.Code(err) != codes.NotFound {
				return goerr.Wrap(err, "failed to get attributes from firestore")
			}
		} else {
			doc = resp
		}

		newLock := lock{
			AlertID:   alert.ID,
			ExpiresAt: timeout.UTC(),
		}

		if doc == nil {
			if err := tx.Create(x.client.Collection(x.attrCollection).Doc(key), newLock); err != nil {
				if status.Code(err) == codes.AlreadyExists {
					return goerr.Wrap(errLockFailed, "lock is already acquired")
				}
				return goerr.Wrap(err, "failed to create lock")
			}
		} else {
			var current lock
			if err := resp.DataTo(&current); err != nil {
				return goerr.Wrap(err, "failed to unmarshal lock")
			}

			if current.ExpiresAt.After(now) {
				return goerr.Wrap(errLockFailed, "lock is already acquired")
			}

			if err := tx.Set(doc.Ref, newLock); err != nil {
				return goerr.Wrap(err, "failed to update lock")
			}
		}

		return nil
	})

	if err != nil {
		return goerr.Wrap(err, "failed firestore transaction")
	}

	return nil
}

// Unlock implements interfaces.Database.
func (x *Client) Unlock(ctx *model.Context, ns types.Namespace) error {
	key := lockKeyPrefix + hashNamespace(ns)

	if _, err := x.client.Collection(x.attrCollection).Doc(key).Delete(ctx); err != nil {
		return goerr.Wrap(err, "failed to delete lock")
	}
	return nil
}

func New(ctx *model.Context, projectID string, collectionPrefix string) (*Client, error) {
	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return nil, goerr.Wrap(err, "Failed to initialize firebase app")
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, goerr.Wrap(err, "Failed to initialize firestore client")
	}

	return &Client{
		client:             client,
		projectID:          projectID,
		attrCollection:     collectionPrefix + ".attr",
		workflowCollection: collectionPrefix + ".workflow",
	}, nil
}

/*
func (x *Client) Migrate(ctx *model.Context) error {
	adminClient, err := apiv1.NewFirestoreAdminClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer adminClient.Close()

	op, err := adminClient.CreateIndex(ctx, &adminpb.CreateIndexRequest{
		Parent: "projects/" + x.projectID + "/databases/(default)/collectionGroups/" + x.collection,
		Index: &adminpb.Index{
			Fields: []*adminpb.Index_IndexField{
				{
					FieldPath: "RecordType",
					ValueMode: &adminpb.Index_IndexField_Order_{
						Order: adminpb.Index_IndexField_ASCENDING,
					},
				},
				{
					FieldPath: "CreatedAt",
					ValueMode: &adminpb.Index_IndexField_Order_{
						Order: adminpb.Index_IndexField_DESCENDING,
					},
				},
			},
			QueryScope: adminpb.Index_COLLECTION,
		},
	})
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			ctx.Logger().Info("Index already created")
			return nil
		}
		return goerr.Wrap(err, "failed to create index")
	}

	idx, err := op.Wait(ctx)
	if err != nil {
		return goerr.Wrap(err, "failed to wait index creation")
	}
	ctx.Logger().Info("Created index", slog.Any("index", idx))

	return nil
}
*/

func (x *Client) Close() error {
	if err := x.client.Close(); err != nil {
		return goerr.Wrap(err, "failed to close firestore client")
	}
	return nil
}

var _ interfaces.Database = &Client{}
