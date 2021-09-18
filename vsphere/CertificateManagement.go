package vsphere

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

func (c *Client) createVcenterTlsCsr(ctx context.Context, spec CertificateManagementVcenterTlsCsrSpec) (string, error) {
	jsonBytes, _ := json.Marshal(spec)
	req, err := c.newRequest(ctx, "POST", "/api/vcenter/certificate-management/vcenter/tls-csr", bytes.NewBuffer(jsonBytes), true)
	if err != nil {
		return "", err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != 201 {
		var apiError APIError
		if err := decodeBody(res, &apiError); err != nil {
			return "", err
		}
		return "", errors.New(
			fmt.Sprintf("status:%s messages:%s id:%s", apiError.Error_Type, apiError.Messages[0]["default_message"], apiError.Messages[0]["id"]),
		)
	}

	var csrStr map[string]string
	if err := decodeBody(res, &csrStr); err != nil {
		return "", err
	}
	return csrStr["csr"], nil
}
