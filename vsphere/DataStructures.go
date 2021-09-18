package vsphere

type APIError struct {
	Error_Type string                   `json:"error_type"`
	Messages   []map[string]interface{} `json:"messages"`
	Data       interface{}              `json:"data"`
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
