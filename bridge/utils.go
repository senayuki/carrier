package bridge

import (
	"fmt"

	"github.com/senayuki/carrier/types"
)

func PreloadCerts(c types.Config) error {
	// load certs
	types.ConfigInstance.CertsAlias = map[string]types.Cert{}
	for _, cert := range c.Certs {
		if _, ok := types.ConfigInstance.CertsAlias[cert.Alias]; !ok {
			types.ConfigInstance.CertsAlias[cert.Alias] = cert
		} else {
			return fmt.Errorf("ambiguous cert alias: '%s'", cert.Alias)
		}
		if cert.Mode != types.CertModeFile && cert.Mode != "" {
			_, _, err := getCertFile(cert.CertConfig)
			if err != nil {
				return err
			}
		}
	}
	// load standalone certs
	for _, forward := range c.Forwards {
		if forward.TLS.Mode != types.CertModeFile && forward.TLS.Mode != "" {
			_, _, err := getCertFile(forward.TLS.CertConfig)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
