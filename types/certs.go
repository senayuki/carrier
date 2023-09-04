package types

import (
	"fmt"

	"github.com/go-acme/lego/v4/cmd"
)

type Cert struct {
	Alias      string `yaml:"alias"`
	CertConfig `yaml:",inline"`
}
type CertMode string

const (
	CertModeFile CertMode = "file"
	CertModeDNS  CertMode = "dns"
	CertModeHTTP CertMode = "http"
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
	Params   map[string]string `yaml:"params"`
}

func (c *CertConfig) ACME() error {
	cmd.NewAccountsStorage()
}
func (c *CertConfig) GetCertFile() (certFile, keyFile string, err error) {
	switch c.Mode {
	case CertModeFile:
		if c.CertPath == "" || c.KeyPath == "" {
			return "", "", fmt.Errorf("TLS cert or key file unset")
		}
		return c.CertPath, c.KeyPath, nil
	case CertModeDNS, CertModeHTTP:
		// TODO acme
		return "", "", nil
	default:
		return "", "", fmt.Errorf("unknown cert mode: %s", c.Mode)
	}
}
