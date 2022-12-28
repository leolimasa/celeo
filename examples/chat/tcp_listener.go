package chat

import (
	"net"

	"github.com/leolimasa/celeo/cell"
)


type TcpListenerOutMsg interface {
	TcpListenerMessageType() string
}

type TcpListenerSocketError struct {
	Error error
}

func (t TcpListenerSocketError) TcpListenerMessageType() string { 
	return "TcpListenerSocketError"
}

type TcpListenerConnectionError struct {
	Error error
}

func (t TcpListenerConnectionError) TcpListenerMessageType() string { 
	return "TcpListenerConnectionError"
}

type TcpListenerNewConnection struct {
	Conn net.Conn
}
func (t TcpListenerNewConnection) TcpListenerMessageType() string { 
	return "TcpListenerNewConnection"
}

type TcpListener = cell.Cell[any, TcpListenerOutMsg]

func NewTcpListener(address string, out chan TcpListenerOutMsg) *TcpListener {
	self := cell.NewCell[any, TcpListenerOutMsg](out)
	go func() {
		defer self.Destroy()
		listener, err := net.Listen("tcp", address)
		if err != nil {
			self.Output <- TcpListenerSocketError{Error: err}
			return
		}
		defer listener.Close()
		// Force listener.Accept to stop blocking if there is a signal in
		// the stop channel
		go func() {
			<-self.Stop
			listener.Close()
		}()
		for {
			conn, err := listener.Accept()
			if err != nil {
				self.Output <- TcpListenerConnectionError{Error: err}
				// TODO handle close error
				return
			}
			self.Output <- TcpListenerNewConnection{Conn: conn}
		}
	}()
	return self
}
