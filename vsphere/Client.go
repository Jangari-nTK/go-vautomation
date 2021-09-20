package vsphere

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/pkg/errors"
)

type Client struct {
	URL        *url.URL
	HTTPClient *http.Client

	SessionID string
	Logger    *log.Logger
}

const version = "v0.1"

var userAgent = fmt.Sprintf("XXXGoClient/%s", version)

func decodeBody(res *http.Response, out interface{}) error {
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	return decoder.Decode(out)
}

func NewClient(urlStr string, logger *log.Logger) (*Client, error) {

	parsedURL, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse URL: %s", urlStr)
	}

	var discardLogger = log.New(ioutil.Discard, "", log.LstdFlags)
	if logger == nil {
		logger = discardLogger
	}

	c := &Client{parsedURL, new(http.Client), "", logger}

	return c, nil
}

func (c *Client) ignoreInsecureTlsCertificate(ignore bool) error {
	if ignore {
		c.HTTPClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	return nil
}

func (c *Client) newRequest(ctx context.Context, method, spath string, body io.Reader, useSessionId bool) (*http.Request, error) {
	u := *c.URL
	u.Path = path.Join(c.URL.Path, spath)

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	if ctx == nil {
		ctx = context.Background()
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	if useSessionId {
		req.Header.Set("vmware-api-session-id", c.SessionID)
	}

	return req, nil
}

func (c *Client) createSession(ctx context.Context, username, password string) error {
	req, err := c.newRequest(ctx, "POST", "/api/session", nil, false)
	if err != nil {
		return err
	}

	req.SetBasicAuth(username, password)
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 201 {
		var apiError Error
		if err := decodeBody(res, &apiError); err != nil {
			return err
		}
		return &apiError
	}

	var sessionId string
	if err := decodeBody(res, &sessionId); err != nil {
		return err
	}
	c.SessionID = sessionId

	return nil
}
