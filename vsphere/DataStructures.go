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
	Authority_Information_Access_Uri []string `json:"authority_information_access_uri"`
	Cert                             string   `json:"cert"`
	Extended_Key_Usage               []string `json:"extended_key_usage"`
	Is_CA                            bool     `json:"is_CA"`
	Issuer_DN                        string   `json:"issuer_dn"`
	Key_Usage                        []string `json:"key_usage"`
	Path_Length_Constraint           int      `json:"path_length_constraint"`
	Serial_Number                    string   `json:"serial_number"`
	Signature_Algorithm              string   `json:"signature_algorithm"`
	Subject_Alternative_Name         []string `json:"subject_alternative_name"`
	Subject_DN                       string   `json:"subject_dn"`
	Thumbprint                       string   `json:"thumbprint"`
	Valid_From                       string   `json:"valid_from"`
	Valid_To                         string   `json:"valid_to"`
	Version                          int      `json:"version"`
}
