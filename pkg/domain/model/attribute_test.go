package model_test

import (
	"testing"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/gt"
)

func TestTidyAttributes(t *testing.T) {
	attrs := model.Attributes{
		{ID: "1", Name: "pattr1", Value: "value1", Type: "type1"},
		{ID: "2", Name: "pattr2", Value: "value2", Type: "type2"},
		{ID: "1", Name: "pattr1_updated", Value: "value1_updated", Type: "type1_updated"},
		{ID: "3", Name: "pattr3", Value: "value3", Type: "type3"},
	}

	expected := model.Attributes{
		{ID: "1", Name: "pattr1_updated", Value: "value1_updated", Type: "type1_updated"},
		{ID: "2", Name: "pattr2", Value: "value2", Type: "type2"},
		{ID: "3", Name: "pattr3", Value: "value3", Type: "type3"},
	}

	result := model.TidyAttributes(attrs)
	gt.A(t, result).Equal(expected)
}
