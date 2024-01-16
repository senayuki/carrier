package lego

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/senayuki/carrier/pkg/consts"
	"github.com/senayuki/carrier/pkg/log"
	"github.com/senayuki/carrier/types"
	"go.uber.org/zap"
)

var logger *zap.Logger

var defaultPath string

func init() {
	logger = log.Logger(consts.LEGO)
}

func New(certConf *types.CertConfig, certData string) (*LegoCMD, error) {
	// Set default path to configPath/cert
	var p = ""
	if certData != "" {
		fileInfo, err := os.Stat(certData)
		if err != nil {
			if !os.IsNotExist(err) {
				return nil, fmt.Errorf("failed to open acmeDir: %s", err.Error())
			}
		} else {
			if !fileInfo.IsDir() {
				return nil, fmt.Errorf("acmeDir is not dir: %s", err.Error())
			}
		}
		p = certData
	} else if cwd, err := os.Getwd(); err == nil {
		p = cwd
	} else {
		p = "."
	}

	defaultPath = p

	lego := &LegoCMD{
		C:    certConf,
		path: p,
		logger: logger.With(
			zap.String("mode", string(certConf.Mode)),
			zap.String("provider", string(certConf.Provider)),
			zap.String("domain", string(certConf.Domain)),
			zap.String("email", certConf.Email),
		),
	}

	return lego, nil
}

// DNSCert cert a domain using DNS API
func (l *LegoCMD) DNSCert() (CertPath string, KeyPath string, err error) {
	defer func() (string, string, error) {
		// Handle any error
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
			return "", "", err
		}
		return CertPath, KeyPath, nil
	}()

	defer func() {
		if err == nil {
			// record instance for renew
			legos[l.C.Domain] = l
		}
	}()

	// Set Env for DNS configuration
	for key, value := range l.C.Env {
		os.Setenv(strings.ToUpper(key), value)
	}

	// First check if the certificate exists
	CertPath, KeyPath, err = checkCertFile(l.C.Domain)
	if err == nil {
		return CertPath, KeyPath, err
	}

	err = l.Run()
	if err != nil {
		return "", "", err
	}
	CertPath, KeyPath, err = checkCertFile(l.C.Domain)
	if err != nil {
		return "", "", err
	}

	return CertPath, KeyPath, nil
}

// HTTPCert cert a domain using http methods
func (l *LegoCMD) HTTPCert() (CertPath string, KeyPath string, err error) {
	defer func() (string, string, error) {
		// Handle any error
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
			return "", "", err
		}
		return CertPath, KeyPath, nil
	}()

	defer func() {
		if err == nil {
			// record instance for renew
			legos[l.C.Domain] = l
		}
	}()

	// First check if the certificate exists
	CertPath, KeyPath, err = checkCertFile(l.C.Domain)
	if err == nil {
		return CertPath, KeyPath, err
	}

	err = l.Run()
	if err != nil {
		return "", "", err
	}

	CertPath, KeyPath, err = checkCertFile(l.C.Domain)
	if err != nil {
		return "", "", err
	}

	return CertPath, KeyPath, nil
}

// RenewCert renew a domain cert
func (l *LegoCMD) RenewCert() (CertPath string, KeyPath string, ok bool, err error) {
	defer func() (string, string, bool, error) {
		// Handle any error
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
			return "", "", false, err
		}
		return CertPath, KeyPath, ok, nil
	}()

	ok, err = l.Renew()
	if err != nil {
		return
	}

	CertPath, KeyPath, err = checkCertFile(l.C.Domain)
	if err != nil {
		return
	}

	return
}

func checkCertFile(domain string) (string, string, error) {
	keyPath := path.Join(defaultPath, "certificates", fmt.Sprintf("%s.key", domain))
	certPath := path.Join(defaultPath, "certificates", fmt.Sprintf("%s.crt", domain))
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return "", "", fmt.Errorf("cert key failed: %s", domain)
	}
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		return "", "", fmt.Errorf("cert cert failed: %s", domain)
	}
	absKeyPath, _ := filepath.Abs(keyPath)
	absCertPath, _ := filepath.Abs(certPath)
	return absCertPath, absKeyPath, nil
}
