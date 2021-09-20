package vsphere

import (
	"encoding/json"
)

type Error struct {
	Error_Type string                   `json:"error_type"`
	Messages   []map[string]interface{} `json:"messages"`
	Data       interface{}              `json:"data"`
}

func (e *Error) Error() string {
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		return err.Error()
	}
	return string(jsonBytes)
}

type CertificateManagementVcenterTlsCsrSpec struct {
	Common_Name       string   `json:"common_name"`
	Country           string   `json:"country"`
	Email_Address     string   `json:"email_address"`
	Key_Size          int      `json:"key_size"`
	Locality          string   `json:"locality"`
	Organization      string   `json:"organization"`
	Organization_Unit string   `json:"organization_unit"`
	State_Or_Province string   `json:"state_or_province"`
	Subject_Alt_Name  []string `json:"subject_alt_name"`
}

type CertificateManagementVcenterTlsInfo struct {
	Authority_Information_Access_Uri []string
	Cert                             string
	Extended_Key_Usage               []string
	Is_CA                            bool
	Issuer_DN                        string
	Key_Usage                        []string
	Path_Length_Constraint           int
	Serial_Number                    string
	Signature_Algorithm              string
	Subject_Alternative_Name         []string
	Subject_DN                       string
	Thumbprint                       string
	Valid_From                       string
	Valid_To                         string
	Version                          int
}
