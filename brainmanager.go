package main

//BrainManager looks after the brain status and switches between brains
type BrainManager struct {
	AvailableBrains *Brains
}

//Initialise prepares the manager
func (brainmanager *BrainManager) Initialise(brains *Brains) (err error) {
	brainmanager.AvailableBrains = brains
	return nil
}
