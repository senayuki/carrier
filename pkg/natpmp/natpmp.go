package natpmp

import (
	"net"
	"time"

	"github.com/jackpal/gateway"
	pmp "github.com/jackpal/go-nat-pmp"
	"github.com/senayuki/carrier/pkg/consts"
	"github.com/senayuki/carrier/pkg/log"
	"go.uber.org/zap"
)

var client *pmp.Client
var gatewayIP net.IP
var logger *zap.Logger

func init() {
	logger = log.Logger(consts.NATPMP)
	var err error
	gatewayIP, err = gateway.DiscoverGateway()
	if err != nil {
		logger.Fatal("DiscoverGateway failed", zap.Error(err))
	}
	client = pmp.NewClient(gatewayIP)
}

func AddPortMapping(port int, protocol string) {
	logger := logger.With(zap.Int16(consts.ListenPort, int16(port)), zap.String(consts.Protocol, protocol))

	logger.Info("AddPortMapping")
	_, err := client.AddPortMapping(protocol, port, port, 60*2)
	if err != nil {
		logger.Fatal("AddPortMapping failed", zap.Error(err))
	}
	ticker := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-ticker.C:
			logger.Debug("AddPortMapping renewal")
			_, err := client.AddPortMapping("tcp", port, port, 60*2)
			if err != nil {
				logger.Fatal("AddPortMapping renewal failed", zap.Error(err))
			}
		}
	}
}
