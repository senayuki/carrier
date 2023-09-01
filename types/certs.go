package types

var Certs = map[string]Cert{}

type Cert struct {
	TLSCertPath string `json:"cert_file"`
	TLSKeyPath  string `json:"key_file"`
}
