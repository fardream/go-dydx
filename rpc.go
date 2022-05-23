package dydx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

func urlJoin(host string, others ...string) string {
	result := host
	for _, v := range others {
		if v == "" {
			continue
		}
		if strings.HasSuffix(host, "/") || strings.HasPrefix(v, "/") {
			result = result + v
		} else {
			result = result + "/" + v
		}
	}
	return result
}

type DydxError struct {
	HttpStatusCode int
	Body           []byte
	Message        string
}

var _ error = (*DydxError)(nil)

func (e *DydxError) Error() string {
	return fmt.Sprintf("http failed: %d %s %s", e.HttpStatusCode, e.Message, e.Body)
}

func getParamsString(input any) (string, error) {
	if input == nil {
		return "", nil
	}

	switch v := input.(type) {
	case string:
		return v, nil
	case url.Values:
		return v.Encode(), nil
	case *url.Values:
		return v.Encode(), nil
	default:
		k, err := query.Values(v)
		if err != nil {
			return "", err
		}

		return k.Encode(), nil
	}
}

func addParamsStr(path, param string) string {
	if len(param) > 0 {
		return fmt.Sprintf("%s?%s", path, param)
	}
	return path
}

func doRequest[TResponse any](ctx context.Context, c *Client, httpMethod, dydxPath string, params any, body []byte, isPublic bool) (*TResponse, error) {
	if c.apiKey == nil {
		return nil, fmt.Errorf("api key is uninitialized")
	}

	timeNow := GetIsoDateStr(time.Now())

	paramstr, err := getParamsString(params)
	if err != nil {
		return nil, fmt.Errorf("failed to get parameter string %#v: %w", params, err)
	}
	path_seg := addParamsStr(fmt.Sprintf("/v3/%s", dydxPath), paramstr)
	fullpath := urlJoin(c.rpcUrl, path_seg)

	timeout_ctx, cancel := context.WithTimeout(ctx, c.timeOut)
	defer cancel()

	log.Debugf("sending %s request to %s", httpMethod, fullpath)

	req, err := http.NewRequestWithContext(timeout_ctx, httpMethod, fullpath, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	if !isPublic {
		signature := c.apiKey.Sign(path_seg, httpMethod, timeNow, body)
		req.Header.Add("DYDX-SIGNATURE", signature)
		req.Header.Add("DYDX-API-KEY", c.apiKey.Key)
		req.Header.Add("DYDX-TIMESTAMP", timeNow)
		req.Header.Add("DYDX-PASSPHRASE", c.apiKey.Passphrase)
	}

	if len(body) > 0 {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	msg, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, &DydxError{HttpStatusCode: resp.StatusCode, Message: resp.Status, Body: msg}
	}

	log.Debugf("msg: %s", msg)

	r := new(TResponse)

	if err := json.Unmarshal(msg, r); err != nil {
		return nil, err
	}

	return r, nil
}
