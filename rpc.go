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

// join Urls
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

// DydxError represents a successful HTTP request with status code >= 400
// This can indicates errors like failed authentication with api key or bad parameters.
type DydxError struct {
	HttpStatusCode int
	Body           []byte
	Message        string
}

var _ error = (*DydxError)(nil)

func (e *DydxError) Error() string {
	return fmt.Sprintf("http failed: %d %s %s", e.HttpStatusCode, e.Message, e.Body)
}

// get the parameter string
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

// doRequest is the main function to process the request.
// dydxPath is used in the signing of the request (when isPublic is true)
func doRequest[TResponse any](ctx context.Context, c *Client, httpMethod, dydxPath string, params any, body []byte, isPublic bool) (*TResponse, error) {
	// get parameter string
	paramstr, err := getParamsString(params)
	if err != nil {
		return nil, fmt.Errorf("failed to get parameter string %#v: %w", params, err)
	}

	// the get dydx path
	path_seg := fmt.Sprintf("/v3/%s", dydxPath)
	if len(paramstr) > 0 {
		path_seg = fmt.Sprintf("%s?%s", path_seg, paramstr)
	}

	fullpath := urlJoin(c.rpcUrl, path_seg)

	// setup timeout
	timeout_ctx, cancel := context.WithTimeout(ctx, c.timeOut)
	defer cancel()

	log.Debugf("sending %s request to %s", httpMethod, fullpath)

	req, err := http.NewRequestWithContext(timeout_ctx, httpMethod, fullpath, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// for private, set the authentication headers
	if !isPublic {
		if c.apiKey == nil {
			return nil, fmt.Errorf("api key is uninitialized")
		}
		// timeNow
		timeNow := GetIsoDateStr(time.Now())
		signature := c.apiKey.Sign(path_seg, httpMethod, timeNow, body)
		req.Header.Add("DYDX-SIGNATURE", signature)
		req.Header.Add("DYDX-API-KEY", c.apiKey.Key)
		req.Header.Add("DYDX-TIMESTAMP", timeNow)
		req.Header.Add("DYDX-PASSPHRASE", c.apiKey.Passphrase)
	}

	// set the content to json
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

	log.Debugf("response from remote: %s", msg)

	r := new(TResponse)

	if err := json.Unmarshal(msg, r); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return r, nil
}
