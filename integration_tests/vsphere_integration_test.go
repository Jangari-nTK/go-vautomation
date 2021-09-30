package integration_tests

import (
	"context"
	"testing"
	"time"

	"github.com/Jangari-nTK/go-vautomation/vsphere"
)

const vcUrl = "https://vcsa.api.lab"
const ssoUser = "administrator@vsphere.local"
const ssoPass = "VMware1!"

func TestCreateSessionSuccess(t *testing.T) {
	c, err := vsphere.NewClient(vcUrl, nil)
	if err != nil {
		t.Fatal("failed: cannot create API client")
	}

	ctx := context.Background()
	err = c.CreateSession(ctx, ssoUser, ssoPass)
	if err != nil {
		t.Fatal("failed: cannot create API session", err.Error())
	}

	if c.SessionID == "" {
		t.Fatal("failed: Session id is not set")
	}
}

func TestCreateSessionFailedByIncorrectCredential(t *testing.T) {
	c, err := vsphere.NewClient(vcUrl, nil)
	if err != nil {
		t.Fatal("failed: cannot create API client")
	}

	ctx := context.Background()
	err = c.CreateSession(ctx, ssoUser, "incorrectPassword")
	if err == nil {
		t.Fatal("failed: unexpected authentication occured")
	}
}

func TestGetVcenterTlsSuccess(t *testing.T) {
	c, err := vsphere.NewClient(vcUrl, nil)
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
	c, err := vsphere.NewClient(vcUrl, nil)
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
	c, err := vsphere.NewClient(vcUrl, nil)
	if err != nil {
		t.Fatal("failed: cannot create API client")
	}

	ctx := context.Background()
	err = c.CreateSession(ctx, ssoUser, ssoPass)
	if err != nil {
		t.Fatal("failed: authentication failed")
	}

	spec := vsphere.CertificateManagementVcenterTlsCsrSpec{
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
