package types

import (
	"fmt"
	"net"
	"strconv"
)

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
		Cache     uint16       `yaml:"cache"`      // seconds, defalut or 0 will disable update IP
		Server    string       `yaml:"server"`     // use system setting by default
		ForceType DNSForceType `yaml:"force_type"` // force to use A or AAAA record
	}
	DNSForceType string
)

type Forward struct {
	Name string `yaml:"name"`

	ListenPort     uint16   `yaml:"listen_port"`
	ListenProtocol Protocol `yaml:"listen_protocol"`

	DstPort     uint16   `yaml:"dst_port"`
	DstHost     string   `yaml:"dst_host"` // ipv4 or ipv6 or domain
	DstProtocol Protocol `yaml:"dst_protocol"`

	DNS DNSSetting `yaml:"dns"`

	TLS            ForwardTLS `yaml:"tls"`
	IgnoreTLSError bool       `yaml:"ignore_tls_error"`

	PortMapping bool `yaml:"port_mapping"` // auto port mapping
}
type ForwardTLS struct {
	RefAlias   string `yaml:"ref_alias"` // perferred, reference to alias of cert
	CertConfig `yaml:",inline"`
}

func (f *Forward) Valid() error {
	if f.DstHost == "" || f.DstPort == 0 {
		return fmt.Errorf("invalid dst_host or dst_port")
	}
	if f.DstProtocol == "" {
		return fmt.Errorf("dst_protocol unset")
	}
	if f.ListenPort == 0 {
		f.ListenPort = f.DstPort
	}
	if f.ListenProtocol == "" {
		f.ListenProtocol = f.DstProtocol
	}
	if f.DstProtocol != f.ListenProtocol {
		if !((f.DstProtocol == ProtocolHTTPS && f.ListenProtocol == ProtocolHTTP) || (f.DstProtocol == ProtocolHTTP && f.ListenProtocol == ProtocolHTTPS)) {
			return fmt.Errorf("cannot convert protocol unless http<->https")
		}
	}
	return nil
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
