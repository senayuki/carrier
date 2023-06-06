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
	tcp4Listener net.Listener
	tcp6Listener net.Listener

	logger *zap.Logger
}

func (p *TCP) Close() error {
	if err := p.tcp4Listener.Close(); err != nil {
		return err
	}
	if err := p.tcp6Listener.Close(); err != nil {
		p.tcp6Listener.Close()
	}
	return nil
}

func (p *TCP) Start() error {
	p.logger = log.Logger(consts.TCPProxy).With(zap.Int16(consts.ListenPort, int16(p.ListenPort)),
		zap.Int16(consts.DestPort, int16(p.DestPort)), zap.String(consts.DestUri, p.DestUri()))

	if p.Forward == nil {
		p.logger.Fatal("Forward must be provided")
	}

	// listen ipv4 and ipv6
	tcp4Addr, err := net.ResolveTCPAddr("tcp4", p.ListenIPv4Addr())
	if err != nil {
		p.logger.Fatal("Resolve TCP over IPv4 failed", zap.Error(err))
	}
	if tcp4Listener, err := net.ListenTCP("tcp4", tcp4Addr); err != nil {
		p.logger.Fatal("Listen TCP over IPv4 failed", zap.Error(err))
	} else {
		p.tcp4Listener = tcp4Listener
	}

	tcp6Addr, err := net.ResolveTCPAddr("tcp6", p.ListenIPv6Addr())
	if err != nil {
		p.logger.Fatal("Resolve TCP over IPv6 failed", zap.Error(err))
	}
	if tcp6Listener, err := net.ListenTCP("tcp6", tcp6Addr); err != nil {
		p.logger.Fatal("Listen TCP over IPv6 failed", zap.Error(err))
	} else {
		p.tcp6Listener = tcp6Listener
	}
	p.logger.Info("Start listening connections")

	go func() {
		for {
			conn, err := p.tcp4Listener.Accept()
			if err != nil {
				p.logger.Error("Accepting IPv4 connection failed", zap.Error(err))
				continue
			}
			go p.handleTCP(conn)
		}
	}()
	go func() {
		for {
			conn, err := p.tcp6Listener.Accept()
			if err != nil {
				p.logger.Error("Accepting IPv6 connection failed", zap.Error(err))
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
	serverConn, err := net.Dial("tcp", p.DestAddr())
	if err != nil {
		p.logger.Error("Connect to dest failed", zap.Error(err), zap.String(consts.DestAddr, p.DestAddr()))
		_ = tcpConn.Close()
		return
	}
	p.pipe(tcpConn, serverConn)
}

// pipe from local socket to remote socket
func (p TCP) pipe(src net.Conn, dest net.Conn) {
	errChan := make(chan error, 1)
	onClose := func(err error) {
		_ = dest.Close()
		_ = src.Close()
	}
	go func() {
		_, err := io.Copy(src, dest)
		errChan <- err
		onClose(err)
	}()
	go func() {
		_, err := io.Copy(dest, src)
		errChan <- err
		onClose(err)
	}()
	<-errChan
}
