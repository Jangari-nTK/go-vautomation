package vsphere

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"
)

const vcUrl = "https://vcsa.api.lab"
const ssoUser = "administrator@vsphere.local"
const ssoPass = "VMware1!"
const sessionId = "bc9c22db3ada2cb3c4726effd93e042b"

type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// Create http.Client mock to avoid REST API call
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestNewClient(t *testing.T) {
	// success
	_, err := NewClient(vcUrl, nil)
	if err != nil {
		t.Fatal("failed: cannot create API client", err.Error())
	}

	// failed by malformed URL
	_, err = NewClient("thisismalformedurl", nil)
	if err == nil {
		t.Fatal("failed: API client is created incorrectly", err.Error())
	}
}
func TestNewRequest(t *testing.T) {
	const spath = "/api/session"

	c, err := NewClient(vcUrl, nil)
	if err != nil {
		t.Fatal("failed: cannot create API client")
	}

	// success
	req, err := c.NewRequest(context.Background(), "POST", spath, nil, false)
	if err != nil {
		t.Fatal("failed: cannot create request", err.Error())
	}
	if req.URL.String() != vcUrl+spath {
		t.Fatal("failed: incorrect request url", req.URL.String())
	}
	if req.Header.Get("vmware-api-session-id") != "" {
		t.Fatal("failed: session id is added unexpectedly")
	}

	// success with nil context
	req, err = c.NewRequest(nil, "POST", spath, nil, false)
	if err != nil {
		t.Fatal("failed: cannot create request", err.Error())
	}
	if req.Context() != context.Background() {
		t.Fatal("failed: Incorrect context is set")
	}
}

func TestCreateSession(t *testing.T) {
	httpClient := NewTestClient(func(req *http.Request) (*http.Response, error) {
		parsedVcUrl, _ := url.Parse(vcUrl)
		if parsedVcUrl.Scheme != req.URL.Scheme || parsedVcUrl.Host != req.URL.Host {
			return nil, errors.New("Given URL is unreachable")
		}

		usr, pwd, _ := req.BasicAuth()
		if usr != ssoUser || pwd != ssoPass {
			return &http.Response{
				StatusCode: 401,
				Body:       nil,
				Header:     make(http.Header),
			}, nil
		}

		return &http.Response{
			StatusCode: 201,
			Body:       ioutil.NopCloser(bytes.NewBufferString(fmt.Sprintf("\"%s\"", sessionId))),
			Header:     make(http.Header),
		}, nil
	})

	c, err := NewClient(vcUrl, nil)
	if err != nil {
		t.Fatal("failed: cannot create API client")
	}
	c.HTTPClient = httpClient
	ctx := context.Background()

	// Success
	err = c.CreateSession(ctx, ssoUser, ssoPass)
	if err != nil {
		t.Fatal("failed: cannot create API session", err.Error())
	}
	if c.SessionID != sessionId {
		t.Fatal("failed: Session id is not set")
	}

	// Failed by invalid username
	err = c.CreateSession(ctx, "username@invalid.local", ssoPass)
	if err == nil {
		t.Fatal("failed: unexpected authentication occured")
	}

	// Failed by incorrect password
	err = c.CreateSession(ctx, ssoUser, "incorrectPassword")
	if err == nil {
		t.Fatal("failed: unexpected authentication occured")
	}

	// Failed by invalid URL
	c, err = NewClient("https://this.url.is.unreachable", nil)
	if err != nil {
		t.Fatal("failed: cannot create API client")
	}
	c.HTTPClient = httpClient
	c.HTTPClient.Timeout = 1 * time.Second
	ctx = context.Background()
	err = c.CreateSession(ctx, ssoUser, ssoPass)
	if err == nil {
		t.Fatal("failed: access incorrect URL unexpectedly")
	}
}

