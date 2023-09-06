package main

import (
	"flag"
	"os"

	"github.com/senayuki/carrier/bridge"
	"github.com/senayuki/carrier/pkg/consts"
	"github.com/senayuki/carrier/pkg/lego"
	"github.com/senayuki/carrier/pkg/log"
	"github.com/senayuki/carrier/types"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func main() {
	defer log.Sync()
	logger := log.Logger(consts.Main)
	// load configs
	{
		config := flag.String("config", "config.yaml", "path to config YAML file")
		flag.Parse()
		logger.Info("loading config", zap.String(consts.Config, *config))
		data, err := os.ReadFile(*config)
		if err != nil {
			logger.Fatal("read config failed", zap.Error(err))
		}
		yaml.Unmarshal(data, &types.ConfigInstance)
		types.ConfigInstance.ConfigLocation = *config
		err = bridge.PreloadCerts(types.ConfigInstance)
		if err != nil {
			logger.Fatal("load certs config failed", zap.Error(err))
		}
	}
	HTTPVHosts := map[uint16]*bridge.HTTPVHost{}
	for _, forward := range types.ConfigInstance.Forwards {
		forward := forward
		err := forward.Valid()
		if err != nil {
			logger.Error("Valid forward rule failed", zap.Error(err), zap.String(consts.ForwardName, forward.Name))
			continue
		}
		switch forward.DstProtocol {
		case types.ProtocolTCP:
			conn := bridge.TCP{Forward: &forward}
			conn.Start()
			defer conn.Close()
		case types.ProtocolHTTP, types.ProtocolHTTPS:
			conn := bridge.HTTP{Forward: &forward}
			conn.Setup()
			if _, ok := HTTPVHosts[forward.ListenPort]; !ok {
				HTTPVHosts[forward.ListenPort] = &bridge.HTTPVHost{
					Hosts:          map[string]*bridge.HTTP{},
					ListenPort:     forward.ListenPort,
					ListenProtocol: forward.ListenProtocol,
				}
			}
			HTTPVHosts[forward.ListenPort].Hosts[forward.ListenHost] = &conn
			defer conn.Close()
		}
	}
	for port, vhost := range HTTPVHosts {
		err := vhost.Start()
		if err != nil {
			logger.Fatal("start vhost failed", zap.Error(err), zap.Uint16(consts.ListenPort, port))
		}
	}
	go lego.RenewCron()
	close := make(chan struct{})
	<-close
}
