package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/mattn/go-shellwords"
)

//BrainManager looks after the brain status and switches between brains
type BrainManager struct {
	Server *IgorServer
}

//Initialise prepares the manager
func (brainmanager *BrainManager) Initialise(server *IgorServer) (err error) {
	brainmanager.Server = server
	return nil
}

//StartBrain starts a brain in the manager. If one is currently running it will be shutdown.
func (brainmanager *BrainManager) StartBrain(brainid string) {
	// TODO: Check for running brains

	// Get Brain
	brain := brainmanager.Server.brains.Get(brainid)

	go runBrain(brain, brainmanager.Server.manager)
}

func runBrain(brain Brain, serverchan chan (*IgorServerMsg)) {
	select {
	case serverchan <- &IgorServerMsg{nil, NewIgorMsg("starting", map[string]string{"brain": brain.ID}, nil)}:
		fmt.Println("Sent Message")
	default:
		return
	}
	cmdparts, err := shellwords.Parse(brain.Start)
	if err != nil {
		select {
		case serverchan <- &IgorServerMsg{nil, NewIgorMsg("stopped", map[string]string{"brain": brain.ID, "error": "Couldn't parse start"}, nil)}:
			fmt.Println("Sent Message")
			return
		default:
			return
		}
	}
	cmd := exec.Command(cmdparts[0], cmdparts[1:]...)
	//stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Start()

	//read, err := stdout.ReadString(`\n`)

	cmd.Wait()
}

func sendServerMsg(serverchan chan (*IgorServerMsg), cmd string, args map[string]string, response interface{}) {
	select {
	case serverchan <- &IgorServerMsg{nil, NewIgorMsg(cmd, args, response)}:
		fmt.Println("Sent Message")
		return
	default:
		return
	}
}
