package vsphere

import (
	"context"
	"fmt"
	"testing"
)

const vcUrl = "https://vcsa.tanzu.local"
const ssoUser = "administrator@vsphere.local"
const ssoPass = "VMware1!"

func TestNewClientSuccess(t *testing.T) {
	_, err := NewClient(vcUrl, nil)
	if err != nil {
		t.Fatal("failed: cannot create API client", err.Error())
	}
}

func TestNewClientFailedByMalformedURL(t *testing.T) {
	_, err := NewClient("thisismalformedurl", nil)
	if err == nil {
		t.Fatal("failed: API client is created incorrectly", err.Error())
	}
}

func TestNewRequestSuccess(t *testing.T) {
	c, err := NewClient(vcUrl, nil)
	if err != nil {
		t.Fatal("failed: cannot create API client")
	}

	const spath = "/api/session"
	req, err := c.newRequest(context.Background(), "POST", spath, nil, false)
	if err != nil {
		t.Fatal("failed: cannot create request", err.Error())
	}
	if req.URL.String() != vcUrl+spath {
		t.Fatal("failed: incorrect request url", req.URL.String())
	}
	if req.Header.Get("vmware-api-session-id") != "" {
		t.Fatal("failed: session id is added unexpectedly")
	}
}

func TestNewRequestWithNilContext(t *testing.T) {
	c, err := NewClient(vcUrl, nil)
	if err != nil {
		t.Fatal("failed: cannot create API client")
	}

	const spath = "/api/session"
	req, err := c.newRequest(nil, "POST", spath, nil, false)
	if err != nil {
		t.Fatal("failed: cannot create request", err.Error())
	}
	if req.Context() != context.Background() {
		t.Fatal("failed: Incorrect context is set")
	}
}

func TestCreateSessionSuccess(t *testing.T) {
	c, err := NewClient(vcUrl, nil)
	if err != nil {
		t.Fatal("failed: cannot create API client")
	}

	ctx := context.Background()
	err = c.createSession(ctx, ssoUser, ssoPass)
	if err != nil {
		t.Fatal("failed: cannot create API session", err.Error())
	}

	if c.SessionID == "" {
		t.Fatal("failed: Session id is not set")
	}
	fmt.Println(c.SessionID)
}

func TestCreateSessionFailedByIncorrectCredential(t *testing.T) {
	c, err := NewClient(vcUrl, nil)
	if err != nil {
		t.Fatal("failed: cannot create API client")
	}

	ctx := context.Background()
	err = c.createSession(ctx, ssoUser, "incorrectPassword")
	if err == nil {
		t.Fatal("failed: unexpected authentication occured")
	}
	fmt.Println(err)
}
