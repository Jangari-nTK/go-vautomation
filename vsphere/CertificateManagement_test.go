package vsphere

import (
	"context"
	"testing"
	"time"
)

func TestGetVcenterTlsSuccess(t *testing.T) {
	c, err := NewClient(vcUrl, nil)
	if err != nil {
		t.Fatal("failed: cannot create API client")
	}

	ctx := context.Background()
	err = c.CreateSession(ctx, ssoUser, ssoPass)
	if err != nil {
		t.Fatal("failed: authentication failed")
	}

	_, err = c.GetVcenterTls(ctx)
	if err != nil {
		t.Fatal("failed: could not retrieve TLS certificate info", err.Error())
	}
}

func TestRenewVcenterTlsSuccess(t *testing.T) {
	c, err := NewClient(vcUrl, nil)
	if err != nil {
		t.Fatal("failed: cannot create API client")
	}

	ctx := context.Background()
	err = c.CreateSession(ctx, ssoUser, ssoPass)
	if err != nil {
		t.Fatal("failed: authentication failed")
	}

	tlsInfo, err := c.GetVcenterTls(ctx)
	if err != nil {
		t.Fatal("failed: could not retrieve TLS certificate info", err.Error())
	}
	valid_from, _ := time.Parse("2006-01-02T15:04:05.999Z", tlsInfo.Valid_From)

	err = c.RenewVcenterTls(ctx, 730)
	if err != nil {
		t.Fatal("failed: ", err.Error())
	}

	t.Log("Waiting for vCenter Server services to be rebooted")

	for {
		time.Sleep(1 * time.Minute)
		err = c.CreateSession(ctx, ssoUser, ssoPass)
		if err != nil {
			t.Log(err.Error())
			continue
		}
		tlsInfo, err = c.GetVcenterTls(ctx)
		if err == nil {
			break
		}
		t.Log(err.Error())
	}
	new_valid_from, _ := time.Parse("2006-01-02T15:04:05.999Z", tlsInfo.Valid_From)

	if !valid_from.Before(new_valid_from) {
		t.Fatal("failed: invalid valid_from date")
	}
}

func TestCreateVcenterTlsCsrSuccess(t *testing.T) {
	c, err := NewClient(vcUrl, nil)
	if err != nil {
		t.Fatal("failed: cannot create API client")
	}

	ctx := context.Background()
	err = c.CreateSession(ctx, ssoUser, ssoPass)
	if err != nil {
		t.Fatal("failed: authentication failed")
	}

	spec := CertificateManagementVcenterTlsCsrSpec{
		"vcsa.api.lab",
		"US",
		"admin@example.com",
		2048,
		"Palo Alto",
		"VMware",
		"VMware Engineering",
		"California",
		[]string{"vcsa.api.lab", "192.168.0.160"},
	}
	_, err = c.CreateVcenterTlsCsr(ctx, spec)
	if err != nil {
		t.Fatal("failed: cannot retrieve TLS CSR", err.Error())
	}
}
