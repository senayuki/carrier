package bridge

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"

	"github.com/senayuki/carrier/pkg/consts"
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
	h.logger.Debug("Close connection", zap.String(consts.Path, r.RequestURI), zap.String(consts.Method, r.Method), zap.Int64(consts.Duration, time.Now().Sub(startAt).Milliseconds()))
}

func (h *HTTP) Start() error {
	h.logger = log.Logger(consts.HTTPProxy).With(zap.Int16(consts.ListenPort, int16(h.ListenPort)),
		zap.Int16(consts.DestPort, int16(h.DestPort)), zap.String(consts.DestUri, h.DestUri()))

	targetUrl, err := url.Parse(h.DestUri())
	if err != nil {
		h.logger.Fatal("Parse dest URI failed", zap.Error(err))
	}

	h.proxy = httputil.NewSingleHostReverseProxy(targetUrl)
	h.proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: h.IgnoreTLSError,
		},
	}

	if h.PortMapping {
		go natpmp.AddPortMapping(int(h.ListenPort))
	}

	if h.ListenProtocol == types.ProtocolHTTPS {
		if err := h.LoadTLSRef(); err != nil {
			h.logger.Fatal("LoadTLSRef failed", zap.Error(err))
		}
		h.logger.Info("Start listening HTTPS connections")
		go func() {
			err = http.ListenAndServeTLS(":"+strconv.Itoa(int(h.ListenPort)), h.TLSCert.TLSCertPath, h.TLSCert.TLSKeyPath, h)
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
