//go:build client
// +build client

package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strings"
	"time"
)

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Could not connect:", err)
	}
	defer client.Close()

	var joinReply struct{ ID int }
	err = client.Call("ChatServer.Join", struct{}{}, &joinReply)
	if err != nil {
		log.Fatal("Join error:", err)
	}
	fmt.Printf("Connected as User %d\n", joinReply.ID)

	go func() {
		for {
			time.Sleep(300 * time.Millisecond)
			var messages []string
			err := client.Call("ChatServer.GetUpdates", joinReply.ID, &messages)
			if err == nil {
				for _, msg := range messages {
					fmt.Println("\n", msg)
					fmt.Print("You: ")
				}
			}
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("You: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "exit" {
			fmt.Println("Exiting chat.")
			break
		}

		args := struct {
			ID   int
			Text string
		}{joinReply.ID, text}

		err := client.Call("ChatServer.SendMessage", args, &struct{}{})
		if err != nil {
			fmt.Println("Failed to send message:", err)
		}
	}
}
