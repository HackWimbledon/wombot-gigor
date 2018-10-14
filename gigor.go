package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	server := NewIgorServer()
	server.run()
}

//IgorServer is a struct for holding Igor server connections
type IgorServer struct {
	clients    map[*IgorClient]bool
	register   chan *IgorClient
	unregister chan *IgorClient
	incoming   chan *IgorServerMsg
	manager    chan *IgorServerMsg

	router       *http.ServeMux
	brains       *Brains
	brainmanager *BrainManager
}

//NewIgorServer initialises brains are returns a working IgorServer
func NewIgorServer() *IgorServer {
	brains := new(Brains)
	err := brains.Initialise()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	brainmanager := new(BrainManager)

	server := &IgorServer{
		clients:      make(map[*IgorClient]bool),
		register:     make(chan *IgorClient),
		unregister:   make(chan *IgorClient),
		incoming:     make(chan *IgorServerMsg),
		manager:      make(chan *IgorServerMsg),
		brains:       brains,
		brainmanager: brainmanager,
	}

	brainmanager.Initialise(server)

	return server
}

func (s *IgorServer) run() {
	go s.startServer()
	for {
		select {
		case client := <-s.register:
			s.clients[client] = true
		case client := <-s.unregister:
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.sendChan)
			}
		case servermsg := <-s.manager:
			for c := range s.clients {
				c.sendChan <- servermsg.message
			}
		case message := <-s.incoming:
			switch message.message.Command {
			case "request":
				if message.message.Args["for"] == "brains" {
					message.client.sendChan <- NewIgorMsg("brains", nil, s.brains.Brains)
				}
			case "start":
				brain := message.message.Args["brain"]
				go s.brainmanager.StartBrain(brain)
			}
		}
	}
}

func (s *IgorServer) startServer() {
	s.router = http.NewServeMux()
	s.router.HandleFunc("/config", s.getConfig)
	s.router.Handle("/", http.FileServer(http.Dir("static")))
	s.router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		servews(s, w, r)
	})

	fmt.Println("Listening on 8080")
	err := http.ListenAndServe(":8080", s.router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// NewIgorMsg creates a igor message for sending
func NewIgorMsg(cmd string, args map[string]string, response interface{}) *IgorMsg {
	igormsg := new(IgorMsg)
	igormsg.Command = cmd
	igormsg.Args = args
	igormsg.Response = response
	return igormsg
}

func (s *IgorServer) getConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	config := new(Config)
	config.WebSocket = "ws://" + r.Host + "/ws"
	json.NewEncoder(w).Encode(config)
}

// Config is a struct for saving and JSON serialising config data
type Config struct {
	WebSocket string `json:"websocket"`
}

// IgorMsg is a struct for JSON serialised messages to and from clients
type IgorMsg struct {
	Command  string            `json:"cmd"`
	Args     map[string]string `json:"args,omitempty"`
	Response interface{}       `json:"resp,omitempty"`
}

// IgorServerMsg is an internal structure for associating a message with a client
type IgorServerMsg struct {
	client  *IgorClient
	message *IgorMsg
}
