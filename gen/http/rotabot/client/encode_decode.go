// Code generated by goa v3.8.3, DO NOT EDIT.
//
// Rotabot HTTP client encoders and decoders
//
// Command:
// $ goa gen github.com/kevinrobayna/rotabot/design

package client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"

	goahttp "goa.design/goa/v3/http"
)

// BuildHealthcheckRequest instantiates a HTTP request object with method and
// path set to call the "Rotabot" service "Healthcheck" endpoint
func (c *Client) BuildHealthcheckRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: HealthcheckRotabotPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("Rotabot", "Healthcheck", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeHealthcheckResponse returns a decoder for responses returned by the
// Rotabot Healthcheck endpoint. restoreBody controls whether the response body
// should be restored after having been read.
func DecodeHealthcheckResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = io.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body string
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("Rotabot", "Healthcheck", err)
			}
			return body, nil
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("Rotabot", "Healthcheck", resp.StatusCode, string(body))
		}
	}
}

// BuildHomeRequest instantiates a HTTP request object with method and path set
// to call the "Rotabot" service "Home" endpoint
func (c *Client) BuildHomeRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: HomeRotabotPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("Rotabot", "Home", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeHomeResponse returns a decoder for responses returned by the Rotabot
// Home endpoint. restoreBody controls whether the response body should be
// restored after having been read.
func DecodeHomeResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = io.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body string
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("Rotabot", "Home", err)
			}
			return body, nil
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("Rotabot", "Home", resp.StatusCode, string(body))
		}
	}
}
