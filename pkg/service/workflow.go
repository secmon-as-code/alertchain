package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/m-mizutani/goerr/v2"
	"github.com/secmon-lab/alertchain/pkg/ctxutil"
	"github.com/secmon-lab/alertchain/pkg/domain/interfaces"
	"github.com/secmon-lab/alertchain/pkg/domain/model"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
)

type WorkflowService struct {
	db interfaces.Database
}

type Workflow struct {
	db interfaces.Database
	wf *model.WorkflowRecord
}

func NewWorkflowService(db interfaces.Database) *WorkflowService {
	return &WorkflowService{db: db}
}

func (x *WorkflowService) Get(ctx context.Context, offset, limit *int) ([]model.WorkflowRecord, error) {
	if offset == nil {
		offset = new(int)
		*offset = 0
	}
	if limit == nil {
		limit = new(int)
		*limit = 20
	}

	return x.db.GetWorkflows(ctx, *offset, *limit)
}

func (x *WorkflowService) Lookup(ctx context.Context, id types.WorkflowID) (*model.WorkflowRecord, error) {
	return nil, nil
}

func attrsToRecord(attrs model.Attributes) []*model.AttributeRecord {
	records := make([]*model.AttributeRecord, len(attrs))
	for i, attr := range attrs {
		var typ *string
		if attr.Type != "" {
			typ = (*string)(&attrs[i].Type)
		}

		records[i] = &model.AttributeRecord{
			ID:      string(attr.ID),
			Key:     string(attr.Key),
			Value:   fmt.Sprintf("%+v", attr.Value),
			Type:    typ,
			Persist: attr.Persist,
			TTL:     int(attr.TTL),
		}
	}

	return records
}

func (x *WorkflowService) Create(ctx context.Context, alert model.Alert) (*Workflow, error) {
	rawData, err := json.Marshal(alert.Data)
	if err != nil {
		return nil, goerr.Wrap(err, "Fail to marshal alert data", goerr.T(types.ErrTagBadRequest))
		// types.AsBadRequestErr(goerr.Wrap(err, "Fail to marshal alert data"))
	}

	var namespace *string
	if alert.Namespace != "" {
		namespace = (*string)(&alert.Namespace)
	}

	if err := x.db.PutAlert(ctx, alert); err != nil {
		return nil, err
	}

	workflow := model.WorkflowRecord{
		ID:        types.NewWorkflowID(),
		CreatedAt: ctxutil.Now(ctx),
		Alert: &model.AlertRecord{
			ID:          alert.ID,
			Schema:      string(alert.Schema),
			Data:        string(rawData),
			CreatedAt:   alert.CreatedAt,
			Title:       alert.Title,
			Source:      alert.Source,
			InitAttrs:   attrsToRecord(alert.Attrs),
			Description: alert.Description,
			Namespace:   namespace,
		},
	}

	if err := x.db.PutWorkflow(ctx, workflow); err != nil {
		return nil, err
	}

	return &Workflow{db: x.db, wf: &workflow}, nil
}

func (x *Workflow) UpdateLastAttrs(ctx context.Context, attrs model.Attributes) error {
	x.wf.Alert.LastAttrs = attrsToRecord(attrs)
	if err := x.db.PutWorkflow(ctx, *x.wf); err != nil {
		return err
	}
	return nil
}

func (x *Workflow) AddAction(ctx context.Context, action *model.Action) error {
	return nil
}
