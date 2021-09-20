package vsphere

import (
	"context"
	"testing"
)

func TestCreateVcenterTlsCsrSuccess(t *testing.T) {
	c, err := NewClient(vcUrl, nil)
	if err != nil {
		t.Fatal("failed: cannot create API client")
	}

	ctx := context.Background()
	err = c.createSession(ctx, ssoUser, ssoPass)
	if err != nil {
		t.Fatal("failed: authentication failed")
	}

	spec := CertificateManagementVcenterTlsCsrSpec{
		"vcsa.tanzu.local",
		"US",
		"admin@example.com",
		2048,
		"Palo Alto",
		"VMware",
		"VMware Engineering",
		"California",
		[]string{"vcsa.tanzu.local", "192.168.0.1"},
	}
	_, err = c.createVcenterTlsCsr(ctx, spec)
	if err != nil {
		t.Fatal("failed: cannot retrieve TLS CSR", err.Error())
	}
}

func TestRenewVcenterTlsSuccess(t *testing.T) {
	c, err := NewClient(vcUrl, nil)
	if err != nil {
		t.Fatal("failed: cannot create API client")
	}

	ctx := context.Background()
	err = c.createSession(ctx, ssoUser, ssoPass)
	if err != nil {
		t.Fatal("failed: authentication failed")
	}

	err = c.renewVcenterTls(ctx, 730)
	if err != nil {
		t.Fatal("failed: ", err.Error())
	}
}
