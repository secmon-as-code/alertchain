package incident_io

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/goerr"
)

type createAlertRequest struct {
	// Required
	Status string `json:"status"`
	Title  string `json:"title"`

	// Optional
	DeduplicationKey string                 `json:"deduplication_key"`
	Description      string                 `json:"description,omitempty"`
	MetaData         map[string]interface{} `json:"metadata,omitempty"`
	SourceURL        string                 `json:"source_url,omitempty"`
}

const (
	baseURL = "https://api.incident.io"
)

type sinkHTTPClient struct {
	Requests []*http.Request
}

func (x *sinkHTTPClient) Reset() {
	x.Requests = nil
}

func (x *sinkHTTPClient) Do(req *http.Request) (*http.Response, error) {
	x.Requests = append(x.Requests, req)
	resp := httptest.NewRecorder()
	resp.WriteHeader(http.StatusAccepted)
	resp.Write([]byte(`{"deduplication_key":"test_key","message":"test_message","status":"test_status"}`))
	return resp.Result(), nil
}

var sink = &sinkHTTPClient{}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func CreateAlert(ctx *model.Context, alert model.Alert, args model.ActionArgs) (any, error) {
	var (
		apiToken            string
		alertSourceConfigID string
	)
	metaData := map[string]any{}
	request := &createAlertRequest{
		Title:            alert.Title,
		Status:           "firing",
		DeduplicationKey: alert.ID.String(),
		Description:      alert.Description,
		MetaData:         map[string]any{},
	}
	for _, attr := range alert.Attrs {
		request.MetaData[attr.Key.String()] = attr.Value
	}

	if err := args.Parse(
		model.ArgDef("secret_api_token", &apiToken),
		model.ArgDef("alert_source_config_id", &alertSourceConfigID),

		model.ArgDef("status", &request.Status, model.ArgOptional()),
		model.ArgDef("title", &request.Title, model.ArgOptional()),
		model.ArgDef("description", &request.Description, model.ArgOptional()),
		model.ArgDef("deduplication_key", &request.DeduplicationKey, model.ArgOptional()),
		model.ArgDef("metadata", &metaData, model.ArgOptional()),
		model.ArgDef("source_url", &request.SourceURL, model.ArgOptional()),
	); err != nil {
		return nil, err
	}

	for k, v := range metaData {
		request.MetaData[k] = v
	}

	url := baseURL + "/v2/alert_events/http/" + alertSourceConfigID
	raw, err := json.Marshal(request)
	if err != nil {
		return nil, goerr.Wrap(err, "fail to marshal incident.io request")
	}
	reqBody := bytes.NewReader(raw)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, reqBody)
	if err != nil {
		return nil, goerr.Wrap(err, "fail to create incident.io http request")
	}
	req.Header.Set("Authorization", "Bearer "+apiToken)

	var client httpClient = http.DefaultClient
	if ctx.Test() {
		client = sink
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, goerr.Wrap(err, "fail to send incident.io request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, goerr.Wrap(err, "unexpected status code from incident.io").With("code", resp.StatusCode).With("body", string(respBody))
	}

	var respBody struct {
		DeduplicationKey string `json:"deduplication_key"`
		Message          string `json:"message"`
		Status           string `json:"status"`
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, goerr.Wrap(err, "fail to read incident.io response")
	}
	if err := json.Unmarshal(data, &respBody); err != nil {
		return nil, goerr.Wrap(err, "fail to decode incident.io response").With("body", string(data))
	}

	return &respBody, nil
}
