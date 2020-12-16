package client

import (
	"bytes"
	"clientServer/massage"
	"encoding/gob"
	"log"
	"net"
	"os"
	"strings"
)

const messageBufferSize = 500

type Client struct {
	hostAndPort string
	connection   net.Conn
}

type MessageClient struct {
	massage.Message
}

func NewClient(hostAndPort string) *Client {
	newClient := &Client{
		hostAndPort: hostAndPort,
	}

	connection, portNumber := net.Dial("tcp", newClient.hostAndPort)
	if portNumber != nil {
		log.Fatal(portNumber)
	}

	newClient.connection = connection
	return newClient
}

func (newClient *Client) AddItem(item string) {
	newItem := strings.TrimSpace(item)
	newClient.sendMessage(newItem)
}

func (newClient *Client) getResponse() *MessageClient {
	tmp := make([]byte, messageBufferSize)
	if _, err := newClient.connection.Read(tmp); err != nil {
		log.Println(err)
		newClient.connection.Close()
		os.Exit(1)
	}
	tempBuff := bytes.NewBuffer(tmp)
	messageClientObj := new(MessageClient)
	decoder := gob.NewDecoder(tempBuff)
	_ = decoder.Decode(messageClientObj)
	return messageClientObj
}

func (newClient *Client) sendMessage(newItem string) {
	msg := MessageClient{massage.Message{Data: newItem, Length: len(newItem)}}
	tempBuffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(tempBuffer)
	encoder.Encode(msg)
	if _, err := newClient.connection.Write(tempBuffer.Bytes()); err != nil {
		log.Printf("failed to send the client request: %v\n", err)
	}
}

func (newClient *Client)PopItem(){
	newClient.AddItem("POP")
	newClient.getResponse()
}

