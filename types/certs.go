package types

type Cert struct {
	Alias       string `yaml:"alias"`
	TLSCertPath string `yaml:"cert_file"`
	TLSKeyPath  string `yaml:"key_file"`
}
