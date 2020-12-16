package server

import (
	"bytes"
	"clientServer/massage"
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

const clientDataSize = 8
const queueSize = 50000
const messageBufferSize = 500

type server struct {
	listener net.Listener
	queue    chan *MessageServer
}

type MessageServer struct {
	massage.Message
}

func NewServer(port string) {
	newServer := &server{
		queue: make(chan *MessageServer, queueSize),
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	newServer.listener = listener
	defer listener.Close()
	defer close(newServer.queue)
	newServer.serve()
}

func (newServer *server) serve() {
	for {
		connection, err := newServer.listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		//write ack
		newServer.sendMessage(connection, "ac\n")
		go newServer.handleClientRequest(connection)
	}
}

func (newServer *server) handleClientRequest(conn net.Conn) {
	newServer.response(conn)
}

func (newServer *server) sendMessage(connection net.Conn, message string) {
	msg := MessageServer{massage.Message{Data: message, Length: len(message)}}
	tempBuffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(tempBuffer)
	_ = encoder.Encode(msg)
	if _, err := connection.Write(tempBuffer.Bytes()); err != nil {
		log.Printf("failed to send the server ack: %v\n", err)
	}
}

func read(conn net.Conn) *MessageServer {
	// create a temp buffer
	tmp := make([]byte, messageBufferSize)
	_, err := conn.Read(tmp)
	if err != nil {
		return nil
	}
	tempBuff := bytes.NewBuffer(tmp)
	messageServerObj := new(MessageServer)
	decoder := gob.NewDecoder(tempBuff)
	// decodes buffer and into a Message struct
	_ = decoder.Decode(messageServerObj)
	cutIfNeeded(messageServerObj)
	return messageServerObj
}
func cutIfNeeded(messageServer *MessageServer) {
	if len(messageServer.Message.Data) > clientDataSize {
		messageServer.Message.Data = messageServer.Message.Data[0:8]
		messageServer.Message.Length = clientDataSize
	}
}

func (newServer *server) response(conn net.Conn) {
	defer conn.Close()
	for {
		itemThatWasRead := read(conn)
		if itemThatWasRead == nil {
			log.Println("client left")
			break
		}
		log.Printf("client entered %s, length: %d bytes\n",
					itemThatWasRead.Data, itemThatWasRead.Length)

		if itemThatWasRead.Data == "POP" {
			select {
			//pop from channel
			case item, ok := <-newServer.queue:
				if ok {
					log.Printf("popped item: %s, length: %d bytes\n", item.Data, item.Length)
					// Responding to the client
					newServer.sendMessage(conn, item.Data+" <- popped value from server\n")
				}
			default:
				log.Println("Client tried to pop but no value ready, moving on.")
				newServer.sendMessage(conn, "You tried to pop but no value ready, moving on.\n")
			}
		} else {
			select {
			case newServer.queue <- itemThatWasRead:
				log.Printf("element entered the queue: %s, length: %d bytes\n",
					itemThatWasRead.Data, itemThatWasRead.Length)
				newServer.sendMessage(conn, itemThatWasRead.Data+"\n")
			default:
				log.Println("Channel full")
				newServer.sendMessage(conn, "Channel full"+"\n")
			}
		}
	}
}
