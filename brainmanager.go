package main

import "fmt"

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
func (brainmanager *BrainManager) StartBrain(brain string) {
	fmt.Println("Start brain")
	select {
	case brainmanager.Server.manager <- &IgorServerMsg{nil, NewIgorMsg("starting", map[string]string{"brain": brain}, nil)}:
		fmt.Println("Sent Message")
	default:
		fmt.Println("Fayl")
	}
	fmt.Println("Message sent back")
}
