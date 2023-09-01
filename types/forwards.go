package types

import (
	"fmt"
	"net"
	"strconv"
)

var Forwards = []Forward{}

type Protocol string

const (
	ProtocolHTTPS   Protocol = "https"
	ProtocolHTTP    Protocol = "http"
	ProtocolTCP     Protocol = "tcp"
	ProtocolDefault Protocol = ""
	// ProtocolUDP    Protocol = "udp"
	// ProtocolTCPUDP Protocol = "tcp_udp"
)

const (
	ForceAAAA  DNSForceType = "AAAA"
	ForceA     DNSForceType = "A"
	ForceUnset DNSForceType = ""
)

type (
	DNSSetting struct {
		Cache     uint16       `yaml:"cache" json:"cache"`           // seconds, defalut or 0 will disable update IP
		Server    string       `yaml:"server" json:"server"`         // use system setting by default
		ForceType DNSForceType `yaml:"force_type" json:"force_type"` // force to use A or AAAA record
	}
	DNSForceType string
)

type Forward struct {
	ListenPort     uint16   `yaml:"listen_prot" json:"listen_prot"`
	ListenProtocol Protocol `yaml:"listen_protocol" json:"listen_protocol"`

	DstPort     uint16   `yaml:"dst_prot" json:"dst_prot"`
	DstHost     string   `yaml:"dst_host" json:"dst_host"` // ipv4 or ipv6 or domain
	DstProtocol Protocol `yaml:"dst_protocol" json:"dst_protocol"`

	DNS DNSSetting `yaml:"dns" json:"dns"`

	TLSCert        Cert   `yaml:"tls_cert" json:"tls_cert"`
	TLSRef         string `yaml:"tls_ref" json:"tls_ref"` // perferred, reference to TLS certificate
	IgnoreTLSError bool   `yaml:"ignore_tls_error" json:"ignore_tls_error"`

	PortMapping bool `yaml:"port_mapping" json:"port_mapping"` // auto port mapping
}

func (f Forward) ListenIPv4Addr() string {
	return net.JoinHostPort("127.0.0.1", strconv.Itoa(int(f.ListenPort)))
}
func (f Forward) ListenIPv6Addr() string {
	return net.JoinHostPort("::1", strconv.Itoa(int(f.ListenPort)))
}
func (f Forward) DstAddr() string {
	return net.JoinHostPort(f.DstHost, strconv.Itoa(int(f.DstPort)))
}
func (f Forward) DstUri() string {
	return fmt.Sprintf("%s://%s", f.DstProtocol, f.DstAddr())
}
func (f *Forward) LoadTLSRef() error {
	if f.TLSRef != "" {
		if cert, ok := Certs[f.TLSRef]; !ok {
			return fmt.Errorf("TLS ref %s not found", f.TLSRef)
		} else {
			f.TLSCert = cert
		}
	} else {
		if f.TLSCert.TLSCertPath == "" || f.TLSCert.TLSKeyPath == "" {
			return fmt.Errorf("TLS cert or key file unset")
		}
	}
	return nil
}
