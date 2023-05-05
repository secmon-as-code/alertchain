package model_test

import (
	"testing"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/gt"
)

func TestTidyParameters(t *testing.T) {
	params := model.Parameters{
		{ID: "1", Name: "param1", Value: "value1", Type: "type1"},
		{ID: "2", Name: "param2", Value: "value2", Type: "type2"},
		{ID: "1", Name: "param1_updated", Value: "value1_updated", Type: "type1_updated"},
		{ID: "3", Name: "param3", Value: "value3", Type: "type3"},
	}

	expected := model.Parameters{
		{ID: "1", Name: "param1_updated", Value: "value1_updated", Type: "type1_updated"},
		{ID: "2", Name: "param2", Value: "value2", Type: "type2"},
		{ID: "3", Name: "param3", Value: "value3", Type: "type3"},
	}

	result := model.TidyParameters(params)
	gt.A(t, result).Equal(expected)
}
