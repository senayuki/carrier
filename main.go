package main

import (
	"flag"
	"io/ioutil"
	"path"

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
	logger.Info("Carrier starting")
	// load configs
	{
		// single config
		config := flag.String("config", "", "path to config YAML file")
		// config dir
		configDir := flag.String("configDir", "./conf", "directory which contains config YAML files, non-recursion")
		// acme dir
		acmeDir := flag.String("acmeDir", "./acme", "directory which storage ACME datas")
		flag.Parse()

		initConfig := types.Config{}
		if config != nil && *config != "" {
			logger.Info("loading config file", zap.String(consts.Config, *config))
			data, err := ioutil.ReadFile(*config)
			if err != nil {
				logger.Fatal("read config file failed", zap.Error(err), zap.String(consts.Config, *config))
			}
			conf := types.Config{}
			yaml.Unmarshal(data, &conf)
			initConfig.Forwards = append(initConfig.Forwards, conf.Forwards...)
			initConfig.Certs = append(initConfig.Certs, conf.Certs...)
		}

		if configDir != nil && *configDir != "" {
			logger.Info("loading config files from directory", zap.String(consts.Config, *configDir))
			files, err := ioutil.ReadDir(*configDir)
			if err != nil {
				logger.Fatal("read config failed", zap.Error(err))
			}
			for _, file := range files {
				if file.IsDir() {
					continue
				}
				filename := path.Join(*configDir, file.Name())
				logger.Info("loading config file", zap.String(consts.Config, filename))
				data, err := ioutil.ReadFile(filename)
				if err != nil {
					logger.Fatal("read config file failed", zap.Error(err), zap.String(consts.Config, filename))
				}
				conf := types.Config{}
				yaml.Unmarshal(data, &conf)
				initConfig.Forwards = append(initConfig.Forwards, conf.Forwards...)
				initConfig.Certs = append(initConfig.Certs, conf.Certs...)
			}
		}

		types.ConfigInstance = initConfig
		types.ConfigInstance.ACMEDir = *acmeDir
		err := bridge.PreloadCerts(types.ConfigInstance)
		if err != nil {
			logger.Fatal("load certs config failed", zap.Error(err))
		}
	}

	logger.Info("Carrier config loaded", zap.Int(consts.ForwardCount, len(types.ConfigInstance.Forwards)), zap.Int(consts.CertCount, len(types.ConfigInstance.Certs)))

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
			conn.Start()
			defer conn.Close()
		}
	}
	go lego.RenewCron()
	close := make(chan struct{})
	<-close
}
