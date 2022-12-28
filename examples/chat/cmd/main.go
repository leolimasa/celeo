package main

import (
	"fmt"

	"github.com/leolimasa/celeo/examples/chat"
)

func main() {
	fmt.Println("Starting up")
	chatApp := chat.NewChatApp("localhost:3000")
	var val string
	fmt.Scan(&val)
	fmt.Println("Destroying")
	chatApp.Destroy()
}
