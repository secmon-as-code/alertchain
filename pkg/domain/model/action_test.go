package model_test

import (
	"testing"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/gt"
)

func TestActionArgsParser(t *testing.T) {
	args := model.ActionArgs{
		"foo": "bar",
		"bar": int64(123),
		"baz": float64(3.14),
		"raw": []byte("raw"),
		"ss":  []string{"a", "b"},
	}

	t.Run("parse all", func(t *testing.T) {
		var foo string
		var bar int64
		var baz float64
		var raw []byte
		var ss []string

		gt.NoError(t, args.Parse(
			model.ArgDef("foo", &foo),
			model.ArgDef("bar", &bar),
			model.ArgDef("baz", &baz),
			model.ArgDef("raw", &raw),
			model.ArgDef("ss", &ss),
		))
		gt.Equal(t, "bar", foo)
		gt.Equal(t, int64(123), bar)
		gt.Equal(t, float64(3.14), baz)
		gt.Equal(t, []byte("raw"), raw)
		gt.A(t, ss).Length(2).Have("a").Have("b")
	})

	t.Run("parse partial", func(t *testing.T) {
		var foo string
		var baz float64
		gt.NoError(t, args.Parse(
			model.ArgDef("foo", &foo),
			model.ArgDef("baz", &baz),
		))
		gt.Equal(t, "bar", foo)
		gt.Equal(t, float64(3.14), baz)
	})

	t.Run("parse error", func(t *testing.T) {
		var foo string

		gt.Error(t, args.Parse(
			model.ArgDef("xxx", &foo),
		))
	})

	t.Run("type error", func(t *testing.T) {
		var foo int64

		gt.Error(t, args.Parse(
			model.ArgDef("foo", &foo),
		))
	})
}
