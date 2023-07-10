package firestore

import (
	"context"
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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Client struct {
	client     *firestore.Client
	collection string
}

const (
	attrKeyPrefix = "attr:"
	lockKeyPrefix = "lock:"
)

// GetAttrs implements interfaces.Database.
func (x *Client) GetAttrs(ctx *model.Context, ns types.Namespace) (model.Attributes, error) {
	key := attrKeyPrefix + string(ns)
	docs, err := x.client.Collection(x.collection).Doc(key).Collection("attributes").Documents(ctx).GetAll()
	if err != nil {
		return nil, goerr.Wrap(err, "failed to get attributes from firestore")
	}

	attrs := make([]model.Attribute, len(docs))

	for i, doc := range docs {
		if doc.Exists() {
			if err := doc.DataTo(&attrs[i]); err != nil {
				return nil, goerr.Wrap(err, "failed to unmarshal attribute from firestore")
			}
		}
	}

	return attrs, nil
}

// PutAttrs implements interfaces.Database.
func (x *Client) PutAttrs(ctx *model.Context, ns types.Namespace, attrs model.Attributes) error {
	err := x.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		key := attrKeyPrefix + string(ns)
		collection := x.client.Collection(x.collection).Doc(key).Collection("attributes")

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

		now := time.Now().UTC().Unix()

		for _, attr := range attrs {
			if attr.ExpiresAt == 0 {
				attr.ExpiresAt = now + types.DefaultAttributeTTL
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

type lock struct {
	AlertID   types.AlertID `firestore:"alert_id"`
	ExpiresAt int64         `firestore:"expires_at"`
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
	key := lockKeyPrefix + string(ns)
	now := time.Now().UTC().Unix()
	alert := ctx.Alert()

	err := x.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		var doc *firestore.DocumentSnapshot
		resp, err := tx.Get(x.client.Collection(x.collection).Doc(key))
		if err != nil {
			if status.Code(err) != codes.NotFound {
				return goerr.Wrap(err, "failed to get attributes from firestore")
			}
		} else {
			doc = resp
		}

		newLock := lock{
			AlertID:   alert.ID,
			ExpiresAt: timeout.UTC().Unix(),
		}

		if doc == nil {
			if err := tx.Create(x.client.Collection(x.collection).Doc(key), newLock); err != nil {
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

			if current.ExpiresAt > now {
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
	key := lockKeyPrefix + string(ns)

	if _, err := x.client.Collection(x.collection).Doc(key).Delete(ctx); err != nil {
		return goerr.Wrap(err, "failed to delete lock")
	}
	return nil
}

func New(ctx *model.Context, projectID string, collection string) (*Client, error) {
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
		client:     client,
		collection: collection,
	}, nil
}

func (x *Client) Close() error {
	if err := x.client.Close(); err != nil {
		return goerr.Wrap(err, "failed to close firestore client")
	}
	return nil
}

var _ interfaces.Database = &Client{}
