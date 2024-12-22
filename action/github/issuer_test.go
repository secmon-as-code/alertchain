package github_test

import (
	"bytes"
	"context"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/m-mizutani/gt"
	"github.com/secmon-lab/alertchain/action/github"
	"github.com/secmon-lab/alertchain/pkg/ctxutil"
	"github.com/secmon-lab/alertchain/pkg/domain/model"
	"github.com/secmon-lab/alertchain/pkg/domain/types"

	gh "github.com/google/go-github/github"
)

func TestIssueTemplate(t *testing.T) {
	var buf bytes.Buffer
	gt.NoError(t, github.ExecuteTemplate(&buf, model.Alert{
		AlertMetaData: model.AlertMetaData{
			Title:       "blue",
			Description: "orange",
			Attrs: []model.Attribute{
				{
					Key:   "magic",
					Value: "five",
				},
				{
					Key:   "star",
					Value: "light",
				},
				{
					Key:   "int",
					Value: 123,
				},
				{
					Key:   "struct",
					Value: struct{ Foo string }{Foo: "bar"},
				},
				{
					Key:   "md-title",
					Value: "*md-test*",
					Type:  types.MarkDown,
				},
			},
		},
		Schema: "fire",
		Raw:    `{"foo": "bar"}`,
	}))

	s := buf.String()
	gt.B(t, strings.Contains(s, "orange")).True()
	gt.B(t, strings.Contains(s, "| magic | `five` |")).True()
	gt.B(t, strings.Contains(s, "| star | `light` |")).True()
	gt.B(t, strings.Contains(s, "| int | `123` |")).True()
	gt.B(t, strings.Contains(s, "| struct | `{bar}` |")).True()
	gt.B(t, strings.Contains(s, `{"foo": "bar"}`)).True()

	gt.B(t, strings.Contains(s, "| md-title | `*md-test*` |")).False()
	gt.B(t, strings.Contains(s, "## Comments")).True()
	gt.B(t, strings.Contains(s, "### md-title")).True()
	gt.B(t, strings.Contains(s, "*md-test*")).True()

	// os.WriteFile("test.md", []byte(s), 0644)
}

func TestIssueTemplateNoMarkdown(t *testing.T) {
	var buf bytes.Buffer
	gt.NoError(t, github.ExecuteTemplate(&buf, model.Alert{
		AlertMetaData: model.AlertMetaData{
			Title:       "blue",
			Description: "orange",
			Attrs: []model.Attribute{
				{
					Key:   "magic",
					Value: "five",
				},
				{
					Key:   "star",
					Value: "light",
				},
				{
					Key:   "int",
					Value: 123,
				},
				{
					Key:   "struct",
					Value: struct{ Foo string }{Foo: "bar"},
				},
			},
		},
		Schema: "fire",
		Raw:    `{"foo": "bar"}`,
	}))

	s := buf.String()
	gt.B(t, strings.Contains(s, "## Comments")).False()
}

func TestIssuer(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_GITHUB_ISSUER"); !ok {
		t.Skip("Skipping test because TEST_GITHUB_ISSUER is not set")
	}

	cfg := model.ActionArgs{
		"app_id":             float64(gt.R1(strconv.Atoi(os.Getenv("TEST_GITHUB_APP_ID"))).NoError(t)),
		"install_id":         float64(gt.R1(strconv.Atoi(os.Getenv("TEST_GITHUB_INSTALL_ID"))).NoError(t)),
		"secret_private_key": os.Getenv("TEST_GITHUB_PRIVATE_KEY"),
		"owner":              os.Getenv("TEST_GITHUB_OWNER"),
		"repo":               os.Getenv("TEST_GITHUB_REPO"),
	}

	requiredVars := []string{"app_id", "install_id", "secret_private_key", "owner", "repo"}
	for _, key := range requiredVars {
		gt.V(t, cfg[key]).NotEqual("")
	}

	ctx := context.Background()
	alert := model.Alert{
		AlertMetaData: model.AlertMetaData{
			Title:       "blue",
			Description: "orange",
			Attrs: []model.Attribute{
				{
					Key:   "magic",
					Value: "five",
				},
			},
		},
		CreatedAt: time.Now(),
		Raw:       `{"foo": "bar"}`,
	}

	args := model.ActionArgs{
		"assignee": "m-mizutani",
		"labels":   []string{"bug", "help wanted", "dummy"},
	}
	resp := gt.R1(github.CreateIssue(ctx, alert, args)).NoError(t)
	issue := gt.Cast[*gh.Issue](t, resp)
	gt.V(t, issue.Title).Must().NotNil().Equal(&alert.Title)
}

const dummyPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAwXBU34FShax6rMKB4jgRny05cW05N8JbAENnmKb5BUlYQHa/
kTOl8IUeAaDCoA4wXhTovpyGxA2MyoLqEHPvfuPauEm/6kAkVcfw0Owux05X5SiC
6F8mxfs9aTHafP2oLFUfNJJepEV/u8p6BHYE0B2CAWe9D6ND50XSqA1aH6EYeWDK
GQLl58GupaY5GpYZuD6BWMyn7rTMaP5HkH8LwXTZAqclyn6mLUwPsSHP80WYK+pd
cOAmnP8+WTqzHvMALjeR1HUvR3D23HwnPRdtSe+can5xX+hXxbFwHDay9iWfmpUG
OCsta4uWIqmYnH2/mW+gAiuDb3seS8Gkf3UGLQIDAQABAoIBACW4rhRXt6vxkoqV
85YVsPoFa6o+zmWdNPm8KzuNdAof32HSxlCebcGVc+CFZO6pVa1DDo/9HhqlOctT
9Cj5Mr7f2AsP9qjLkUpZDxDuvcCH+oPpfn2p8HmzIKqe2ih9nonmn4s079fA5cPN
HDY6fX3IA04a2Ldv8xHqf8XdtLFt9flwDmLn5PUfsEjmcWOmdtLdyRx8/zY8081e
ueODmypnrXFF9denq29WgRtWbmfEWdOh5e/eZRixn8u0gxMcAk0/g+v3AaLmsN8P
UCRaXcO+SZJMKBgFBRdFCGWsouJB9JdR+RWfQVvv5Z/Dfb5qtL5dyDOXhlTnyfSs
ys0qagECgYEA3qLtdnW7v4SYC9jh+/XVedH10yuNeb7FBuHjgs+1U2hqExqe+LU2
osidWeM1WDuvzK9fTzcrzGRfUFTWLUDuc3smz1kBMigC1YaB22uuFDe4/05G68fG
AB9ZkzWT9N1igyzwIt7fJp569naw7R5koaqtXfP2Sp00+kchvUpXfg0CgYEA3m1J
GrCIdGitp29yOZN9B44Nye9EjghGGw60tDPRW8NPrw041G6rtQz8lTyV+kb0OwVn
DERawMy4wj3/G96EUi7SjyHlA2oZrs1HyxBCcd4hg7Ll0410h1GHOTslp+yYqvIK
KIFFzBXfos2uhh2y0FmYYZjo+/9sEA6+14s5wKECgYEAmvyDIN8u91Ff44d1Mkjd
9rMyVXJRR7qFQJhKIItmKI1corX6ixrj0QileajRPv42EODZEbVPmTcanzqf6trz
5JKL3vaP/ZGa/3hmuBBLHCn6cEjW2Fa3QOiSHAfFW0YuyTCkbzIF2MWkxiS0YC2z
UlQV4nzuLN0pvz17gGHbbJUCgYATKzfxpOUdoyfUFjax35QW4pctoAE4fF4OVuYb
4ZtZXSuw2mLba+5AXC4obmA+gX7q1zxaQknP89S4aL9jl3mv23kp/LHP6YTtG6Pk
TDJtvccFopVL9hTk1JHizMYiArHliZZ2hy2MuRXc4fz4cfbHHfGT96mcjhayC5NG
4CjKAQKBgC9+IejTbSX81oi6ZLLJ4jV/PHjUnvsydM7ri/JQQdl6OqdgkBQOsN5F
g8q/tHQfeJbEoHTJUxPoCeMar/F30A2BFT0aKr2A9rcwDKF4WZl/zb5MkwP6o/rs
jM1rsSGIP5FFS056O92OpA3f3r7MPd2LFTBrQoxNIIqn9Lq+F+dX
-----END RSA PRIVATE KEY-----`

func TestIssuerDryRun(t *testing.T) {
	ctx := ctxutil.SetDryRun(context.Background(), true)
	_, err := github.CreateIssue(ctx, model.Alert{}, model.ActionArgs{
		"app_id":             float64(123),
		"install_id":         float64(123),
		"secret_private_key": dummyPrivateKey,
		"owner":              "owner",
		"repo":               "repo",
	})
	gt.NoError(t, err)
}

func TestIssuerValidationFail(t *testing.T) {
	testCases := map[string]struct {
		cfg model.ActionArgs
	}{
		"missing app_id": {
			cfg: model.ActionArgs{
				"install_id":         float64(123),
				"secret_private_key": dummyPrivateKey,
				"owner":              "owner",
				"repo":               "repo",
			},
		},
		"missing install_id": {
			cfg: model.ActionArgs{
				"app_id":             float64(123),
				"secret_private_key": dummyPrivateKey,
				"owner":              "owner",
				"repo":               "repo",
			},
		},
		"missing private_key": {
			cfg: model.ActionArgs{
				"app_id":     float64(123),
				"install_id": float64(123),
				"owner":      "owner",
				"repo":       "repo",
			},
		},
		"missing owner": {
			cfg: model.ActionArgs{
				"app_id":             float64(123),
				"install_id":         float64(123),
				"secret_private_key": dummyPrivateKey,
				"repo":               "repo",
			},
		},
		"missing repo": {
			cfg: model.ActionArgs{
				"app_id":             float64(123),
				"install_id":         float64(123),
				"secret_private_key": dummyPrivateKey,
				"owner":              "owner",
			},
		},
		"app_id is not a float64": {
			cfg: model.ActionArgs{
				"app_id":             "123",
				"install_id":         float64(123),
				"secret_private_key": dummyPrivateKey,
				"owner":              "owner",
				"repo":               "repo",
			},
		},
		"install_id is not a float64": {
			cfg: model.ActionArgs{
				"app_id":             float64(123),
				"install_id":         "123",
				"secret_private_key": dummyPrivateKey,
				"owner":              "owner",
				"repo":               "repo",
			},
		},
		"app_id is not a float64, but int": {
			cfg: model.ActionArgs{
				"app_id":             123,
				"install_id":         float64(123),
				"secret_private_key": dummyPrivateKey,
				"owner":              "owner",
				"repo":               "repo",
			},
		},
		"install_id is not a float64, but int": {
			cfg: model.ActionArgs{
				"app_id":             float64(123),
				"install_id":         123,
				"secret_private_key": dummyPrivateKey,
				"owner":              "owner",
				"repo":               "repo",
			},
		},
		"private_key is not RSA format": {
			cfg: model.ActionArgs{
				"app_id":             float64(123),
				"install_id":         float64(123),
				"secret_private_key": "xxx",
				"owner":              "owner",
				"repo":               "repo",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctx := ctxutil.SetDryRun(context.Background(), true)
			_, err := github.CreateIssue(ctx, model.Alert{}, tc.cfg)
			gt.Error(t, err)
		})
	}
}

func TestRenderReference(t *testing.T) {
	t.Run("test rendering .Refs by template", func(t *testing.T) {
		alert := model.Alert{
			AlertMetaData: model.AlertMetaData{
				Refs: []model.Reference{
					{
						Title: "Test Title",
						URL:   "http://test.url",
					},
				},
			},
		}

		var buf bytes.Buffer
		if err := github.ExecuteTemplate(&buf, alert); err != nil {
			t.Fatal(err)
		}

		got := buf.String()
		if !strings.Contains(got, "# References") {
			t.Errorf("expected 'References' as H1 to be included in the output, got %s", got)
		}
		if !strings.Contains(got, "- [Test Title](http://test.url)") {
			t.Errorf("expected title and link 'http://test.url' to be included in the output, got %s", got)
		}
	})
}
