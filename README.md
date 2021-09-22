# go-vautomation (WIP)

This repository contains Golang bindings for vSphere Automation API.

Currently, following API bindings in [vCenter REST APIs](https://developer.vmware.com/docs/vsphere-automation/latest/vcenter/index.html) are implemented.

- CIS
  - Session
    - [Create Session](https://developer.vmware.com/docs/vsphere-automation/latest/cis/api/session/post/)
- Certificate Management
  - TLS
    - [Get vCenter TLS](https://developer.vmware.com/docs/vsphere-automation/latest/vcenter/api/vcenter/certificate-management/vcenter/tls/get/)
    - [Renew vCenter TLS](https://developer.vmware.com/docs/vsphere-automation/latest/vcenter/api/vcenter/certificate-management/vcenter/tlsactionrenew/post/)
  - TLS CSR
    - [Create vCenter TLS CSR](https://developer.vmware.com/docs/vsphere-automation/latest/vcenter/api/vcenter/certificate-management/vcenter/tls-csr/post/)