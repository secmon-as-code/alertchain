package memory

import (
	"sort"
	"sync"
	"time"

	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

type lock struct {
	mutex     sync.Mutex
	expiresAt time.Time
}

type Client struct {
	attrs     map[types.Namespace]map[types.AttrID]*model.Attribute
	locks     map[types.Namespace]*lock
	workflows map[types.WorkflowID]model.WorkflowRecord
	alerts    map[types.AlertID]*model.Alert

	attrMutex     sync.RWMutex
	lockMutex     sync.Mutex
	workflowMutex sync.RWMutex
	alertMutex    sync.RWMutex
}

func New() *Client {
	return &Client{
		attrs:     map[types.Namespace]map[types.AttrID]*model.Attribute{},
		locks:     map[types.Namespace]*lock{},
		workflows: map[types.WorkflowID]model.WorkflowRecord{},
		alerts:    map[types.AlertID]*model.Alert{},
	}
}

// Close implements interfaces.Database.
func (x *Client) Close() error {
	return nil
}

// GetAttrs implements interfaces.Database.
func (x *Client) GetAttrs(ctx *model.Context, ns types.Namespace) (model.Attributes, error) {
	x.attrMutex.RLock()
	defer x.attrMutex.RUnlock()

	attrs, ok := x.attrs[ns]
	if !ok {
		return nil, nil
	}

	var ret model.Attributes
	for _, a := range attrs {
		ret = append(ret, *a)
	}

	return ret, nil
}

// PutAttrs implements interfaces.Database.
func (x *Client) PutAttrs(ctx *model.Context, ns types.Namespace, attrs model.Attributes) error {
	x.attrMutex.Lock()
	if _, ok := x.attrs[ns]; !ok {
		x.attrs[ns] = map[types.AttrID]*model.Attribute{}
	}
	x.attrMutex.Unlock()

	for i, src := range attrs {
		if dst, ok := x.attrs[ns][src.ID]; ok {
			dst.Value = src.Value
		} else {
			x.attrs[ns][src.ID] = &attrs[i]
		}
	}

	return nil
}

func (x *Client) PutWorkflow(ctx *model.Context, workflow model.WorkflowRecord) error {
	x.workflowMutex.Lock()
	defer x.workflowMutex.Unlock()

	x.workflows[workflow.ID] = workflow
	return nil
}

func (x *Client) GetWorkflows(ctx *model.Context, offset, limit int) ([]model.WorkflowRecord, error) {
	x.workflowMutex.RLock()
	defer x.workflowMutex.RUnlock()

	workflows := make([]model.WorkflowRecord, 0, len(x.workflows))
	for _, wf := range x.workflows {
		workflows = append(workflows, wf)
	}
	// sort
	sort.Slice(workflows, func(i, j int) bool {
		return workflows[i].CreatedAt.After(workflows[j].CreatedAt)
	})

	if offset >= len(workflows) {
		return nil, nil
	}

	end := offset + limit
	if end > len(workflows) {
		end = len(workflows)
	}

	return workflows[offset:end], nil
}

func (x *Client) GetWorkflow(ctx *model.Context, id types.WorkflowID) (*model.WorkflowRecord, error) {
	for _, wf := range x.workflows {
		if wf.ID == id {
			return &wf, nil
		}
	}

	return nil, nil
}

func (x *Client) PutAlert(ctx *model.Context, alert model.Alert) error {
	x.alertMutex.Lock()
	defer x.alertMutex.Unlock()

	x.alerts[alert.ID] = &alert
	return nil
}

func (x *Client) GetAlert(ctx *model.Context, id types.AlertID) (*model.Alert, error) {
	x.alertMutex.RLock()
	defer x.alertMutex.RUnlock()

	if alert, ok := x.alerts[id]; ok {
		return alert, nil
	}

	return nil, nil
}

// Lock implements interfaces.Database.
func (x *Client) Lock(ctx *model.Context, ns types.Namespace, timeout time.Time) error {
	locked := make(chan struct{}, 1)
	go func() {
		for i := 0; ; i++ {
			x.lockMutex.Lock()
			if l, ok := x.locks[ns]; !ok || l.expiresAt.Before(time.Now()) {
				x.locks[ns] = &lock{}
			}
			x.lockMutex.Unlock()

			if x.locks[ns].mutex.TryLock() {
				x.locks[ns].expiresAt = timeout
				break
			}

			time.Sleep(10 * time.Millisecond)
		}
		locked <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		// cancelled
	case <-locked:
		// locked
	}

	return nil
}

// Unlock implements interfaces.Database.
func (x *Client) Unlock(ctx *model.Context, ns types.Namespace) error {
	if _, ok := x.locks[ns]; !ok {
		return nil
	}

	x.locks[ns].mutex.Unlock()
	return nil
}

var _ interfaces.Database = (*Client)(nil)
