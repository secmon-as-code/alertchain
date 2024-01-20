package server_test

import (
	"bytes"
	_ "embed"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"

	"github.com/m-mizutani/alertchain/pkg/controller/server"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra/policy"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/gt"
)

//go:embed testdata/authz.rego
var authzRego string

func newServer(t *testing.T, policyData string) *server.Server {
	authz := gt.R1(policy.New(
		policy.WithPolicyData("authz.rego", policyData),
		policy.WithPackage("authz"),
	)).NoError(t)
	srv := server.New(func(ctx *model.Context, schema types.Schema, data any) ([]*model.Alert, error) {
		return nil, nil
	}, server.WithAuthzPolicy(authz))

	return srv
}

func TestAuthorize(t *testing.T) {
	srv := newServer(t, authzRego)
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
			Expect: http.StatusForbidden,
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

func TestVerifyGoogleIDToken(t *testing.T) {
	var email string
	if err := utils.LoadEnv(
		utils.EnvDef("TEST_GOOGLE_CLOUD_ACCOUNT_EMAIL", &email),
	); err != nil {
		t.Skip("Skip TestVerifyGoogleIDToken because TEST_GOOGLE_CLOUD_ACCOUNT_EMAIL is not set")
	}

	srv := newServer(t, strings.ReplaceAll(authzRego, "__GOOGLE_CLOUD_ACCOUNT_EMAIL__", email))

	cmd := exec.Command("gcloud", "auth", "print-identity-token")
	validToken := strings.TrimSpace(string(gt.R1(cmd.Output()).NoError(t)))

	testCases := map[string]struct {
		path   string
		build  func(r *http.Request)
		expect int
	}{
		"valid": {
			path: "/alert/raw/test",
			build: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer "+validToken)
			},
			expect: http.StatusOK,
		},
		"invalid token": {
			path: "/alert/raw/test",
			build: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer invalid-token")
			},
			expect: http.StatusForbidden,
		},
		"no token": {
			path:   "/alert/raw/test",
			build:  func(r *http.Request) {},
			expect: http.StatusForbidden,
		},
		"invalid path": {
			path: "/admin",
			build: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer "+validToken)
			},
			expect: http.StatusForbidden,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest("POST", tc.path, bytes.NewReader([]byte("{}")))
			tc.build(req)
			w := httptest.NewRecorder()
			srv.ServeHTTP(w, req)
			gt.N(t, w.Result().StatusCode).Equal(tc.expect)
			t.Log(name, w.Result().StatusCode)
		})
	}
}
