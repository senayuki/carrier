package types

var ConfigInstance = Config{}

type Config struct {
	Forwards       []Forward       `yaml:"forwards"`
	Certs          []Cert          `yaml:"certs"`
	CertsAlias     map[string]Cert `yaml:"-"`
	ConfigLocation string          `yaml:"-"`
}
