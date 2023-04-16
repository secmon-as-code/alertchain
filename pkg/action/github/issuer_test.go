package github_test

import (
	"bytes"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/m-mizutani/alertchain/pkg/action/github"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/gt"

	gh "github.com/google/go-github/github"
)

func TestIssueTemplate(t *testing.T) {
	var buf bytes.Buffer
	gt.NoError(t, github.ExecuteTemplate(&buf, model.Alert{
		AlertMetaData: model.AlertMetaData{
			Title:       "blue",
			Description: "orange",
			Params: []model.Parameter{
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
	gt.B(t, strings.Contains(s, "orange")).True()
	gt.B(t, strings.Contains(s, "| magic | `five` |")).True()
	gt.B(t, strings.Contains(s, "| star | `light` |")).True()
	gt.B(t, strings.Contains(s, "| int | `123` |")).True()
	gt.B(t, strings.Contains(s, "| struct | `{bar}` |")).True()
	gt.B(t, strings.Contains(s, `{"foo": "bar"}`)).True()
}

func TestIssuer(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_GITHUB_ISSUER"); !ok {
		t.Skip("Skipping test because TEST_GITHUB_ISSUER is not set")
	}

	cfg := model.ActionConfigValues{
		"app_id":      gt.R1(strconv.Atoi(os.Getenv("TEST_GITHUB_APP_ID"))).NoError(t),
		"install_id":  gt.R1(strconv.Atoi(os.Getenv("TEST_GITHUB_INSTALL_ID"))).NoError(t),
		"private_key": os.Getenv("TEST_GITHUB_PRIVATE_KEY"),
		"owner":       os.Getenv("TEST_GITHUB_OWNER"),
		"repo":        os.Getenv("TEST_GITHUB_REPO"),
	}

	requiredVars := []string{"app_id", "install_id", "private_key", "owner", "repo"}
	for _, key := range requiredVars {
		gt.V(t, cfg[key]).NotEqual("")
	}

	factory := &github.IssuerFactory{}
	issuer := gt.R1(factory.New("test", cfg)).NoError(t)

	ctx := model.NewContext()
	alert := model.Alert{
		AlertMetaData: model.AlertMetaData{
			Title:       "blue",
			Description: "orange",
			Params: []model.Parameter{
				{
					Key:   "magic",
					Value: "five",
				},
			},
		},
		CreatedAt: time.Now(),
		Raw:       `{"foo": "bar"}`,
	}

	params := model.ActionArgs{}
	resp := gt.R1(issuer.Run(ctx, alert, params)).NoError(t)
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

func TestIssuerFactoryPass(t *testing.T) {
	factory := &github.IssuerFactory{}
	action, err := factory.New("test-id", model.ActionConfigValues{
		"app_id":      float64(123),
		"install_id":  float64(123),
		"private_key": dummyPrivateKey,
		"owner":       "owner",
		"repo":        "repo",
	})
	gt.NoError(t, err)
	gt.V(t, action.ID()).Equal("test-id")
}

func TestIssuerFactoryFail(t *testing.T) {
	factory := &github.IssuerFactory{}
	gt.V(t, factory.Name()).Equal("github-issuer")

	testCases := map[string]struct {
		cfg model.ActionConfigValues
	}{
		"missing app_id": {
			cfg: model.ActionConfigValues{
				"install_id":  float64(123),
				"private_key": dummyPrivateKey,
				"owner":       "owner",
				"repo":        "repo",
			},
		},
		"missing install_id": {
			cfg: model.ActionConfigValues{
				"app_id":      float64(123),
				"private_key": dummyPrivateKey,
				"owner":       "owner",
				"repo":        "repo",
			},
		},
		"missing private_key": {
			cfg: model.ActionConfigValues{
				"app_id":     float64(123),
				"install_id": float64(123),
				"owner":      "owner",
				"repo":       "repo",
			},
		},
		"missing owner": {
			cfg: model.ActionConfigValues{
				"app_id":      float64(123),
				"install_id":  float64(123),
				"private_key": dummyPrivateKey,
				"repo":        "repo",
			},
		},
		"missing repo": {
			cfg: model.ActionConfigValues{
				"app_id":      float64(123),
				"install_id":  float64(123),
				"private_key": dummyPrivateKey,
				"owner":       "owner",
			},
		},
		"app_id is not a int": {
			cfg: model.ActionConfigValues{
				"app_id":      "123",
				"install_id":  float64(123),
				"private_key": dummyPrivateKey,
				"owner":       "owner",
				"repo":        "repo",
			},
		},
		"install_id is not a int": {
			cfg: model.ActionConfigValues{
				"app_id":      float64(123),
				"install_id":  "123",
				"private_key": dummyPrivateKey,
				"owner":       "owner",
				"repo":        "repo",
			},
		},
		"private_key is not RSA format": {
			cfg: model.ActionConfigValues{
				"app_id":      float64(123),
				"install_id":  float64(123),
				"private_key": "xxx",
				"owner":       "owner",
				"repo":        "repo",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := factory.New("test", tc.cfg)
			gt.Error(t, err)
		})
	}
}
