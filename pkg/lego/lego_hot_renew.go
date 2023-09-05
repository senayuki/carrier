package lego

import (
	"time"

	"github.com/senayuki/carrier/pkg/consts"
	"github.com/senayuki/carrier/types"
	"go.uber.org/zap"
)

var legos map[string]*LegoCMD = map[string]*LegoCMD{} // for renew

const RENEW_CRON = "renew-cron"

func checkCert() {
	for _, lego := range legos {
		switch lego.C.Mode {
		case types.CertModeDNS, types.CertModeHTTP, types.CertModeTLS:
			_, _, _, err := lego.RenewCert()
			if err != nil {
				lego.logger.Error("renew cert failed", zap.Error(err), zap.String(consts.SubComponent, RENEW_CRON))
			}
		}
	}
}

func RenewCron() {
	// try to renew every 60mins
	duration := 60 * time.Minute
	for {
		time.Sleep(duration)
		logger.Info("renew cron start", zap.String(consts.SubComponent, RENEW_CRON))
		checkCert()
		logger.Info("renew cron end", zap.String(consts.SubComponent, RENEW_CRON))
	}
}
