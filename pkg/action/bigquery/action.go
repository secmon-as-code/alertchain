package bigquery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/secmon-lab/alertchain/pkg/ctxutil"
	"github.com/secmon-lab/alertchain/pkg/domain/model"
	"github.com/secmon-lab/alertchain/pkg/domain/types"

	"cloud.google.com/go/bigquery"
	"github.com/m-mizutani/goerr"
	"google.golang.org/api/googleapi"
)

type DataRecord struct {
	ID        string    `bigquery:"id"`
	AlertID   string    `bigquery:"alert_id"`
	CreatedAt time.Time `bigquery:"created_at"`
	Tags      []string  `bigquery:"tags"`
	Data      string    `bigquery:"data"`
}

func InsertData(ctx context.Context, alert model.Alert, args model.ActionArgs) (any, error) {
	table, err := setupTable(ctx, args)
	if err != nil {
		return nil, err
	}

	data, ok := args["data"]
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidArgument, "data is required")
	}
	if data == nil {
		return nil, goerr.Wrap(types.ErrActionInvalidArgument, "data must not be nil")
	}

	var tags []string
	if v, ok := args["tags"].([]string); ok {
		tags = v
	}

	raw, err := json.Marshal(data)
	if err != nil {
		return nil, goerr.Wrap(err, "Fail to marshal data")
	}

	row := DataRecord{
		ID:        uuid.NewString(),
		AlertID:   alert.ID.String(),
		CreatedAt: ctxutil.Now(ctx),
		Tags:      tags,
		Data:      string(raw),
	}

	schema, err := bigquery.InferSchema(row)
	if err != nil {
		return nil, goerr.Wrap(err, "Fail to infer schema")
	}
	for i := range schema {
		if schema[i].Name == "data" {
			schema[i].Type = bigquery.JSONFieldType
		}
		if schema[i].Name == "tags" {
			schema[i].Required = false
		}
	}

	return nil, insert(ctx, table, schema, row)
}

type AlertRecord struct {
	ID          types.AlertID    `bigquery:"id"`
	Schema      types.Schema     `bigquery:"schema"`
	CreatedAt   time.Time        `bigquery:"created_at"`
	Title       string           `bigquery:"title"`
	Description string           `bigquery:"description"`
	Source      string           `bigquery:"source"`
	Namespace   types.Namespace  `bigquery:"namespace"`
	Attrs       []AttrRecord     `bigquery:"attrs"`
	Refs        model.References `bigquery:"refs"`
	Data        string           `bigquery:"data"`
}

type AttrRecord struct {
	ID     string `bigquery:"id"`
	Key    string `bigquery:"key"`
	Value  string `bigquery:"value"`
	Type   string `bigquery:"type"`
	TTL    int    `bigquery:"ttl"`
	Global bool   `bigquery:"global"`
}

func InsertAlert(ctx context.Context, alert model.Alert, args model.ActionArgs) (any, error) {
	table, err := setupTable(ctx, args)
	if err != nil {
		return nil, err
	}

	raw, err := json.Marshal(alert.Data)
	if err != nil {
		return nil, goerr.Wrap(err, "Fail to marshal data")
	}

	row := AlertRecord{
		ID:          alert.ID,
		Schema:      alert.Schema,
		CreatedAt:   alert.CreatedAt,
		Title:       alert.Title,
		Description: alert.Description,
		Source:      alert.Source,
		Namespace:   alert.Namespace,
		Refs:        alert.Refs,
		Data:        string(raw),
	}

	for _, attr := range alert.Attrs {
		row.Attrs = append(row.Attrs, AttrRecord{
			ID:     attr.ID.String(),
			Key:    string(attr.Key),
			Value:  fmt.Sprintf("%v", attr.Value),
			Type:   string(attr.Type),
			TTL:    attr.TTL,
			Global: attr.Global,
		})
	}

	schema, err := bigquery.InferSchema(row)
	if err != nil {
		return nil, goerr.Wrap(err, "Fail to infer schema")
	}
	for i := range schema {
		if schema[i].Name == "data" {
			schema[i].Type = bigquery.JSONFieldType
		}
	}

	return nil, insert(ctx, table, schema, row)
}

func setupTable(ctx context.Context, args model.ActionArgs) (*bigquery.Table, error) {
	projectID, ok := args["project_id"].(string)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidArgument, "project_id is required")
	}
	datasetID, ok := args["dataset_id"].(string)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidArgument, "dataset_id is required")
	}
	tableID, ok := args["table_id"].(string)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidArgument, "table_id is required")
	}

	c, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, goerr.Wrap(err, "Fail to create BigQuery client")
	}

	dataSet := c.Dataset(datasetID)

	return dataSet.Table(tableID), nil
}

func insert(ctx context.Context, table *bigquery.Table, schema bigquery.Schema, data any) error {
	if _, err := table.Metadata(ctx); err != nil {
		if gerr, ok := err.(*googleapi.Error); !ok || gerr.Code != 404 {
			return goerr.Wrap(err, "failed to get metadata of table")
		}

		// Table not found
		meta := &bigquery.TableMetadata{
			Schema: schema,
			TimePartitioning: &bigquery.TimePartitioning{
				Field: "created_at",
			},
		}
		if err := table.Create(ctx, meta); err != nil {
			if gerr, ok := err.(*googleapi.Error); !ok || gerr.Code != 409 {
				return goerr.Wrap(err, "failed to create table of data")
			}
			// ignore 409 error
		}
	}

	if err := insertWithRetry(ctx, table, data); err != nil {
		return goerr.Wrap(err, "Fail to insert data").With("table", table)
	}

	return nil
}

func insertWithRetry(ctx context.Context, table *bigquery.Table, data any) error {
	// Define the maximum number of retries and the initial delay.
	const maxRetries = 3
	initialDelay := time.Millisecond * 100

	// Attempt to insert data with exponential backoff retry logic.
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// Wait for the delay period before retrying.
			time.Sleep(initialDelay)
			// Increase the delay for the next retry.
			initialDelay *= 2
		}

		inserter := table.Inserter()
		err := inserter.Put(ctx, data)
		if err == nil {
			// Data inserted successfully, no need to retry.
			return nil
		}

		if e, ok := err.(*googleapi.Error); ok && e.Code != 404 {
			return err
		}

		ctxutil.Logger(ctx).Warn("Table not found, retrying",
			"table", table.FullyQualifiedName(),
			"err", err,
		)
	}

	// Data insertion failed after all retries.
	return errors.New("insert failed: exceeded retry limit")
}
