package model_test

import (
	"testing"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/gt"
)

func TestInvalidPath(t *testing.T) {
	c := model.Commit{
		Attribute: model.Attribute{
			Key: "hoge",
		},
		Path: "$.invalid",
	}

	data := map[string]interface{}{
		"hoge": "fuga",
	}

	v, err := c.ToAttr(data)
	gt.NoError(t, err)
	gt.EQ(t, v, nil)
}

func TestEmptyPathWithValue(t *testing.T) {
	c := model.Commit{
		Attribute: model.Attribute{
			Key:   "hoge",
			Value: "fuga",
		},
		Path: "",
	}

	data := map[string]interface{}{
		"hoge": "fuga",
	}

	v, err := c.ToAttr(data)
	gt.NoError(t, err)
	gt.EQ(t, v.Key, "hoge")
	gt.EQ(t, v.Value, "fuga")
}

func TestEmptyPathWithoutValue(t *testing.T) {
	c := model.Commit{
		Attribute: model.Attribute{
			Key: "hoge",
		},
		Path: "",
	}

	data := map[string]interface{}{
		"hoge": "fuga",
	}

	v, err := c.ToAttr(data)
	gt.Error(t, err)
	gt.EQ(t, v, nil)
}

func TestValidPath(t *testing.T) {
	c := model.Commit{
		Attribute: model.Attribute{
			Key: "hoge",
		},
		Path: "$.hoge",
	}

	data := map[string]interface{}{
		"hoge": "fuga",
	}

	v, err := c.ToAttr(data)
	gt.NoError(t, err)
	gt.EQ(t, v.Key, "hoge")
	gt.EQ(t, v.Value, "fuga")
}

func TestGetArrayObject(t *testing.T) {
	c := model.Commit{
		Attribute: model.Attribute{
			Key: "hoge",
		},
		Path: "$.array[0]",
	}

	data := map[string]interface{}{
		"array": []interface{}{
			"fuga",
		},
	}

	v, err := c.ToAttr(data)
	gt.NoError(t, err)
	gt.EQ(t, v.Key, "hoge")
	gt.EQ(t, v.Value, "fuga")
}
