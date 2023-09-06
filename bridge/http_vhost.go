package bridge

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/senayuki/carrier/pkg/consts"
	"github.com/senayuki/carrier/pkg/log"
	"github.com/senayuki/carrier/pkg/natpmp"
	"github.com/senayuki/carrier/types"
	"go.uber.org/zap"
)

type HTTPVHost struct {
	Hosts          map[string]*HTTP
	ListenPort     uint16
	ListenProtocol types.Protocol
	logger         *zap.Logger
}

func (h *HTTPVHost) Close() error {
	return nil
}
func (h *HTTPVHost) ChooseHost(host string) *HTTP {
	hostConfig, ok := h.Hosts[host]
	if !ok {
		hostConfig, ok = h.Hosts[""] // empty host config as default
		if !ok {
			return nil
		}
	}
	return hostConfig
}
func (h *HTTPVHost) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	logger := h.logger.With(zap.String(consts.Host, host))
	startAt := time.Now()
	hostConfig := h.ChooseHost(host)
	if hostConfig == nil {
		logger.Info("unknown host")
		return
	}
	logger.Info("Receive connection", zap.String(consts.Path, r.RequestURI), zap.String(consts.Method, r.Method))
	hostConfig.proxy.ServeHTTP(w, r)
	logger.Debug("Close connection", zap.String(consts.Path, r.RequestURI), zap.String(consts.Method, r.Method), zap.Int64(consts.Duration, time.Since(startAt).Milliseconds()))
}

// valid hosts set before start
func (h *HTTPVHost) Valid() error {
	// all hosts have to listen on the same protocol
	var hostBefore *HTTP
	// only one default endponint(empty host) per port
	var emptyHostCount int
	for _, host := range h.Hosts {
		if hostBefore == nil {
			hostBefore = host
		}
		if hostBefore.ListenProtocol != host.ListenProtocol {
			return fmt.Errorf("%s(%s) haves different protocol from %s(%s) at port %d", host.Name, host.ListenProtocol, hostBefore.Name, hostBefore.ListenProtocol, host.ListenPort)
		}
		if host.ListenHost == "" {
			emptyHostCount++
			if emptyHostCount > 1 {
				return fmt.Errorf("port %d has more than one host unset", host.ListenPort)
			}
		}
	}
	return nil
}

func (h *HTTPVHost) Start() error {
	h.logger = log.Logger(consts.HTTPVHost).With(zap.Int16(consts.ListenPort, int16(h.ListenPort)))

	err := h.Valid()
	if err != nil {
		h.logger.Error("Valid vhost rule failed", zap.Error(err))
		return err
	}

	for _, host := range h.Hosts {
		// port mapping
		if host.PortMapping {
			go natpmp.AddPortMapping(int(h.ListenPort), "tcp")
			break
		}
	}

	if h.ListenProtocol == types.ProtocolHTTPS {
		h.logger.Info("Start listening HTTPS connections")
		go func() {
			server := &http.Server{Addr: fmt.Sprintf(":%d", h.ListenPort), Handler: h, TLSConfig: &tls.Config{
				// read tls cert realtime
				GetCertificate: func(client *tls.ClientHelloInfo) (*tls.Certificate, error) {
					host := client.ServerName
					hostConfig := h.ChooseHost(host)
					if hostConfig == nil {
						h.logger.Info("unknown host", zap.String(consts.Host, host))
						return nil, fmt.Errorf("unknown host: %s", host)
					}
					cert, err := tls.LoadX509KeyPair(hostConfig.TLS.CertPath, hostConfig.TLS.KeyPath)
					if err != nil {
						return nil, err
					}
					return &cert, nil
				},
			}}
			err := server.ListenAndServeTLS("", "")
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
