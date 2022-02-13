package spanner

import (
	"fmt"
	"os"
	"testing"

	"github.com/m-mizutani/alertchain/pkg/domain/types"

	"cloud.google.com/go/spanner"
	"github.com/m-mizutani/goerr"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
)

const (
	tblAlert      = "Alerts"
	tblAttribute  = "Attributes"
	tblAnnotation = "Annotations"
	tblReferences = "References"
)

type Client struct {
	client *spanner.Client
}

func New(ctx *types.Context, database string, options ...option.ClientOption) (*Client, error) {
	client, err := spanner.NewClient(ctx, database, options...)
	if err != nil {
		return nil, goerr.Wrap(err).With("database", database)
	}

	return &Client{
		client: client,
	}, nil
}

func NewTestDB(t *testing.T) *Client {
	dbName := fmt.Sprintf("projects/%s/instances/%s/databases/%s",
		os.Getenv("SPANNER_PROJECT_ID"),
		os.Getenv("SPANNER_INSTANCE_ID"),
		os.Getenv("SPANNER_DATABASE_ID"),
	)
	client, err := New(types.NewContext(), dbName)
	require.NoError(t, err)
	return client
}
