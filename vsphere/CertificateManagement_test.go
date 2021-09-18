package vsphere

import (
	"context"
	"fmt"
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
	csrStr, err := c.createVcenterTlsCsr(ctx, spec)
	if err != nil {
		t.Fatal("failed: cannot retrieve TLS CSR", err.Error())
	}
	fmt.Println(csrStr)
}
