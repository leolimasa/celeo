package chat

import (
	"fmt"
	"net"
	"strings"

	"github.com/leolimasa/celeo/cell"
)

// ---------------------------------------------------------------------------
// Outputs
// ---------------------------------------------------------------------------

type TcpLineReaderOutMsg interface {
	TcpLineReaderOutMsgType() string
}

type TcpLineReaderReceived struct {
	Line string
}

func (s TcpLineReaderReceived) TcpLineReaderOutMsgType() string {
	return "TcpLineReaderOutMsg"
}

type TcpLineReaderState struct {
	buffer string
}

type TcpLineReader = cell.Cell[any, TcpLineReaderOutMsg]

// ---------------------------------------------------------------------------
// Cell
// ---------------------------------------------------------------------------

func NewTcpLineReader(conn net.Conn, out chan TcpLineReaderOutMsg) *TcpLineReader {
	self := cell.NewCell[any, TcpLineReaderOutMsg](out)
	go func() {
		defer self.Destroy()
		tcpReader := NewTcpReader(conn, make(chan TcpReaderOutMsg, 512))
		self.AddChild(tcpReader)
		state := TcpLineReaderState{
			buffer: "",
		}
		for {
			select {
			case <-self.Stop:
				return
			case readerMsg := <-tcpReader.Output:
				switch msg := readerMsg.(type) {
				case TcpReaderBytesReceived:
					str := string(msg.Bytes)
					newLinePos := strings.Index(str, "\n")
					if newLinePos != -1 {
						prefix := str[0:newLinePos]
						suffix := str[newLinePos+1:]
						self.Output <- TcpLineReaderReceived{Line: state.buffer + prefix}
						fmt.Print("Line sent.")
						state.buffer = suffix
					}
				}
			}
		}
	}()
	return self
}
