package bridge

import (
	"net"

	"github.com/senayuki/carrier/pkg/consts"
	"github.com/senayuki/carrier/pkg/log"
	"github.com/senayuki/carrier/types"
	"go.uber.org/zap"
)

type UDP struct {
	*types.Forward
	udpListener *net.UDPConn

	logger *zap.Logger
}

func (p *UDP) Close() error {
	if err := p.udpListener.Close(); err != nil {
		return err
	}
	return nil
}

func (p *UDP) Start() error {
	p.logger = log.Logger(consts.UDPProxy).With(zap.Int16(consts.ListenPort, int16(p.ListenPort)),
		zap.Int16(consts.DstPort, int16(p.DstPort)), zap.String(consts.DstUri, p.DstUri()))

	if p.Forward == nil {
		p.logger.Fatal("Forward must be provided")
	}

	// listen ipv4 and ipv6
	udp4Addr, err := net.ResolveUDPAddr("udp", p.ListenIPv4Addr())
	if err != nil {
		p.logger.Fatal("Resolve UDP failed", zap.Error(err))
	}
	if udp4Listener, err := net.ListenUDP("udp", udp4Addr); err != nil {
		p.logger.Fatal("Listen UDP failed", zap.Error(err))
	} else {
		p.udpListener = udp4Listener
	}
	p.logger.Info("Start listening connections")

	go func() {
		for {
			buf := make([]byte, 1024)
			n, src, _ := p.udpListener.ReadFromUDP(buf)
			p.logger.Info("New UDP connection", zap.String(consts.SourceAddr, src.String()))
			dstAddr, _ := net.ResolveUDPAddr("udp", p.DstAddr())
			_, _ = p.udpListener.WriteToUDP(buf[:n], dstAddr)
		}
	}()
	return nil
}
