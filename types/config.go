package types

import "fmt"

var ConfigInstance = Config{}

type Config struct {
	Forwards   []Forward       `yaml:"forwards"`
	Certs      []Cert          `yaml:"certs"`
	CertsAlias map[string]Cert `yaml:"-"`
}

func (c *Config) LoadCerts() error {
	c.CertsAlias = map[string]Cert{}
	for _, cert := range c.Certs {
		if _, ok := c.CertsAlias[cert.Alias]; !ok {
			c.CertsAlias[cert.Alias] = cert
		} else {
			return fmt.Errorf("ambiguous cert alias: '%s'", cert.Alias)
		}
	}
	return nil
}
