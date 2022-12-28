package chat

import (
	"net"

	"github.com/leolimasa/celeo/cell"
)

// ---------------------------------------------------------------------------
// Outputs
// ---------------------------------------------------------------------------

type TcpReaderOutMsg interface {
	TcpReaderMsgOutType() string
}

type TcpReaderBytesReceived struct {
	Bytes []byte
}

func (t TcpReaderBytesReceived) TcpReaderMsgOutType() string {
	return "TcpReaderMsgOut"
}

type TcpReader = cell.Cell[any, TcpReaderOutMsg]

// ---------------------------------------------------------------------------
// Cell
// ---------------------------------------------------------------------------

func NewTcpReader(conn net.Conn, out chan TcpReaderOutMsg) *TcpReader {
	self := cell.NewCell[any, TcpReaderOutMsg](out)
	go func() {
		defer self.Destroy()
		buffer := make([]byte, 1024)
		for {
			select {
			case <-self.Stop:
				return
			default:
				conn.Read(buffer)
				self.Output <- TcpReaderBytesReceived{Bytes: buffer}
			}
		}
	}()
	return self
}
