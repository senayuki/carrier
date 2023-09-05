package types

type Cert struct {
	Alias      string `yaml:"alias"`
	CertConfig `yaml:",inline"`
}
type CertMode string

const (
	CertModeFile CertMode = "file"
	CertModeDNS  CertMode = "dns"
	CertModeHTTP CertMode = "http"
	CertModeTLS  CertMode = "tls"
)

type CertConfig struct {
	Mode CertMode `yaml:"mode"`
	// for file mode
	CertPath string `yaml:"cert_file"`
	KeyPath  string `yaml:"key_file"`
	// for acme(dns/http) mode
	Provider string            `yaml:"provider"`
	Domain   string            `yaml:"domain"`
	Email    string            `yaml:"email"`
	Env      map[string]string `yaml:"env"`
}
