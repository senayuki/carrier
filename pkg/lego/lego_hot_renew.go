package lego

import (
	"log"
	"time"

	"github.com/senayuki/carrier/types"
)

var legos map[string]*LegoCMD = map[string]*LegoCMD{} // for renew

func checkCert() {
	for _, lego := range legos {
		log.Println(lego)
		switch lego.C.Mode {
		case types.CertModeDNS, types.CertModeHTTP, types.CertModeTLS:
			_, _, _, err := lego.RenewCert()
			if err != nil {
				log.Print(err)
			}
		}
	}
}

func RenewCron() {
	// try to renew every 60mins
	duration := 1 * time.Second
	for {
		time.Sleep(duration)
		checkCert()
	}
}
