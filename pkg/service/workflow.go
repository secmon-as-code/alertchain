package service

import (
	"encoding/json"
	"fmt"

	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/goerr"
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

func (x *WorkflowService) Get(ctx *model.Context, offset, limit *int) ([]model.WorkflowRecord, error) {
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

func (x *WorkflowService) Lookup(ctx *model.Context, id types.WorkflowID) (*model.WorkflowRecord, error) {
	return nil, nil
}

func (x *WorkflowService) Create(ctx *model.Context, alert model.Alert) (*Workflow, error) {
	rawData, err := json.Marshal(alert.Data)
	if err != nil {
		return nil, goerr.Wrap(err, "Fail to marshal alert data")
	}

	initAttrs := make([]*model.AttributeRecord, len(alert.Attrs))
	for i, attr := range alert.Attrs {
		var typ *string
		if attr.Type != "" {
			typ = (*string)(&attr.Type)
		}

		initAttrs[i] = &model.AttributeRecord{
			ID:     string(attr.ID),
			Key:    string(attr.Key),
			Value:  fmt.Sprintf("%+v", attr.Value),
			Type:   typ,
			Global: attr.Global,
			TTL:    int(attr.TTL),
		}
	}

	var namespace *string
	if alert.Namespace != "" {
		namespace = (*string)(&alert.Namespace)
	}

	workflow := model.WorkflowRecord{
		ID:        types.NewWorkflowID().String(),
		CreatedAt: ctx.Now(),
		Alert: &model.AlertRecord{
			ID:          string(alert.ID),
			Schema:      string(alert.Schema),
			Data:        string(rawData),
			CreatedAt:   alert.CreatedAt,
			Title:       alert.Title,
			Source:      alert.Source,
			InitAttrs:   initAttrs,
			Description: alert.Description,
			Namespace:   namespace,
		},
	}

	if err := x.db.PutWorkflow(ctx, workflow); err != nil {
		return nil, err
	}

	return &Workflow{db: x.db, wf: &workflow}, nil
}

func (x *Workflow) AddAction(ctx *model.Context, action *model.ActionRecord) error {

	return nil
}
