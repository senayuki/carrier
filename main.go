package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"

	"github.com/senayuki/carrier/bridge"
	"github.com/senayuki/carrier/pkg/consts"
	"github.com/senayuki/carrier/pkg/log"
	"github.com/senayuki/carrier/types"
	"go.uber.org/zap"
)

func main() {
	defer log.Sync()
	// load configs
	{
		certsJSON := flag.String("certs-json", "certs.json", "path to certs JSON file")
		forwardsJSON := flag.String("forwards-json", "forwards.json", "path to forwards JSON file")
		flag.Parse()
		logger := log.Logger(consts.Main)
		logger.Info("loading certs-json", zap.String(consts.Config, *certsJSON))
		logger.Info("loading forwards-json", zap.String(consts.Config, *forwardsJSON))
		data, err := ioutil.ReadFile(*certsJSON)
		if err != nil {
			logger.Fatal("read certs-json failed", zap.Error(err))
		}
		json.Unmarshal(data, &types.Certs)
		data, err = ioutil.ReadFile(*forwardsJSON)
		if err != nil {
			logger.Fatal("read forwards-json failed", zap.Error(err))
		}
		json.Unmarshal(data, &types.Forwards)
	}
	for _, forward := range types.Forwards {
		forward := forward
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
	close := make(chan struct{})
	<-close
}
