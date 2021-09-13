package vsphere

import (
	"testing"

	vsphere "github.com/Jangari-nTK/go-vsphere-automation/vsphere"
)

func TestNewClient(t *testing.T) {
	var c vsphere.Client
	vsphere.NewClient()
}
