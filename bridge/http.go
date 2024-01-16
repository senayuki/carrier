package bridge

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"

	"github.com/senayuki/carrier/pkg/consts"
	"github.com/senayuki/carrier/pkg/lego"
	"github.com/senayuki/carrier/pkg/log"
	"github.com/senayuki/carrier/pkg/natpmp"

	"github.com/senayuki/carrier/types"
	"go.uber.org/zap"
)

type HTTP struct {
	*types.Forward
	proxy *httputil.ReverseProxy

	logger *zap.Logger
}

func (h *HTTP) Close() error {
	return nil
}
func (h *HTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startAt := time.Now()
	h.logger.Info("Receive connection", zap.String(consts.Path, r.RequestURI), zap.String(consts.Method, r.Method))
	h.proxy.ServeHTTP(w, r)
	h.logger.Debug("Close connection", zap.String(consts.Path, r.RequestURI), zap.String(consts.Method, r.Method), zap.Int64(consts.Duration, time.Since(startAt).Milliseconds()))
}

func (h *HTTP) Start() error {
	h.logger = log.Logger(consts.HTTPProxy).With(zap.Int16(consts.ListenPort, int16(h.ListenPort)),
		zap.Int16(consts.DstPort, int16(h.DstPort)), zap.String(consts.DstUri, h.DstUri()),
		zap.String(consts.ForwardName, h.Name))

	targetUrl, err := url.Parse(h.DstUri())
	if err != nil {
		h.logger.Fatal("Parse dst URI failed", zap.Error(err))
	}

	h.proxy = httputil.NewSingleHostReverseProxy(targetUrl)
	h.proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: h.IgnoreTLSError,
		},
	}

	if h.PortMapping {
		go natpmp.AddPortMapping(int(h.ListenPort), "tcp")
	}

	if h.ListenProtocol == types.ProtocolHTTPS {
		if err := loadCert(h.Forward); err != nil {
			h.logger.Fatal("Load cert failed", zap.Error(err))
		}
		h.logger.Info("Start listening HTTPS connections")
		go func() {
			server := &http.Server{Addr: fmt.Sprintf(":%d", h.ListenPort), Handler: h, TLSConfig: &tls.Config{
				// read tls cert realtime
				GetCertificate: func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
					cert, err := tls.LoadX509KeyPair(h.TLS.CertPath, h.TLS.KeyPath)
					if err != nil {
						return nil, err
					}
					return &cert, nil
				},
			}}
			err = server.ListenAndServeTLS("", "")
			if err != nil {
				h.logger.Fatal("ListenAndServeTLS failed", zap.Error(err))
			}
		}()
	} else {
		h.logger.Info("Start listening HTTP connections")
		go func() {
			err := http.ListenAndServe(":"+strconv.Itoa(int(h.ListenPort)), h)
			h.logger.Fatal("ListenAndServe failed", zap.Error(err))
		}()
	}
	return nil
}

func loadCert(f *types.Forward) error {
	var cert types.CertConfig

	if f.TLS.RefAlias != "" {
		if certRef, ok := types.ConfigInstance.CertsAlias[f.TLS.RefAlias]; !ok {
			return fmt.Errorf("TLS ref alias '%s' not found", f.TLS.RefAlias)
		} else {
			cert = certRef.CertConfig
		}
	} else {
		cert = f.TLS.CertConfig
	}

	certPath, keyPath, err := getCertFile(cert)
	if err != nil {
		return fmt.Errorf("get cert file failed: %s", err)
	} else {
		f.TLS.CertPath = certPath
		f.TLS.KeyPath = keyPath
	}
	return nil
}

func getCertFile(cert types.CertConfig) (certFile, keyFile string, err error) {
	switch cert.Mode {
	case types.CertModeFile:
		if cert.CertPath == "" || cert.KeyPath == "" {
			return "", "", fmt.Errorf("TLS cert or key file unset")
		}
		return cert.CertPath, cert.KeyPath, nil
	case types.CertModeDNS:
		lego, err := lego.New(&cert, types.ConfigInstance.ACMEDir)
		if err != nil {
			return "", "", err
		}
		return lego.DNSCert()
	case types.CertModeHTTP, types.CertModeTLS:
		lego, err := lego.New(&cert, types.ConfigInstance.ACMEDir)
		if err != nil {
			return "", "", err
		}
		return lego.HTTPCert()
	default:
		return "", "", fmt.Errorf("unknown cert mode: %s", cert.Mode)
	}
}
