package model_test

import (
	_ "embed"
	"testing"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/gt"
)

//go:embed testdata/config1.jsonnet
var config1 string

func TestParseConfig(t *testing.T) {
	var cfg model.Config
	var vars []model.EnvVar = []model.EnvVar{
		{
			Key:   "COLOR",
			Value: "orange",
		},
	}
	gt.NoError(t, model.ParseConfig("testdata/config1.jsonnet", config1, vars, &cfg))
	gt.Array(t, cfg.Actions).Length(1).Elem(0, func(t testing.TB, v model.ActionConfig) {
		gt.V(t, v.ID).Equal("test-scc")
		gt.V(t, v.Name).Equal("scc")
		gt.Map(t, v.Config).EqualAt("data", "orange")
	})
}

func TestParseConfigWithNoFileName(t *testing.T) {
	var cfg model.Config
	var vars []model.EnvVar = []model.EnvVar{
		{
			Key:   "COLOR",
			Value: "orange",
		},
	}
	gt.NoError(t, model.ParseConfig("", config1, vars, &cfg))
}
