package actions

import (
	"fmt"

	"github.com/m-mizutani/alertchain/pkg/actions/github"
	"github.com/m-mizutani/alertchain/pkg/actions/otx"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/goerr"
)

var repository map[string]model.ActionFactory = map[string]model.ActionFactory{}

func Register(id string, factory model.ActionFactory) {
	if _, ok := repository[id]; ok {
		panic(fmt.Sprintf("action factory '%s' is already registered", id))
	}

	repository[id] = factory
}

func init() {
	Register(github.CreateIssueID, github.NewCreateIssue)
	Register(otx.InquiryID, otx.NewInquiry)
}

func New(id string, cfg model.ActionConfig) (model.Action, error) {
	factory, ok := repository[id]
	if !ok {
		return nil, goerr.Wrap(types.ErrActionNotFound).With("action", id).With("config", cfg)
	}

	newAction, err := factory(cfg)
	if err != nil {
		return nil, err
	}

	return newAction, nil
}
