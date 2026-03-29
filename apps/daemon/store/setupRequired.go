package store

import (
	"sync"

	"github.com/corecollectives/mist/models"
)

type SetupState struct {
	setupRequired bool
	mu            sync.RWMutex
}

var state = &SetupState{setupRequired: true}

func InitSetupRequired() error {
	count, err := models.GetUserCount()
	if err != nil {
		return err
	}
	state.mu.Lock()
	state.setupRequired = count == 0
	state.mu.Unlock()
	return nil

}

func SetSetupRequired(setupRequired bool) {
	state.mu.Lock()
	state.setupRequired = setupRequired
	state.mu.Unlock()

}

func IsSetupRequired() bool {
	state.mu.RLock()
	defer state.mu.RUnlock()
	return state.setupRequired
}
