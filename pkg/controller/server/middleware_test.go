package server_test

import (
	_ "embed"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/m-mizutani/alertchain/pkg/controller/server"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra/policy"
	"github.com/m-mizutani/gt"
)

//go:embed testdata/authz.rego
var authzRego string

func TestAuthorize(t *testing.T) {
	authz := gt.R1(policy.New(
		policy.WithPolicyData("authz.rego", authzRego),
		policy.WithPackage("authz"),
	)).NoError(t)
	srv := server.New(func(ctx *model.Context, schema types.Schema, data any) ([]*model.Alert, error) {
		return nil, nil
	}, server.WithAuthzPolicy(authz))

	testCases := map[string]struct {
		NewReq func() *http.Request
		Expect int
	}{
		"valid": {
			NewReq: func() *http.Request {
				return httptest.NewRequest("GET", "/health", nil)
			},
			Expect: http.StatusOK,
		},
		"not found": {
			NewReq: func() *http.Request {
				return httptest.NewRequest("GET", "/not-found", nil)
			},
			Expect: http.StatusNotFound,
		},
		"unauthorized by path": {
			NewReq: func() *http.Request {
				return httptest.NewRequest("GET", "/admin", nil)
			},
			Expect: http.StatusForbidden,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			req := tc.NewReq()
			w := httptest.NewRecorder()
			srv.ServeHTTP(w, req)
			gt.N(t, w.Result().StatusCode).Equal(tc.Expect)
		})
	}
}
