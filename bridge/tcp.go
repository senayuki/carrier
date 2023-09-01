package bridge

import (
	"io"
	"net"

	"github.com/senayuki/carrier/pkg/consts"
	"github.com/senayuki/carrier/pkg/log"
	"github.com/senayuki/carrier/types"
	"go.uber.org/zap"
)

type TCP struct {
	*types.Forward
	tcpListener net.Listener

	logger *zap.Logger
}

func (p *TCP) Close() error {
	if err := p.tcpListener.Close(); err != nil {
		return err
	}
	return nil
}

func (p *TCP) Start() error {
	p.logger = log.Logger(consts.TCPProxy).With(zap.Int16(consts.ListenPort, int16(p.ListenPort)),
		zap.Int16(consts.DstPort, int16(p.DstPort)), zap.String(consts.DstUri, p.DstUri()))

	if p.Forward == nil {
		p.logger.Fatal("Forward must be provided")
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", p.ListenIPv4Addr())
	if err != nil {
		p.logger.Fatal("Resolve TCP failed", zap.Error(err))
	}
	if tcpListener, err := net.ListenTCP("tcp", tcpAddr); err != nil {
		p.logger.Fatal("Listen TCP failed", zap.Error(err))
	} else {
		p.tcpListener = tcpListener
	}

	p.logger.Info("Start listening connections")

	go func() {
		for {
			conn, err := p.tcpListener.Accept()
			if err != nil {
				p.logger.Error("Accepting connection failed", zap.Error(err))
				continue
			}
			go p.handleTCP(conn)
		}
	}()
	return nil
}

func (p TCP) handleTCP(tcpConn net.Conn) {
	defer tcpConn.Close()

	p.logger.Info("New TCP connection", zap.String(consts.SourceAddr, tcpConn.RemoteAddr().String()))
	serverConn, err := net.Dial("tcp", p.DstAddr())
	if err != nil {
		p.logger.Error("Connect to dst failed", zap.Error(err), zap.String(consts.DstAddr, p.DstAddr()))
		_ = tcpConn.Close()
		return
	}
	p.pipe(tcpConn, serverConn)
}

// pipe from local socket to remote socket
func (p TCP) pipe(src net.Conn, dst net.Conn) {
	errChan := make(chan error, 1)
	onClose := func(err error) {
		_ = dst.Close()
		_ = src.Close()
	}
	go func() {
		_, err := io.Copy(src, dst)
		errChan <- err
		onClose(err)
	}()
	go func() {
		_, err := io.Copy(dst, src)
		errChan <- err
		onClose(err)
	}()
	<-errChan
}
