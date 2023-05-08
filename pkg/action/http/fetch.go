package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/goerr"
)

func Fetch(ctx *model.Context, _ model.Alert, args model.ActionArgs) (any, error) {
	method, ok := args["method"].(string)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidArgument, "method is required")
	}

	url, ok := args["url"].(string)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidArgument, "url is required")
	}

	var reqBody io.Reader
	if data, ok := args["data"].(string); ok {
		reqBody = strings.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, goerr.Wrap(err, "Fail to create HTTP request")
	}

	if v, ok := args["header"].(map[string]string); ok {
		for k, v := range v {
			req.Header.Add(k, v)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, goerr.Wrap(err, "Fail to send HTTP request")
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			ctx.Logger().Warn("Fail to close HTTP response body", "err", err)
		}
	}()

	var result any
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, goerr.Wrap(err, "Fail to read HTTP response body")
	}

	switch resp.Header.Get("Content-Type") {
	case "application/json":
		if err := json.Unmarshal(respBody, &result); err != nil {
			return nil, goerr.Wrap(err, "Fail to parse JSON response")
		}

	case "application/octet-stream":
		result = respBody

	default:
		result = string(respBody)
	}

	return result, nil
}
