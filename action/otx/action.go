package otx

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/m-mizutani/goerr/v2"
	"github.com/secmon-lab/alertchain/pkg/ctxutil"
	"github.com/secmon-lab/alertchain/pkg/domain/model"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
)

func ReplaceHTTPClient(client *http.Client) {
	httpClient = client
}

var httpClient = http.DefaultClient

func Indicator(ctx context.Context, _ model.Alert, args model.ActionArgs) (any, error) {
	api_key, ok := args["secret_api_key"].(string)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidArgument, "secret_api_key is required")
	}

	indicatorType, ok := args["type"].(string)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidArgument, "type is required")
	}
	if !isValidType(indicatorType) {
		return nil, goerr.Wrap(types.ErrActionInvalidArgument, "type must be one of ipv4, ipv6, domain, hostname, file, url")
	}

	indicator, ok := args["indicator"].(string)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidArgument, "indicator is required")
	}

	section, ok := args["section"].(string)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidArgument, "section is required")
	}
	if !isValidSection(section) {
		return nil, goerr.Wrap(types.ErrActionInvalidArgument, "section must be one of general, reputation, geo, malware, url_list, passive_dns, http_scans")
	}

	url := "https://otx.alienvault.com/api/v1/indicators/" + indicatorType + "/" + indicator + "/" + section

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, goerr.Wrap(err, "Fail to create HTTP request for OTX")
	}
	req.Header.Set("X-OTX-API-KEY", api_key)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, goerr.Wrap(err, "Fail to send HTTP request to OTX")
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			ctxutil.Logger(ctx).Warn("Fail to close HTTP response body", "err", err)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, goerr.New("OTX returns non-200 status code",
			goerr.V("status", resp.StatusCode),
			goerr.V("url", url),
			goerr.V("body", string(body)),
		)
	}

	var result any
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, goerr.Wrap(err, "Fail to read HTTP response body from OTX")
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, goerr.Wrap(err, "Fail to parse JSON response from OTX", goerr.V("body", string(respBody)))
	}

	return result, nil
}

func isValidSection(section string) bool {
	sections := []string{"general", "reputation", "geo", "malware", "url_list", "passive_dns", "http_scans"}
	for _, s := range sections {
		if section == s {
			return true
		}
	}
	return false
}

func isValidType(t string) bool {
	categories := []string{"IPv4", "IPv6", "domain", "hostname", "file", "url"}
	for _, c := range categories {
		if t == c {
			return true
		}
	}
	return false
}
