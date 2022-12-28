package chat 

import (
	"fmt"
	"github.com/leolimasa/celeo/cell"
)


type ChatAppState struct {
	tcpListener *TcpListener
	messages    []string
	clients     *Client
}

func NewChatApp(address string) *cell.Cell[any, any] {
	self := cell.NewCell[any](make(chan any, 512))
	go func() {
		defer self.Destroy()
		tcpListener := NewTcpListener(address, make(chan TcpListenerOutMsg, 512))
		clients := cell.NewCellCollection[ClientInMsg, ClientOutMsg]()

		self.AddChild(tcpListener)
		self.AddChild(clients)

		state := ChatAppState{
			tcpListener: tcpListener,
			clients:     clients,
		}
		for {
			select {
			case <-self.Stop:
				return
			case tcpListenerMsg := <-state.tcpListener.Output:
				switch msg := tcpListenerMsg.(type) {
				case TcpListenerNewConnection:
					client := NewClient(msg.Conn, clients.Output)
					// Send initial messages
					for _, msg := range state.messages {
						client.Input <- ClientSendMessage{Contents: msg}
					}
					state.clients.AddChild(client)
					fmt.Println("Client connected.")
				}
			case clientMsg := <-state.clients.Output:
				switch msg := clientMsg.(type) {
				case ClientReceiveMessage:
					// Broadcast message to all clients
					state.clients.Input <- ClientSendMessage {Contents: msg.Contents}
				case ClientClosed:
					fmt.Println("Client quit.")
				}
			}
		}
	}()
	return self
}

