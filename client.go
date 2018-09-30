package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type IgorClient struct {
	server   *IgorServer
	conn     *websocket.Conn
	sendChan chan *IgorMsg
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (c *IgorClient) readPump() {
	defer func() {
		c.server.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		igormsg := new(IgorMsg)
		err := c.conn.ReadJSON(&igormsg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.server.incoming <- &IgorServerMsg{c, *igormsg}
	}
}

// func (s *IgorServer) readPump() {
// 	go func() {
// 		for {
// 			select {
// 			case msg := <-c.conn.readChan:
// 				// This is an incoming message
// 				switch msg.Command {
// 				case "request":
// 					fmt.Println("Request recieved")
// 					if msg.Args["for"] == "brains" {
// 						s.writeChan <- newIgorMsg("brains", nil, s.brains.Brains)
// 					}
// 				}
// 			}

// 		}
// 	}()
// }

func (c *IgorClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.sendChan:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.WriteJSON(message)

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func newIgorMsg(cmd string, args map[string]string, response interface{}) *IgorMsg {
	igormsg := new(IgorMsg)
	igormsg.Command = cmd
	igormsg.Args = args
	igormsg.Response = response
	return igormsg
}

func servews(s *IgorServer, w http.ResponseWriter, r *http.Request) {
	fmt.Println("In ws")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	client := IgorClient{server: s, conn: conn, sendChan: make(chan *IgorMsg)}
	client.server.register <- &client

	go client.readPump()
	go client.writePump()
}

// 	ticker := time.NewTicker(time.Second)
// 	defer func() {
// 		fmt.Println("Deferred ticker close")
// 		ticker.Stop()
// 		conn.Close()
// 	}()

// 	go func() {
// 		for {
// 			igormsg := new(IgorMsg)
// 			err := conn.ReadJSON(&igormsg)
// 			if err != nil {
// 				log.Println("read:", err)
// 				if websocket.IsCloseError(err, 1006) {
// 					fmt.Println("Closing reader")
// 					s.Closing = true
// 					break
// 				}
// 			} else {
// 				fmt.Printf("Read: %+v\n", igormsg)
// 				s.readChan <- igormsg
// 			}
// 		}
// 	}()

// out:
// 	for {
// 		select {
// 		case msg := <-s.writeChan:
// 			if s.Closing {
// 				break out
// 			}
// 			err := conn.WriteJSON(msg)
// 			if err != nil {
// 				if websocket.IsCloseError(err) {
// 					fmt.Println("Closeing writer")
// 					break out
// 				}
// 				log.Println("write:", err)
// 			} else {
// 				fmt.Printf("Write: %+v\n", msg)
// 			}
// 		case <-ticker.C:
// 			fmt.Println("Tick")
// 			if s.Closing {
// 				break out
// 			}
// 		}
// 	}
// }
