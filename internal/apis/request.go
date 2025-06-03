package apis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/buraev/barelog"
)

var WarningError = errors.New("non-critical error encountered during request")

func Request(logPrefix string, client *http.Client, req *http.Request) ([]byte, error) {
	ctx, cancel := context.WithTimeout(req.Context(), 1*time.Minute)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			barelog.Warn(logPrefix, "connection timed out for", req.URL.Path)
			return []byte{}, WarningError
		}
		if errors.Is(err, context.DeadlineExceeded) {
			barelog.Warn(logPrefix, "request timed out for", req.URL.Path)
			return []byte{}, WarningError
		}
		if errors.Is(err, io.ErrUnexpectedEOF) {
			barelog.Warn(logPrefix, "unexpected EOF from", req.URL.Path)
			return []byte{}, WarningError
		}
		if strings.Contains(err.Error(), "read: connection reset by peer") {
			barelog.Warn(logPrefix, "tcp connection reset by peer from", req.URL.Path)
			return []byte{}, WarningError
		}
		return []byte{}, fmt.Errorf("%w sending request to %s failed", err, req.URL.String())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("%w reading response body failed", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		barelog.Warn(
			logPrefix,
			resp.StatusCode,
			fmt.Sprintf("(%s)", strings.ToLower(http.StatusText(resp.StatusCode))),
			"from",
			req.URL.String(),
		)
		return []byte{}, WarningError
	}
	return body, nil
}

func RequestJSON[T any](logPrefix string, client *http.Client, req *http.Request) (T, error) {
	var data T

	body, err := Request(logPrefix, client, req)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		barelog.Debug(string(body))
		return data, fmt.Errorf("%w failed to parse json", err)
	}

	return data, nil
}
