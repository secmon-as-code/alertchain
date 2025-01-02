package config

import (
	"context"

	"github.com/m-mizutani/goerr"
	"github.com/secmon-lab/alertchain/pkg/domain/interfaces"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
	"github.com/secmon-lab/alertchain/pkg/infra/firestore"
	"github.com/secmon-lab/alertchain/pkg/infra/memory"
	"github.com/secmon-lab/alertchain/pkg/utils"
	"github.com/urfave/cli/v3"
)

type Database struct {
	dbType              string
	firestoreProjectID  string
	firestoreDatabaseID string
}

func (x *Database) Flags() []cli.Flag {
	category := "Database"

	return []cli.Flag{
		&cli.StringFlag{
			Name:        "db-type",
			Usage:       "Database type (memory, firestore)",
			Category:    category,
			Aliases:     []string{"t"},
			Sources:     cli.EnvVars("ALERTCHAIN_DB_TYPE"),
			Value:       "memory",
			Destination: &x.dbType,
		},
		&cli.StringFlag{
			Name:        "firestore-project-id",
			Usage:       "Project ID of Firestore",
			Category:    category,
			Sources:     cli.EnvVars("ALERTCHAIN_FIRESTORE_PROJECT_ID"),
			Destination: &x.firestoreProjectID,
		},
		&cli.StringFlag{
			Name:        "firestore-database-id",
			Usage:       "Prefix of Firestore database ID",
			Category:    category,
			Sources:     cli.EnvVars("ALERTCHAIN_FIRESTORE_DATABASE_ID"),
			Destination: &x.firestoreDatabaseID,
		},
	}
}

func (x *Database) New(ctx context.Context) (interfaces.Database, func(), error) {
	nopCloser := func() {}

	switch x.dbType {
	case "memory":
		return memory.New(), nopCloser, nil

	case "firestore":
		if x.firestoreProjectID == "" {
			return nil, nopCloser, goerr.Wrap(types.ErrInvalidOption, "firestore-project-id is required for firestore")
		}
		if x.firestoreDatabaseID == "" {
			return nil, nopCloser, goerr.Wrap(types.ErrInvalidOption, "firestore-collection-prefix is required for firestore")
		}

		client, err := firestore.New(ctx, x.firestoreProjectID, x.firestoreDatabaseID)
		if err != nil {
			return nil, nopCloser, goerr.Wrap(err, "failed to initialize firestore client")
		}
		return client, func() { utils.SafeClose(ctx, client) }, nil

	default:
		return nil, nopCloser, goerr.Wrap(types.ErrInvalidOption, "invalid db-type").With("db-type", x.dbType)
	}
}
