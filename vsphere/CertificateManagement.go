package vsphere

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/pkg/errors"
)

func (c *Client) RenewVcenterTls(ctx context.Context, duration int) error {
	if duration > 730 {
		return errors.New("invalid duration")
	}
	if duration <= 0 {
		return errors.New("duration must be greater than 0")
	}
	jsonBytes, _ := json.Marshal(map[string]int{
		"duration": duration,
	})
	req, err := c.NewRequest(ctx, "POST", "/api/vcenter/certificate-management/vcenter/tls", bytes.NewBuffer(jsonBytes), true)
	if err != nil {
		return err
	}
	params := req.URL.Query()
	params.Add("action", "renew")
	req.URL.RawQuery = params.Encode()

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 204 {
		var apiError Error
		if err := decodeBody(res, &apiError); err != nil {
			return err
		}
		return &apiError
	}

	return nil
}

func (c *Client) GetVcenterTls(ctx context.Context) (*CertificateManagementVcenterTlsInfo, error) {
	req, err := c.NewRequest(ctx, "GET", "/api/vcenter/certificate-management/vcenter/tls", nil, true)
	if err != nil {
		return nil, err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		var apiError Error
		if err := decodeBody(res, &apiError); err != nil {
			return nil, err
		}
		return nil, &apiError
	}

	var tlsInfo CertificateManagementVcenterTlsInfo
	if err := decodeBody(res, &tlsInfo); err != nil {
		return nil, err
	}
	return &tlsInfo, nil
}

func (c *Client) CreateVcenterTlsCsr(ctx context.Context, spec CertificateManagementVcenterTlsCsrSpec) (string, error) {
	jsonBytes, _ := json.Marshal(spec)
	req, err := c.NewRequest(ctx, "POST", "/api/vcenter/certificate-management/vcenter/tls-csr", bytes.NewBuffer(jsonBytes), true)
	if err != nil {
		return "", err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != 201 {
		var apiError Error
		if err := decodeBody(res, &apiError); err != nil {
			return "", err
		}
		return "", &apiError
	}

	var csrStr map[string]string
	if err := decodeBody(res, &csrStr); err != nil {
		return "", err
	}
	return csrStr["csr"], nil
}
