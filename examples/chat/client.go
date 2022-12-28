package chat

import (
	"fmt"
	"net"

	"github.com/leolimasa/celeo/cell"
)

// ---------------------------------------------------------------------------
// Inputs
// ---------------------------------------------------------------------------

type ClientInMsg interface {
	ClientInMsgType() string
}

type ClientSendMessage struct {
	Contents string
}

func (s ClientSendMessage) ClientInMsgType() string {
	return "ClientSendMessage"
}

// ---------------------------------------------------------------------------
// Outputs
// ---------------------------------------------------------------------------

type ClientOutMsg interface {
	ClientOutMsgType() string
}

type ClientReceiveMessage struct {
	Contents string
}

func (s ClientReceiveMessage) ClientOutMsgType() string {
	return "ClientReceiveMessage"
}

type ClientClosed struct{}

func (s ClientClosed) ClientOutMsgType() string {
	return "ClientClosed"
}

// ---------------------------------------------------------------------------
// Cell
// ---------------------------------------------------------------------------

type Client = cell.Cell[ClientInMsg, ClientOutMsg]

type ClientState struct {
	tcpLineReader *TcpLineReader
}

func NewClient(conn net.Conn, out chan ClientOutMsg) *Client {
	self := cell.NewCell[ClientInMsg, ClientOutMsg](out)
	go func() {
		defer self.Destroy()
		tcpLineReader := NewTcpLineReader(conn, make(chan TcpLineReaderOutMsg, 512))
		self.AddChild(tcpLineReader)
		state := ClientState{tcpLineReader: tcpLineReader}
		for {
			select {
			case <-self.Stop:
				return
			case inputMsg := <-self.Input:
				switch msg := inputMsg.(type) {
				case ClientSendMessage:
					conn.Write([]byte(msg.Contents))
				}
			case readerMsg := <-state.tcpLineReader.Output:
				switch msg := readerMsg.(type) {
				case TcpLineReaderReceived:
					fmt.Println("Client received line.")
					switch msg.Line {
					case "quit":
						// TODO come up with a mechanism to auto destroy immediately after the go func
						// returns.
						self.Output <- ClientClosed{}
						return
					default:
						self.Output <- ClientReceiveMessage{Contents: msg.Line}
					}
				}
			}
		}
	}()
	return self
}
