package lego_test

import (
	"testing"

	"github.com/senayuki/carrier/pkg/lego"
	"github.com/senayuki/carrier/types"
)

func TestLegoClient(t *testing.T) {
	_, err := lego.New(&types.CertConfig{}, "")
	if err != nil {
		t.Error(err)
	}
}

func TestLegoDNSCert(t *testing.T) {
	lego, err := lego.New(&types.CertConfig{
		Domain:   "node1.test.com",
		Provider: "alidns",
		Email:    "test@gmail.com",
		Env: map[string]string{
			"ALICLOUD_ACCESS_KEY": "aaa",
			"ALICLOUD_SECRET_KEY": "bbb",
		},
	}, "")
	if err != nil {
		t.Error(err)
	}

	certPath, keyPath, err := lego.DNSCert()
	if err != nil {
		t.Error(err)
	}
	t.Log(certPath)
	t.Log(keyPath)
}

func TestLegoHTTPCert(t *testing.T) {
	lego, err := lego.New(&types.CertConfig{
		Mode:   types.CertModeHTTP,
		Domain: "node1.test.com",
		Email:  "test@gmail.com",
	}, "")
	if err != nil {
		t.Error(err)
	}

	certPath, keyPath, err := lego.HTTPCert()
	if err != nil {
		t.Error(err)
	}
	t.Log(certPath)
	t.Log(keyPath)
}

func TestLegoRenewCert(t *testing.T) {
	lego, err := lego.New(&types.CertConfig{
		Domain:   "node1.test.com",
		Email:    "test@gmail.com",
		Provider: "alidns",
		Env: map[string]string{
			"ALICLOUD_ACCESS_KEY": "aaa",
			"ALICLOUD_SECRET_KEY": "bbb",
		},
	}, "")
	if err != nil {
		t.Error(err)
	}
	lego.C.Mode = types.CertModeHTTP
	certPath, keyPath, ok, err := lego.RenewCert()
	if err != nil {
		t.Error(err)
	}
	t.Log(certPath)
	t.Log(keyPath)
	t.Log(ok)

	lego.C.Mode = types.CertModeDNS
	certPath, keyPath, ok, err = lego.RenewCert()
	if err != nil {
		t.Error(err)
	}
	t.Log(certPath)
	t.Log(keyPath)
	t.Log(ok)
}
