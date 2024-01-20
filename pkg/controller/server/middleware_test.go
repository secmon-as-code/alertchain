package server_test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/m-mizutani/alertchain/pkg/controller/server"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra/policy"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/gt"
)

//go:embed testdata/authz.rego
var authzRego string

func newServer(t *testing.T, policyData string, options ...server.Option) *server.Server {

	authz := gt.R1(policy.New(
		policy.WithPolicyData("authz.rego", policyData),
		policy.WithPackage("authz"),
	)).NoError(t)

	options = append(options, server.WithAuthzPolicy(authz))

	srv := server.New(func(ctx *model.Context, schema types.Schema, data any) ([]*model.Alert, error) {
		return nil, nil
	}, options...)

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

func hmacSign(secretKey, timestamp string, data []byte) string {
	msg := fmt.Sprintf("%s%s", data, timestamp)
	fmt.Println("msg=", msg)

	computed := hmac.New(sha256.New, []byte(secretKey))
	computed.Write([]byte(msg))
	computedBytes := computed.Sum(nil)

	return base64.StdEncoding.EncodeToString(computedBytes)
}

func TestVerifySignedBody(t *testing.T) {
	data := []byte(`{"foo":"bar"}`)
	ts := time.Now().Format("2006-01-02T15:04:05Z")

	const hmacSecret = "Caprice_of_the_Leaves"
	sign := hmacSign(hmacSecret, ts, data)

	srv := newServer(t, authzRego,
		server.WithEnv(func() types.EnvVars {
			return types.EnvVars{
				"CLOUDSTRIKE_HAWK_KEY": hmacSecret,
			}
		}),
	)

	testCases := map[string]struct {
		path   string
		build  func(r *http.Request)
		data   []byte
		expect int
	}{
		"valid": {
			path: "/alert/raw/cloudstrike_hawk",
			build: func(r *http.Request) {
				r.Header.Set("X-Cs-Delivery-Timestamp", ts)
				r.Header.Set("X-Cs-Primary-Signature", sign)
			},
			data:   data,
			expect: http.StatusOK,
		},
		"invalid timestamp": {
			path: "/alert/raw/cloudstrike_hawk",
			build: func(r *http.Request) {
				r.Header.Set("X-Cs-Delivery-Timestamp", "invalid")
				r.Header.Set("X-Cs-Primary-Signature", sign)
			},
			data:   data,
			expect: http.StatusForbidden,
		},
		"invalid signature": {
			path: "/alert/raw/cloudstrike_hawk",
			build: func(r *http.Request) {
				r.Header.Set("X-Cs-Delivery-Timestamp", ts)
				r.Header.Set("X-Cs-Primary-Signature", "invalid")
			},
			data:   data,
			expect: http.StatusForbidden,
		},
		"invalid path": {
			path: "/admin",
			build: func(r *http.Request) {
				r.Header.Set("X-Cs-Delivery-Timestamp", ts)
				r.Header.Set("X-Cs-Primary-Signature", sign)
			},
			data:   data,
			expect: http.StatusForbidden,
		},
		"invalid data": {
			path: "/alert/raw/cloudstrike_hawk",
			build: func(r *http.Request) {
				r.Header.Set("X-Cs-Delivery-Timestamp", ts)
				r.Header.Set("X-Cs-Primary-Signature", sign)
			},
			data:   []byte("invalid"),
			expect: http.StatusForbidden,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest("POST", tc.path, bytes.NewReader(tc.data))
			tc.build(req)
			w := httptest.NewRecorder()
			srv.ServeHTTP(w, req)
			gt.N(t, w.Result().StatusCode).Equal(tc.expect)
			t.Log(name, w.Result().StatusCode)
		})
	}
}
