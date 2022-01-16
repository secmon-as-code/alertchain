package policy

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/m-mizutani/goerr"
)

type Remote struct {
	url        string
	httpClient HTTPClient
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewRemote(url string) (*Remote, error) {
	return NewRemoteWithHTTPClient(url, nil)
}

func NewRemoteWithHTTPClient(url string, client HTTPClient) (*Remote, error) {
	if err := validation.Validate(url, validation.Required, is.URL); err != nil {
		return nil, goerr.Wrap(err, "invalid URL for remote policy")
	}

	if client == nil {
		client = &http.Client{}
	}

	return &Remote{
		url:        url,
		httpClient: client,
	}, nil
}

type httpInput struct {
	Input interface{} `json:"input"`
}

type httpOutput struct {
	Result interface{} `json:"result"`
}

func (x *Remote) Eval(ctx context.Context, in interface{}, out interface{}) error {
	input := httpInput{
		Input: in,
	}
	rawInput, err := json.Marshal(input)
	if err != nil {
		return goerr.Wrap(err, "fail to marshal rego input for remote inquiry")
	}

	req, err := http.NewRequest(http.MethodPost, x.url, bytes.NewReader(rawInput))
	if err != nil {
		return goerr.Wrap(err, "fail to create a http request for remote inquiry")
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := x.httpClient.Do(req)
	if err != nil {
		return goerr.Wrap(err, "fail http request to OPA server").With("url", x.url)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return goerr.Wrap(err, "unexpected http code from OPA server").
			With("code", resp.StatusCode).
			With("body", string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return goerr.Wrap(err, "fail to read body from OPA server")
	}

	var output httpOutput
	if err := json.Unmarshal(body, &output); err != nil {
		return goerr.Wrap(err, "fail to parse OPA server result").With("body", string(body))
	}

	rawOutput, err := json.Marshal(output.Result)
	if err != nil {
		return goerr.Wrap(err, "fail to re-marshal result filed in OPA response")
	}

	if err := json.Unmarshal(rawOutput, out); err != nil {
		return goerr.Wrap(err, "fail to unmarshal OPA server result to out")
	}

	return nil
}