func TestGetVcenterTlsSuccess(t *testing.T) {
	httpClient := NewTestClient(func(req *http.Request) (*http.Response, error) {
		if req.Header.Get("vmware-api-session-id") != sessionId {
			t.Log("Invalid SessionID")
			return &http.Response{
				StatusCode: 401,
				Body:       nil,
				Header:     make(http.Header),
			}, nil
		}
		var tlsInfo CertificateManagementVcenterTlsInfo
		json.Unmarshal([]byte(`{
			"authority_information_access_uri": [
			  "URIName: https://vcsa.api.lab/afd/vecs/ca"
			],
			"cert": "-----BEGIN CERTIFICATE-----\nMIIDxDCCAqygAwIBAgIJAOoCufh9g22QMA0GCSqGSIb3DQEBCwUAMIGTMQswCQYD\nVQQDDAJDQTEXMBUGCgmSJomT8ixkARkWB3ZzcGhlcmUxFTATBgoJkiaJk/IsZAEZ\nFgVsb2NhbDELMAkGA1UEBhMCVVMxEzARBgNVBAgMCkNhbGlmb3JuaWExFTATBgNV\nBAoMDHZjc2EuYXBpLmxhYjEbMBkGA1UECwwSVk13YXJlIEVuZ2luZWVyaW5nMB4X\nDTIxMDkzMDExMjI0MloXDTIzMDkzMDExMjI0MlowJDEVMBMGA1UEAwwMdmNzYS5h\ncGkubGFiMQswCQYDVQQGEwJVUzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoC\nggEBANY9IqsvVL84AkFwhZ1QtVPfJW9vhokVO1Uz1NSfoQUutc8EguRW5HtYex50\nqZUErTurtb5p/+FtHrR9SN4SN7o7f9yT2Qi9qsCgcCERDMyXrDYPKC2POoXkawPU\nbmN4zefaKxTNCGu2lW04ZYqiARGdLZNQsFyYPnsTOczitb6CguVQ/hKrrKw3fEBG\n25LC855eLdwJrH4hkLhB3moKR08UQVhuOSPhMqim+3Pt38el/AVmeeI93Accp9U4\nDdvja3jfZF2gkbL6eGNz0aoOQcLf41y4aHOyBFXWeEAeO87G3pOnQ55qGreFs+UI\nYCEuqVBTnILhWKXiSjYv79OEKBkCAwEAAaOBiDCBhTAXBgNVHREEEDAOggx2Y3Nh\nLmFwaS5sYWIwCwYDVR0PBAQDAgOoMB8GA1UdIwQYMBaAFIUUIC6dRGzq3LiBg/ts\n01eSf99PMDwGCCsGAQUFBwEBBDAwLjAsBggrBgEFBQcwAoYgaHR0cHM6Ly92Y3Nh\nLmFwaS5sYWIvYWZkL3ZlY3MvY2EwDQYJKoZIhvcNAQELBQADggEBAIXYoLsMsdDL\nP6J0/xWBJW14lsf+Ht/4fkuHnGizPptIPEV7tq8Z1dON4cqZGi9fmG3N0VV0m/mN\nOI1iE+o8Owj12KIs9I1gj/E641Ipi+DDH+htLAtUDv8biOcClasjkTAAYMGs8Mq8\ng8iLLWLh2KwiQhGd7Af/zPCAM055hQsD25o1v8SZ4pZlwyyvkZq200laA3SbxDZo\nk8fCcXF9tnJIysD+Y2NxK96Qf7u1CaHFUdAQLufu5VHqNEeO/bkwouu7EZ2r4c4I\nI5Oc6PuVqA0cBMvCGsjM/cYrimVZT7iUUrhUr5hDtTAo+0Ewn8n83WwXfE6Enjeo\nmfkQ3qkd4ME=\n-----END CERTIFICATE-----",
			"extended_key_usage": [
			  ""
			],
			"is_CA": false,
			"issuer_dn": "OU=VMware Engineering, O=vcsa.api.lab, ST=California, C=US, DC=local, DC=vsphere, CN=CA",
			"key_usage": [
			  "digitalSignature",
			  "keyEncipherment",
			  "keyAgreement"
			],
			"path_length_constraint": -1,
			"serial_number": "ea02b9f87d836d90",
			"signature_algorithm": "SHA256WITHRSA",
			"subject_alternative_name": [
			  "vcsa.api.lab"
			],
			"subject_dn": "C=US, CN=vcsa.api.lab",
			"thumbprint": "87CA0701ACBB2E43F52C25512151A323228AEDD5",
			"valid_from": "2021-09-30T11:22:42.000Z",
			"valid_to": "2023-09-30T11:22:42.000Z",
			"version": 3
		  }`), &tlsInfo)
		tlsInfoBytes, _ := json.Marshal(tlsInfo)

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBuffer(tlsInfoBytes)),
			Header:     make(http.Header),
		}, nil
	})

	c, err := NewClient(vcUrl, nil)
	if err != nil {
		t.Fatal("failed: cannot create API client")
	}
	c.HTTPClient = httpClient

	// Success
	ctx := context.Background()
	c.SessionID = sessionId
	_, err = c.GetVcenterTls(ctx)
	if err != nil {
		t.Fatal("failed: could not retrieve TLS certificate info", err.Error())
	}

	// Failed by unauthenticated user
	c.SessionID = "invalidSessionId"
	_, err = c.GetVcenterTls(ctx)
	if err == nil {
		t.Fatal("failed: unauthenticated access passed unexpectedly")
	}
}
