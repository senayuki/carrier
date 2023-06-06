package types

import (
	"fmt"
	"net"
	"strconv"
)

var Certs = map[string]Cert{}
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

type Forward struct {
	ListenPort     uint16   `json:"listen_port"`
	ListenProtocol Protocol `json:"listen_protocol"`

	DestPort     uint16   `json:"dest_port"`
	DestHost     string   `json:"dest_host"` // domain or ipv4 or ipv6
	DestProtocol Protocol `json:"dest_protocol"`

	TLSCert        Cert   `json:"tls_cert"`
	TLSRef         string `json:"tls_ref"` // perferred, reference to TLS certificate
	IgnoreTLSError bool   `json:"ignore_tls_error"`

	PortMapping bool `json:"port_mapping"` // auto port mapping
}

func (f Forward) ListenIPv4Addr() string {
	return net.JoinHostPort("127.0.0.1", strconv.Itoa(int(f.ListenPort)))
}
func (f Forward) ListenIPv6Addr() string {
	return net.JoinHostPort("::1", strconv.Itoa(int(f.ListenPort)))
}
func (f Forward) DestAddr() string {
	return net.JoinHostPort(f.DestHost, strconv.Itoa(int(f.DestPort)))
}
func (f Forward) DestUri() string {
	return fmt.Sprintf("%s://%s", f.DestProtocol, f.DestAddr())
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

type Cert struct {
	TLSCertPath string `json:"cert_file"`
	TLSKeyPath  string `json:"key_file"`
}
